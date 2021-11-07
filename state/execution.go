package state

import (
	"errors"
	"fmt"
	"time"

	abci "github.com/reapchain/reapchain-core/abci/types"
	cryptoenc "github.com/reapchain/reapchain-core/crypto/encoding"
	"github.com/reapchain/reapchain-core/libs/fail"
	"github.com/reapchain/reapchain-core/libs/log"
	mempl "github.com/reapchain/reapchain-core/mempool"
	tmstate "github.com/reapchain/reapchain-core/proto/reapchain/state"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
	"github.com/reapchain/reapchain-core/proxy"
	"github.com/reapchain/reapchain-core/types"
)

//-----------------------------------------------------------------------------
// BlockExecutor handles block execution and state updates.
// It exposes ApplyBlock(), which validates & executes the block, updates state w/ ABCI responses,
// then commits and updates the mempool atomically, then saves state.

// BlockExecutor provides the context and accessories for properly executing a block.
type BlockExecutor struct {
	// save state, validators, consensus params, abci responses here
	store Store

	// execute the app against this
	proxyApp proxy.AppConnConsensus

	// events
	eventBus types.BlockEventPublisher

	// manage the mempool lock during commit
	// and update both with block results after commit.
	mempool mempl.Mempool
	evpool  EvidencePool

	logger log.Logger

	metrics *Metrics
}

type BlockExecutorOption func(executor *BlockExecutor)

func BlockExecutorWithMetrics(metrics *Metrics) BlockExecutorOption {
	return func(blockExec *BlockExecutor) {
		blockExec.metrics = metrics
	}
}

// NewBlockExecutor returns a new BlockExecutor with a NopEventBus.
// Call SetEventBus to provide one.
func NewBlockExecutor(
	stateStore Store,
	logger log.Logger,
	proxyApp proxy.AppConnConsensus,
	mempool mempl.Mempool,
	evpool EvidencePool,
	options ...BlockExecutorOption,
) *BlockExecutor {
	res := &BlockExecutor{
		store:    stateStore,
		proxyApp: proxyApp,
		eventBus: types.NopEventBus{},
		mempool:  mempool,
		evpool:   evpool,
		logger:   logger,
		metrics:  NopMetrics(),
	}

	for _, option := range options {
		option(res)
	}

	return res
}

func (blockExec *BlockExecutor) Store() Store {
	return blockExec.store
}

// SetEventBus - sets the event bus for publishing block related events.
// If not called, it defaults to types.NopEventBus.
func (blockExec *BlockExecutor) SetEventBus(eventBus types.BlockEventPublisher) {
	blockExec.eventBus = eventBus
}

// CreateProposalBlock calls state.MakeBlock with evidence from the evpool
// and txs from the mempool. The max bytes must be big enough to fit the commit.
// Up to 1/10th of the block space is allcoated for maximum sized evidence.
// The rest is given to txs, up to the max gas.
func (blockExec *BlockExecutor) CreateProposalBlock(
	height int64,
	state State, commit *types.Commit,
	proposerAddr []byte,
) (*types.Block, *types.PartSet) {

	maxBytes := state.ConsensusParams.Block.MaxBytes
	maxGas := state.ConsensusParams.Block.MaxGas

	evidence, evSize := blockExec.evpool.PendingEvidence(state.ConsensusParams.Evidence.MaxBytes)

	// Fetch a limited amount of valid txs
	maxDataBytes := types.MaxDataBytes(maxBytes, evSize, state.Validators.Size())

	txs := blockExec.mempool.ReapMaxBytesMaxGas(maxDataBytes, maxGas)

	return state.MakeBlock(height, txs, commit, evidence, proposerAddr)
}

// ValidateBlock validates the given block against the given state.
// If the block is invalid, it returns an error.
// Validation does not mutate state, but does require historical information from the stateDB,
// ie. to verify evidence from a validator at an old height.
func (blockExec *BlockExecutor) ValidateBlock(state State, block *types.Block) error {
	err := validateBlock(state, block)
	if err != nil {
		return err
	}
	return blockExec.evpool.CheckEvidence(block.Evidence.Evidence)
}

// ApplyBlock validates the block against the state, executes it against the app,
// fires the relevant events, commits the app, and saves the new state and responses.
// It returns the new state and the block height to retain (pruning older blocks).
// It's the only function that needs to be called
// from outside this package to process and commit an entire block.
// It takes a blockID to avoid recomputing the parts hash.
func (blockExec *BlockExecutor) ApplyBlock(
	state State, blockID types.BlockID, block *types.Block,
) (State, int64, error) {

	if err := validateBlock(state, block); err != nil {
		return state, 0, ErrInvalidBlock(err)
	}

	startTime := time.Now().UnixNano()
	abciResponses, err := execBlockOnProxyApp(
		//@@@logging: executed block
		blockExec.logger, blockExec.proxyApp, block, blockExec.store, state.InitialHeight,
	)
	endTime := time.Now().UnixNano()
	blockExec.metrics.BlockProcessingTime.Observe(float64(endTime-startTime) / 1000000)
	if err != nil {
		return state, 0, ErrProxyAppConn(err)
	}

	fail.Fail() // XXX

	// Save the results before we commit.
	if err := blockExec.store.SaveABCIResponses(block.Height, abciResponses); err != nil {
		return state, 0, err
	}

	fail.Fail() // XXX

	var validators []*types.Validator

	abciStandingMemberUpdates := abciResponses.EndBlock.StandingMemberUpdates
	err = validateStandingMemberUpdates(abciStandingMemberUpdates, state.ConsensusParams.Validator)
	if err != nil {
		return state, 0, fmt.Errorf("error in standing member updates: %v", err)
	}

	standingMemberUpdates, err := types.PB2TM.StandingMemberUpdates(abciStandingMemberUpdates)
	if err != nil {
		return state, 0, err
	}
	if len(standingMemberUpdates) > 0 {
		blockExec.logger.Debug("updates to standing members", "updates", types.StandingMemberListString(standingMemberUpdates))
	}

	abciSteeringMemberCandidateUpdates := abciResponses.EndBlock.SteeringMemberCandidateUpdates
	err = validateSteeringMemberCandidateUpdates(abciSteeringMemberCandidateUpdates, state.ConsensusParams.Validator)
	if err != nil {
		return state, 0, fmt.Errorf("error in steering member candidate updates: %v", err)
	}

	steeringMemberCandidateUpdates, err := types.PB2TM.SteeringMemberCandidateUpdates(abciSteeringMemberCandidateUpdates)
	if err != nil {
		return state, 0, err
	}
	if len(steeringMemberCandidateUpdates) > 0 {
		blockExec.logger.Debug("updates to steering member candidates", "updates", types.SteeringMemberCandidateListString(steeringMemberCandidateUpdates))
	}

	abciQrnUpdates := abciResponses.EndBlock.QrnUpdates
	err = validateQrnUpdates(abciQrnUpdates, state.ConsensusParams.Validator)
	if err != nil {
		return state, 0, fmt.Errorf("error in steering member candidate updates: %v", err)
	}

	qrnUpdates, err := types.PB2TM.QrnUpdates(abciQrnUpdates)
	if err != nil {
		return state, 0, err
	}
	if len(qrnUpdates) > 0 {
		blockExec.logger.Debug("updates to qrns", "updates", types.QrnListString(qrnUpdates))
	}

	// Update the state with the block and responses.
	state, err = updateState(state, blockID, &block.Header, abciResponses, validators, standingMemberUpdates, steeringMemberCandidateUpdates, qrnUpdates)
	if err != nil {
		return state, 0, fmt.Errorf("commit failed for application: %v", err)
	}

	//@@@logging: committed state
	// Lock mempool, commit app state, update mempoool.
	appHash, retainHeight, err := blockExec.Commit(state, block, abciResponses.DeliverTxs)
	if err != nil {
		return state, 0, fmt.Errorf("commit failed for application: %v", err)
	}

	// Update evpool with the latest state.
	blockExec.evpool.Update(state, block.Evidence.Evidence)

	fail.Fail() // XXX

	// Update the app hash and save the state.
	state.AppHash = appHash

	blockExec.logger.Error("ApplyBlock",
		"state.QrnSet.Height", state.QrnSet.Height,
		"InitialHeight", state.InitialHeight,
		"LastBlockHeight", state.LastBlockHeight)

	if err := blockExec.store.Save(state); err != nil {
		return state, 0, err
	}
	fail.Fail() // XXX

	// Events are fired after everything else.
	// NOTE: if we crash between Commit and Save, events wont be fired during replay
	fireEvents(blockExec.logger, blockExec.eventBus, block, abciResponses, validators)

	return state, retainHeight, nil
}

// Commit locks the mempool, runs the ABCI Commit message, and updates the
// mempool.
// It returns the result of calling abci.Commit (the AppHash) and the height to retain (if any).
// The Mempool must be locked during commit and update because state is
// typically reset on Commit and old txs must be replayed against committed
// state before new txs are run in the mempool, lest they be invalid.
func (blockExec *BlockExecutor) Commit(
	state State,
	block *types.Block,
	deliverTxResponses []*abci.ResponseDeliverTx,
) ([]byte, int64, error) {
	blockExec.mempool.Lock()
	defer blockExec.mempool.Unlock()

	// while mempool is Locked, flush to ensure all async requests have completed
	// in the ABCI app before Commit.
	err := blockExec.mempool.FlushAppConn()
	if err != nil {
		blockExec.logger.Error("client error during mempool.FlushAppConn", "err", err)
		return nil, 0, err
	}

	// Commit block, get hash back
	res, err := blockExec.proxyApp.CommitSync()
	if err != nil {
		blockExec.logger.Error("client error during proxyAppConn.CommitSync", "err", err)
		return nil, 0, err
	}

	// ResponseCommit has no error code - just data
	blockExec.logger.Info(
		"committed state",
		"height", block.Height,
		"num_txs", len(block.Txs),
		"app_hash", fmt.Sprintf("%X", res.Data),
	)

	// Update mempool.
	err = blockExec.mempool.Update(
		block.Height,
		block.Txs,
		deliverTxResponses,
		TxPreCheck(state),
		TxPostCheck(state),
	)

	return res.Data, res.RetainHeight, err
}

// Executes block's transactions on proxyAppConn.
// Returns a list of transaction results and updates to the validator set
func execBlockOnProxyApp(
	logger log.Logger,
	proxyAppConn proxy.AppConnConsensus,
	block *types.Block,
	store Store,
	initialHeight int64,
) (*tmstate.ABCIResponses, error) {
	var validTxs, invalidTxs = 0, 0

	txIndex := 0
	abciResponses := new(tmstate.ABCIResponses)
	dtxs := make([]*abci.ResponseDeliverTx, len(block.Txs))

	abciResponses.DeliverTxs = dtxs

	// Execute transactions and get hash.
	proxyCb := func(req *abci.Request, res *abci.Response) {
		if r, ok := res.Value.(*abci.Response_DeliverTx); ok {
			// TODO: make use of res.Log
			// TODO: make use of this info
			// Blocks may include invalid txs.
			txRes := r.DeliverTx
			if txRes.Code == abci.CodeTypeOK {
				validTxs++
			} else {
				logger.Debug("invalid tx", "code", txRes.Code, "log", txRes.Log)
				invalidTxs++
			}

			abciResponses.DeliverTxs[txIndex] = txRes
			txIndex++
		}
	}
	proxyAppConn.SetResponseCallback(proxyCb)

	commitInfo := getBeginBlockValidatorInfo(block, store, initialHeight)

	byzVals := make([]abci.Evidence, 0)
	for _, evidence := range block.Evidence.Evidence {
		byzVals = append(byzVals, evidence.ABCI()...)
	}

	// Begin block
	var err error
	pbh := block.Header.ToProto()
	if pbh == nil {
		return nil, errors.New("nil header")
	}

	// Send block information
	abciResponses.BeginBlock, err = proxyAppConn.BeginBlockSync(abci.RequestBeginBlock{
		Hash:                block.Hash(),
		Header:              *pbh,
		LastCommitInfo:      commitInfo,
		ByzantineValidators: byzVals,
	})
	if err != nil {
		logger.Error("error in proxyAppConn.BeginBlock", "err", err)
		return nil, err
	}

	// run txs of block
	for _, tx := range block.Txs {
		proxyAppConn.DeliverTxAsync(abci.RequestDeliverTx{Tx: tx})
		if err := proxyAppConn.Error(); err != nil {
			return nil, err
		}
	}

	// End block.
	abciResponses.EndBlock, err = proxyAppConn.EndBlockSync(abci.RequestEndBlock{Height: block.Height})
	if err != nil {
		logger.Error("error in proxyAppConn.EndBlock", "err", err)
		return nil, err
	}
	logger.Info("executed block", "height", block.Height, "num_valid_txs", validTxs, "num_invalid_txs", invalidTxs)
	return abciResponses, nil
}

func getBeginBlockValidatorInfo(block *types.Block, store Store,
	initialHeight int64) abci.LastCommitInfo {
	voteInfos := make([]abci.VoteInfo, block.LastCommit.Size())
	// Initial block -> LastCommitInfo.Votes are empty.
	// Remember that the first LastCommit is intentionally empty, so it makes
	// sense for LastCommitInfo.Votes to also be empty.
	if block.Height > initialHeight {
		lastValSet, err := store.LoadValidators(block.Height - 1)
		if err != nil {
			panic(err)
		}

		// Sanity check that commit size matches validator set size - only applies
		// after first block.
		var (
			commitSize = block.LastCommit.Size()
			valSetLen  = len(lastValSet.Validators)
		)
		if commitSize != valSetLen {
			panic(fmt.Sprintf(
				"commit size (%d) doesn't match valset length (%d) at height %d\n\n%v\n\n%v",
				commitSize, valSetLen, block.Height, block.LastCommit.Signatures, lastValSet.Validators,
			))
		}

		for i, val := range lastValSet.Validators {
			commitSig := block.LastCommit.Signatures[i]
			voteInfos[i] = abci.VoteInfo{
				Validator:       types.TM2PB.Validator(val),
				SignedLastBlock: !commitSig.Absent(),
			}
		}
	}

	return abci.LastCommitInfo{
		Round: block.LastCommit.Round,
		Votes: voteInfos,
	}
}

func validateValidatorUpdates(abciUpdates []abci.ValidatorUpdate,
	params tmproto.ValidatorParams) error {
	for _, valUpdate := range abciUpdates {
		if valUpdate.GetPower() < 0 {
			return fmt.Errorf("voting power can't be negative %v", valUpdate)
		} else if valUpdate.GetPower() == 0 {
			// continue, since this is deleting the validator, and thus there is no
			// pubkey to check
			continue
		}

		// Check if validator's pubkey matches an ABCI type in the consensus params
		pk, err := cryptoenc.PubKeyFromProto(valUpdate.PubKey)
		if err != nil {
			return err
		}

		if !types.IsValidPubkeyType(params, pk.Type()) {
			return fmt.Errorf("validator %v is using pubkey %s, which is unsupported for consensus",
				valUpdate, pk.Type())
		}
	}
	return nil
}

func validateStandingMemberUpdates(standingMemberUpdatesAbci []abci.StandingMemberUpdate, params tmproto.ValidatorParams) error {
	for _, valUpdate := range standingMemberUpdatesAbci {
		// Check if validator's pubkey matches an ABCI type in the consensus params
		pk, err := cryptoenc.PubKeyFromProto(valUpdate.PubKey)
		if err != nil {
			return err
		}

		if !types.IsValidPubkeyType(params, pk.Type()) {
			return fmt.Errorf("validator %v is using pubkey %s, which is unsupported for consensus",
				valUpdate, pk.Type())
		}
	}
	return nil
}

func validateSteeringMemberCandidateUpdates(steeringMemberCandidateUpdatesAbci []abci.SteeringMemberCandidateUpdate, params tmproto.ValidatorParams) error {
	for _, valUpdate := range steeringMemberCandidateUpdatesAbci {
		// Check if validator's pubkey matches an ABCI type in the consensus params
		pk, err := cryptoenc.PubKeyFromProto(valUpdate.PubKey)
		if err != nil {
			return err
		}

		if !types.IsValidPubkeyType(params, pk.Type()) {
			return fmt.Errorf("validator %v is using pubkey %s, which is unsupported for consensus",
				valUpdate, pk.Type())
		}
	}
	return nil
}

func validateQrnUpdates(qrnUpdatesAbci []abci.QrnUpdate, params tmproto.ValidatorParams) error {
	for _, qrnUpdate := range qrnUpdatesAbci {
		pk, err := cryptoenc.PubKeyFromProto(qrnUpdate.StandingMemberPubKey)
		if err != nil {
			return err
		}

		if !types.IsValidPubkeyType(params, pk.Type()) {
			return fmt.Errorf("validator %v is using pubkey %s, which is unsupported for consensus",
				qrnUpdate, pk.Type())
		}
	}
	return nil
}

// updateState returns a new State updated according to the header and responses.
func updateState(
	state State,
	blockID types.BlockID,
	header *types.Header,
	abciResponses *tmstate.ABCIResponses,
	validatorUpdates []*types.Validator,
	standingMemberUpdates []*types.StandingMember,
	steeringMemberCandidateUpdates []*types.SteeringMemberCandidate,
	qrnUpdates []*types.Qrn,
) (State, error) {
	// Copy the valset so we can apply changes from EndBlock
	// and update s.LastValidators and s.Validators.
	nValSet := state.NextValidators.Copy()
	standingMemberSet := state.StandingMemberSet.Copy()
	steeringMemberCandidateSet := state.SteeringMemberCandidateSet.Copy()

	// Update the validator set with the latest abciResponses.
	lastHeightValsChanged := state.LastHeightValidatorsChanged
	if len(validatorUpdates) > 0 {
		nValSet = types.NewValidatorSet(validatorUpdates)

		// Change results from this height but only applies to the next next height.
		lastHeightValsChanged = header.Height + 1 + 1
	}

	lastHeightStandingMembersChanged := state.LastHeightStandingMembersChanged
	if len(standingMemberUpdates) > 0 {
		err := standingMemberSet.UpdateWithChangeSet(standingMemberUpdates)
		if err != nil {
			return state, fmt.Errorf("error changing standing member set: %v", err)
		}
		lastHeightStandingMembersChanged = header.Height + 1
	}

	lastHeightSteeringMemberCandidatesChanged := state.LastHeightSteeringMemberCandidatesChanged
	if len(steeringMemberCandidateUpdates) > 0 {
		err := steeringMemberCandidateSet.UpdateWithChangeSet(steeringMemberCandidateUpdates)
		if err != nil {
			return state, fmt.Errorf("error changing steering member candidate set: %v", err)
		}
		lastHeightSteeringMemberCandidatesChanged = header.Height + 1
	}

	// Update the params with the latest abciResponses.
	nextParams := state.ConsensusParams
	lastHeightParamsChanged := state.LastHeightConsensusParamsChanged
	if abciResponses.EndBlock.ConsensusParamUpdates != nil {
		// NOTE: must not mutate s.ConsensusParams
		nextParams = types.UpdateConsensusParams(state.ConsensusParams, abciResponses.EndBlock.ConsensusParamUpdates)
		err := types.ValidateConsensusParams(nextParams)
		if err != nil {
			return state, fmt.Errorf("error updating consensus params: %v", err)
		}

		state.Version.Consensus.App = nextParams.Version.AppVersion

		// Change results from this height but only applies to the next height.
		lastHeightParamsChanged = header.Height + 1
	}

	// Update to consensus round with the latest abciResponses.
	nextConsensusRound := state.ConsensusRound
	lastHeightConsensusRoundChanged := state.LastHeightConsensusRoundChanged

	if abciResponses.EndBlock.ConsensusRoundUpdates != nil {
		nextConsensusRound = types.UpdateConsensusRound(state.ConsensusRound, abciResponses.EndBlock.ConsensusRoundUpdates)
		err := types.ValidateConsensusRound(nextConsensusRound, header.Height)
		if err != nil {
			return state, fmt.Errorf("error updating consensus round: %v", err)
		}

		// Change results from this height but only applies to the next height.
		lastHeightParamsChanged = header.Height + 1
	}

	// fmt.Println("stompesi ---------------------------------")
	// fmt.Println(header.Height)
	// fmt.Println(nextConsensusRound.ConsensusStartBlockHeight)
	// fmt.Println(nextConsensusRound.Peorid)
	// fmt.Println(nextConsensusRound.ConsensusStartBlockHeight + int64(nextConsensusRound.Peorid) + 1)
	// fmt.Println("stompesi ---------------------------------")

	if nextConsensusRound.ConsensusStartBlockHeight+int64(nextConsensusRound.Peorid)-1 == header.Height {
		i := 0
		validatorSize := len(standingMemberSet.StandingMembers)
		if state.SettingSteeringMember != nil {
			fmt.Println("stompesi---------------SteeringMemberIndexes")
			fmt.Println(len(state.SettingSteeringMember.SteeringMemberIndexes))
			validatorSize = validatorSize + len(state.SettingSteeringMember.SteeringMemberIndexes)
		} else {
			fmt.Println("stompesi---------------SteeringMemberIndexes-empty")
		}

		validators := make([]*types.Validator, validatorSize)

		if state.SettingSteeringMember != nil {
			for _, steeringMemberIndex := range state.SettingSteeringMember.SteeringMemberIndexes {
				validators[i] = types.NewValidator(state.SteeringMemberCandidateSet.SteeringMemberCandidates[steeringMemberIndex].PubKey, 10)
				i++
			}
		}

		for _, standingMember := range standingMemberSet.StandingMembers {
			validators[i] = types.NewValidator(standingMember.PubKey, 10)
			i++
		}

		state.IsSetSteeringMember = false

		nextConsensusRound.ConsensusStartBlockHeight = header.Height + 1

		nextConsensusStartBlockHeight := nextConsensusRound.ConsensusStartBlockHeight + int64(nextConsensusRound.Peorid)

		state.QrnSet = state.NextQrnSet.Copy()
		state.NextQrnSet = types.NewQrnSet(nextConsensusStartBlockHeight, standingMemberSet, nil)

		state.VrfSet = state.NextVrfSet.Copy()
		state.NextVrfSet = types.NewVrfSet(nextConsensusStartBlockHeight, steeringMemberCandidateSet, nil)

		nValSet = types.NewValidatorSet(validators)

		// Change results from this height but only applies to the next next height.
		lastHeightValsChanged = header.Height + 1 + 1
	}

	// if len(qrnUpdates) > 0 {
	// 	state.QrnSet = types.NewQrnSet(state.QrnSet.Height, state.StandingMemberSet, qrnUpdates)
	// }

	nextVersion := state.Version

	// NOTE: the AppHash has not been populated.
	// It will be filled on state.Save.
	return State{
		Version:         nextVersion,
		ChainID:         state.ChainID,
		InitialHeight:   state.InitialHeight,
		LastBlockHeight: header.Height,
		LastBlockID:     blockID,
		LastBlockTime:   header.Time,

		NextValidators:              nValSet,
		Validators:                  state.NextValidators.Copy(),
		LastValidators:              state.Validators.Copy(),
		LastHeightValidatorsChanged: lastHeightValsChanged,

		ConsensusParams:                  nextParams,
		LastHeightConsensusParamsChanged: lastHeightParamsChanged,
		LastResultsHash:                  ABCIResponsesResultsHash(abciResponses),
		AppHash:                          nil,

		StandingMemberSet:                standingMemberSet.Copy(),
		LastHeightStandingMembersChanged: lastHeightStandingMembersChanged,

		SteeringMemberCandidateSet:                steeringMemberCandidateSet.Copy(),
		LastHeightSteeringMemberCandidatesChanged: lastHeightSteeringMemberCandidatesChanged,

		ConsensusRound:                  nextConsensusRound,
		LastHeightConsensusRoundChanged: lastHeightConsensusRoundChanged,

		NextQrnSet: state.NextQrnSet.Copy(),
		QrnSet:     state.QrnSet.Copy(),

		NextVrfSet: state.NextVrfSet.Copy(),
		VrfSet:     state.VrfSet.Copy(),

		SettingSteeringMember: state.SettingSteeringMember.Copy(),
		IsSetSteeringMember:   state.IsSetSteeringMember,
	}, nil
}

// Fire NewBlock, NewBlockHeader.
// Fire TxEvent for every tx.
// NOTE: if Reapchain crashes before commit, some or all of these events may be published again.
func fireEvents(
	logger log.Logger,
	eventBus types.BlockEventPublisher,
	block *types.Block,
	abciResponses *tmstate.ABCIResponses,
	validatorUpdates []*types.Validator,
) {
	if err := eventBus.PublishEventNewBlock(types.EventDataNewBlock{
		Block:            block,
		ResultBeginBlock: *abciResponses.BeginBlock,
		ResultEndBlock:   *abciResponses.EndBlock,
	}); err != nil {
		logger.Error("failed publishing new block", "err", err)
	}

	if err := eventBus.PublishEventNewBlockHeader(types.EventDataNewBlockHeader{
		Header:           block.Header,
		NumTxs:           int64(len(block.Txs)),
		ResultBeginBlock: *abciResponses.BeginBlock,
		ResultEndBlock:   *abciResponses.EndBlock,
	}); err != nil {
		logger.Error("failed publishing new block header", "err", err)
	}

	if len(block.Evidence.Evidence) != 0 {
		for _, ev := range block.Evidence.Evidence {
			if err := eventBus.PublishEventNewEvidence(types.EventDataNewEvidence{
				Evidence: ev,
				Height:   block.Height,
			}); err != nil {
				logger.Error("failed publishing new evidence", "err", err)
			}
		}
	}

	for i, tx := range block.Data.Txs {
		if err := eventBus.PublishEventTx(types.EventDataTx{TxResult: abci.TxResult{
			Height: block.Height,
			Index:  uint32(i),
			Tx:     tx,
			Result: *(abciResponses.DeliverTxs[i]),
		}}); err != nil {
			logger.Error("failed publishing event TX", "err", err)
		}
	}

	if len(validatorUpdates) > 0 {
		if err := eventBus.PublishEventValidatorSetUpdates(
			types.EventDataValidatorSetUpdates{ValidatorUpdates: validatorUpdates}); err != nil {
			logger.Error("failed publishing event", "err", err)
		}
	}
}

//----------------------------------------------------------------------------------------------------
// Execute block without state. TODO: eliminate

// ExecCommitBlock executes and commits a block on the proxyApp without validating or mutating the state.
// It returns the application root hash (result of abci.Commit).
func ExecCommitBlock(
	appConnConsensus proxy.AppConnConsensus,
	block *types.Block,
	logger log.Logger,
	store Store,
	initialHeight int64,
) ([]byte, error) {
	_, err := execBlockOnProxyApp(logger, appConnConsensus, block, store, initialHeight)
	if err != nil {
		logger.Error("failed executing block on proxy app", "height", block.Height, "err", err)
		return nil, err
	}

	// Commit block, get hash back
	res, err := appConnConsensus.CommitSync()
	if err != nil {
		logger.Error("client error during proxyAppConn.CommitSync", "err", res)
		return nil, err
	}

	// ResponseCommit has no error or log, just data
	return res.Data, nil
}
