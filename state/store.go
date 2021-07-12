package state

import (
	"errors"
	"fmt"

	"github.com/gogo/protobuf/proto"
	dbm "github.com/tendermint/tm-db"

	abci "github.com/reapchain/reapchain-core/abci/types"
	tmmath "github.com/reapchain/reapchain-core/libs/math"
	tmos "github.com/reapchain/reapchain-core/libs/os"
	tmstate "github.com/reapchain/reapchain-core/proto/reapchain/state"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
	"github.com/reapchain/reapchain-core/types"
)

const (
	// persist validators every valSetCheckpointInterval blocks to avoid
	// LoadValidators taking too much time.
	// https://github.com/reapchain/reapchain-core/pull/3438
	// 100000 results in ~ 100ms to get 100 validators (see BenchmarkLoadValidators)
	valSetCheckpointInterval = 100000
)

//------------------------------------------------------------------------

func calcValidatorsKey(height int64) []byte {
	return []byte(fmt.Sprintf("validatorsKey:%v", height))
}

func calcConsensusParamsKey(height int64) []byte {
	return []byte(fmt.Sprintf("consensusParamsKey:%v", height))
}

func calcConsensusRoundKey(height int64) []byte {
	return []byte(fmt.Sprintf("consensusRoundKey:%v", height))
}

func calcABCIResponsesKey(height int64) []byte {
	return []byte(fmt.Sprintf("abciResponsesKey:%v", height))
}

//----------------------

//go:generate mockery --case underscore --name Store

// Store defines the state store interface
//
// It is used to retrieve current state and save and load ABCI responses,
// validators and consensus parameters
type Store interface {
	// LoadFromDBOrGenesisFile loads the most recent state.
	// If the chain is new it will use the genesis file from the provided genesis file path as the current state.
	LoadFromDBOrGenesisFile(string) (State, error)
	// LoadFromDBOrGenesisDoc loads the most recent state.
	// If the chain is new it will use the genesis doc as the current state.
	LoadFromDBOrGenesisDoc(*types.GenesisDoc) (State, error)
	// Load loads the current state of the blockchain
	Load() (State, error)
	// LoadValidators loads the validator set at a given height
	LoadValidators(int64) (*types.ValidatorSet, error)
	// LoadABCIResponses loads the abciResponse for a given height
	LoadABCIResponses(int64) (*tmstate.ABCIResponses, error)
	// LoadConsensusParams loads the consensus params for a given height
	LoadConsensusParams(int64) (tmproto.ConsensusParams, error)
	// Save overwrites the previous state with the updated one
	Save(State) error
	// SaveABCIResponses saves ABCIResponses for a given height
	SaveABCIResponses(int64, *tmstate.ABCIResponses) error
	// Bootstrap is used for bootstrapping state when not starting from a initial height.
	Bootstrap(State) error
	// PruneStates takes the height from which to start prning and which height stop at
	PruneStates(int64, int64) error

	LoadStandingMembers(int64) (*types.StandingMemberSet, error)
	LoadQrns(int64) (*types.QrnSet, error)
}

// dbStore wraps a db (github.com/tendermint/tm-db)
type dbStore struct {
	db dbm.DB
}

var _ Store = (*dbStore)(nil)

// NewStore creates the dbStore of the state pkg.
func NewStore(db dbm.DB) Store {
	return dbStore{db}
}

// LoadStateFromDBOrGenesisFile loads the most recent state from the database,
// or creates a new one from the given genesisFilePath.
func (store dbStore) LoadFromDBOrGenesisFile(genesisFilePath string) (State, error) {
	state, err := store.Load()
	if err != nil {
		return State{}, err
	}
	if state.IsEmpty() {
		var err error
		state, err = MakeGenesisStateFromFile(genesisFilePath)
		if err != nil {
			return state, err
		}
	}

	return state, nil
}

// LoadStateFromDBOrGenesisDoc loads the most recent state from the database,
// or creates a new one from the given genesisDoc.
func (store dbStore) LoadFromDBOrGenesisDoc(genesisDoc *types.GenesisDoc) (State, error) {
	//fmt.Println("stompesi-start-2")
	state, err := store.Load()
	if err != nil {
		return State{}, err
	}

	if state.IsEmpty() {
		var err error
		state, err = MakeGenesisState(genesisDoc)
		if err != nil {
			return state, err
		}
	}

	return state, nil
}

// LoadState loads the State from the database.
func (store dbStore) Load() (State, error) {
	//fmt.Println("stompesi-start-3")
	return store.loadState(stateKey)
}

func (store dbStore) loadState(key []byte) (state State, err error) {
	//fmt.Println("stompesi-start-4", key)
	buf, err := store.db.Get(key)
	if err != nil {
		//fmt.Println("stompesi-start-4-1-return")
		return state, err
	}

	if len(buf) == 0 {
		//fmt.Println("stompesi-start-4-2-return")
		return state, nil
	}

	//fmt.Println("stompesi-start-kkk")
	sp := new(tmstate.State)

	//fmt.Println("stompesi-start-start-Unmarshal", sp)
	err = proto.Unmarshal(buf, sp)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadState: Data has been corrupted or its spec has changed:
		%v\n`, err))
	}

	//fmt.Println("stompesi-start-end-Unmarshal", sp)

	sm, err := StateFromProto(sp)
	if err != nil {
		return state, err
	}

	return *sm, nil
}

// Save persists the State, the ValidatorsInfo, and the ConsensusParamsInfo to the database.
// This flushes the writes (e.g. calls SetSync).
func (store dbStore) Save(state State) error {
	return store.save(state, stateKey)
}

func (store dbStore) save(state State, key []byte) error {
	//fmt.Println("stompesi-init-save")

	nextHeight := state.LastBlockHeight + 1
	// If first block, save validators for the block.
	if nextHeight == 1 {
		//fmt.Println("stompesi-init-save-1")
		nextHeight = state.InitialHeight
		// This extra logic due to Reapchain validator set changes being delayed 1 block.
		// It may get overwritten due to InitChain validator updates.
		if err := store.saveValidatorsInfo(nextHeight, nextHeight, state.Validators); err != nil {
			return err
		}
	}

	//fmt.Println("stompesi-init-save-3")
	if err := store.saveStandingMembersInfo(nextHeight, state.LastHeightConsensusParamsChanged, state.StandingMembers); err != nil {
		return err
	}

	if err := store.saveQrnsInfo(nextHeight, state.LastHeightConsensusParamsChanged, state.Qrns); err != nil {
		return err
	}

	//fmt.Println("2stompesi-asdfdf", state.ConsensusRoundInfo)

	if err := store.saveConsensusRoundInfo(nextHeight, state.LastHeightConsensusParamsChanged, state.ConsensusRoundInfo.ToProto()); err != nil {
		return err
	}

	// Save next validators.
	if err := store.saveValidatorsInfo(nextHeight+1, state.LastHeightValidatorsChanged, state.NextValidators); err != nil {
		return err
	}

	// Save next consensus params.
	if err := store.saveConsensusParamsInfo(nextHeight,
		state.LastHeightConsensusParamsChanged, state.ConsensusParams); err != nil {
		return err
	}
	//fmt.Println("stompesi-init-save-2", state.StandingMembers)

	err := store.db.SetSync(key, state.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// BootstrapState saves a new state, used e.g. by state sync when starting from non-zero height.
func (store dbStore) Bootstrap(state State) error {
	height := state.LastBlockHeight + 1
	if height == 1 {
		height = state.InitialHeight
	}

	if height > 1 && !state.LastValidators.IsNilOrEmpty() {
		if err := store.saveValidatorsInfo(height-1, height-1, state.LastValidators); err != nil {
			return err
		}
	}

	if err := store.saveValidatorsInfo(height, height, state.Validators); err != nil {
		return err
	}

	if err := store.saveStandingMembersInfo(height, height, state.StandingMembers); err != nil {
		return err
	}

	if err := store.saveConsensusRoundInfo(height, state.LastHeightConsensusParamsChanged, state.ConsensusRoundInfo.ToProto()); err != nil {
		return err
	}

	if err := store.saveValidatorsInfo(height+1, height+1, state.NextValidators); err != nil {
		return err
	}

	if err := store.saveConsensusParamsInfo(height,
		state.LastHeightConsensusParamsChanged, state.ConsensusParams); err != nil {
		return err
	}

	//fmt.Println("stompesi-init-SetSync", state)
	return store.db.SetSync(stateKey, state.Bytes())
}

// PruneStates deletes states between the given heights (including from, excluding to). It is not
// guaranteed to delete all states, since the last checkpointed state and states being pointed to by
// e.g. `LastHeightChanged` must remain. The state at to must also exist.
//
// The from parameter is necessary since we can't do a key scan in a performant way due to the key
// encoding not preserving ordering: https://github.com/reapchain/reapchain-core/issues/4567
// This will cause some old states to be left behind when doing incremental partial prunes,
// specifically older checkpoints and LastHeightChanged targets.
func (store dbStore) PruneStates(from int64, to int64) error {
	if from <= 0 || to <= 0 {
		return fmt.Errorf("from height %v and to height %v must be greater than 0", from, to)
	}
	if from >= to {
		return fmt.Errorf("from height %v must be lower than to height %v", from, to)
	}
	valInfo, err := loadValidatorsInfo(store.db, to)
	if err != nil {
		return fmt.Errorf("validators at height %v not found: %w", to, err)
	}
	paramsInfo, err := store.loadConsensusParamsInfo(to)
	if err != nil {
		return fmt.Errorf("consensus params at height %v not found: %w", to, err)
	}

	smInfo, err := loadStandingMembersInfo(store.db, to)
	if err != nil {
		return fmt.Errorf("validators at height %v not found: %w", to, err)
	}

	qrnInfo, err := loadQrnsInfo(store.db, to)
	if err != nil {
		return fmt.Errorf("validators at height %v not found: %w", to, err)
	}

	keepQrns := make(map[int64]bool)
	if qrnInfo.QrnSet == nil {
		keepQrns[qrnInfo.LastHeightChanged] = true
		keepQrns[lastStoredHeightFor(to, qrnInfo.LastHeightChanged)] = true // qrn last checkpoint too
	}

	keepVals := make(map[int64]bool)
	if valInfo.ValidatorSet == nil {
		keepVals[valInfo.LastHeightChanged] = true
		keepVals[lastStoredHeightFor(to, valInfo.LastHeightChanged)] = true // keep last checkpoint too
	}

	keepSms := make(map[int64]bool)
	if smInfo.StandingMemberSet == nil {
		keepSms[smInfo.LastHeightChanged] = true
		keepSms[lastStoredHeightFor(to, smInfo.LastHeightChanged)] = true // keep last checkpoint too
	}

	keepParams := make(map[int64]bool)
	if paramsInfo.ConsensusParams.Equal(&tmproto.ConsensusParams{}) {
		keepParams[paramsInfo.LastHeightChanged] = true
	}

	batch := store.db.NewBatch()
	defer batch.Close()
	pruned := uint64(0)

	// We have to delete in reverse order, to avoid deleting previous heights that have validator
	// sets and consensus params that we may need to retrieve.
	for h := to - 1; h >= from; h-- {
		// For heights we keep, we must make sure they have the full validator set or consensus
		// params, otherwise they will panic if they're retrieved directly (instead of
		// indirectly via a LastHeightChanged pointer).
		if keepVals[h] {
			v, err := loadValidatorsInfo(store.db, h)
			if err != nil || v.ValidatorSet == nil {
				vip, err := store.LoadValidators(h)
				if err != nil {
					return err
				}

				pvi, err := vip.ToProto()
				if err != nil {
					return err
				}

				v.ValidatorSet = pvi
				v.LastHeightChanged = h

				bz, err := v.Marshal()
				if err != nil {
					return err
				}
				err = batch.Set(calcValidatorsKey(h), bz)
				if err != nil {
					return err
				}
			}
		} else {
			err = batch.Delete(calcValidatorsKey(h))
			if err != nil {
				return err
			}
		}

		if keepSms[h] {
			v, err := loadStandingMembersInfo(store.db, h)
			if err != nil || v.StandingMemberSet == nil {
				vip, err := store.LoadStandingMembers(h)
				if err != nil {
					return err
				}

				pvi, err := vip.ToProto()
				if err != nil {
					return err
				}

				v.StandingMemberSet = pvi
				v.LastHeightChanged = h

				bz, err := v.Marshal()
				if err != nil {
					return err
				}
				err = batch.Set(calcStandingMembersKey(h), bz)
				if err != nil {
					return err
				}
			}
		} else {
			err = batch.Delete(calcStandingMembersKey(h))
			if err != nil {
				return err
			}
		}

		if keepQrns[h] {
			v, err := loadQrnsInfo(store.db, h)
			if err != nil || v.QrnSet == nil {
				vip, err := store.LoadQrns(h)
				if err != nil {
					return err
				}

				pvi, err := vip.ToProto()
				if err != nil {
					return err
				}

				v.QrnSet = pvi
				v.LastHeightChanged = h

				bz, err := v.Marshal()
				if err != nil {
					return err
				}
				err = batch.Set(calcQrnsKey(h), bz)
				if err != nil {
					return err
				}
			}
		} else {
			err = batch.Delete(calcQrnsKey(h))
			if err != nil {
				return err
			}
		}

		if keepParams[h] {
			p, err := store.loadConsensusParamsInfo(h)
			if err != nil {
				return err
			}

			if p.ConsensusParams.Equal(&tmproto.ConsensusParams{}) {
				p.ConsensusParams, err = store.LoadConsensusParams(h)
				if err != nil {
					return err
				}

				p.LastHeightChanged = h
				bz, err := p.Marshal()
				if err != nil {
					return err
				}

				err = batch.Set(calcConsensusParamsKey(h), bz)
				if err != nil {
					return err
				}
			}
		} else {
			err = batch.Delete(calcConsensusParamsKey(h))
			if err != nil {
				return err
			}
		}

		err = batch.Delete(calcABCIResponsesKey(h))
		if err != nil {
			return err
		}
		pruned++

		// avoid batches growing too large by flushing to database regularly
		if pruned%1000 == 0 && pruned > 0 {
			err := batch.Write()
			if err != nil {
				return err
			}
			batch.Close()
			batch = store.db.NewBatch()
			defer batch.Close()
		}
	}

	err = batch.WriteSync()
	if err != nil {
		return err
	}

	return nil
}

//------------------------------------------------------------------------

// ABCIResponsesResultsHash returns the root hash of a Merkle tree of
// ResponseDeliverTx responses (see ABCIResults.Hash)
//
// See merkle.SimpleHashFromByteSlices
func ABCIResponsesResultsHash(ar *tmstate.ABCIResponses) []byte {
	return types.NewResults(ar.DeliverTxs).Hash()
}

// LoadABCIResponses loads the ABCIResponses for the given height from the
// database. If not found, ErrNoABCIResponsesForHeight is returned.
//
// This is useful for recovering from crashes where we called app.Commit and
// before we called s.Save(). It can also be used to produce Merkle proofs of
// the result of txs.
func (store dbStore) LoadABCIResponses(height int64) (*tmstate.ABCIResponses, error) {
	buf, err := store.db.Get(calcABCIResponsesKey(height))
	if err != nil {
		return nil, err
	}
	if len(buf) == 0 {

		return nil, ErrNoABCIResponsesForHeight{height}
	}

	abciResponses := new(tmstate.ABCIResponses)
	err = abciResponses.Unmarshal(buf)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadABCIResponses: Data has been corrupted or its spec has
                changed: %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return abciResponses, nil
}

// SaveABCIResponses persists the ABCIResponses to the database.
// This is useful in case we crash after app.Commit and before s.Save().
// Responses are indexed by height so they can also be loaded later to produce
// Merkle proofs.
//
// Exposed for testing.
func (store dbStore) SaveABCIResponses(height int64, abciResponses *tmstate.ABCIResponses) error {
	var dtxs []*abci.ResponseDeliverTx
	// strip nil values,
	for _, tx := range abciResponses.DeliverTxs {
		if tx != nil {
			dtxs = append(dtxs, tx)
		}
	}
	abciResponses.DeliverTxs = dtxs

	bz, err := abciResponses.Marshal()
	if err != nil {
		return err
	}

	//fmt.Println("stompesi-init-SetSync", bz)
	err = store.db.SetSync(calcABCIResponsesKey(height), bz)
	if err != nil {
		return err
	}

	return nil
}

//-----------------------------------------------------------------------------

// LoadValidators loads the ValidatorSet for a given height.
// Returns ErrNoValSetForHeight if the validator set can't be found for this height.
func (store dbStore) LoadValidators(height int64) (*types.ValidatorSet, error) {
	valInfo, err := loadValidatorsInfo(store.db, height)
	if err != nil {
		return nil, ErrNoValSetForHeight{height}
	}
	if valInfo.ValidatorSet == nil {
		lastStoredHeight := lastStoredHeightFor(height, valInfo.LastHeightChanged)
		valInfo2, err := loadValidatorsInfo(store.db, lastStoredHeight)
		if err != nil || valInfo2.ValidatorSet == nil {
			return nil,
				fmt.Errorf("couldn't find validators at height %d (height %d was originally requested): %w",
					lastStoredHeight,
					height,
					err,
				)
		}

		vs, err := types.ValidatorSetFromProto(valInfo2.ValidatorSet)
		if err != nil {
			return nil, err
		}

		vs.IncrementProposerPriority(tmmath.SafeConvertInt32(height - lastStoredHeight)) // mutate
		vi2, err := vs.ToProto()
		if err != nil {
			return nil, err
		}

		valInfo2.ValidatorSet = vi2
		valInfo = valInfo2
	}

	vip, err := types.ValidatorSetFromProto(valInfo.ValidatorSet)
	if err != nil {
		return nil, err
	}

	return vip, nil
}

func lastStoredHeightFor(height, lastHeightChanged int64) int64 {
	checkpointHeight := height - height%valSetCheckpointInterval
	return tmmath.MaxInt64(checkpointHeight, lastHeightChanged)
}

// CONTRACT: Returned ValidatorsInfo can be mutated.
func loadValidatorsInfo(db dbm.DB, height int64) (*tmstate.ValidatorsInfo, error) {
	buf, err := db.Get(calcValidatorsKey(height))
	if err != nil {
		return nil, err
	}

	if len(buf) == 0 {
		return nil, errors.New("value retrieved from db is empty")
	}

	v := new(tmstate.ValidatorsInfo)
	err = v.Unmarshal(buf)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadValidators: Data has been corrupted or its spec has changed:
                %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return v, nil
}

// saveValidatorsInfo persists the validator set.
//
// `height` is the effective height for which the validator is responsible for
// signing. It should be called from s.Save(), right before the state itself is
// persisted.
func (store dbStore) saveValidatorsInfo(height, lastHeightChanged int64, valSet *types.ValidatorSet) error {
	if lastHeightChanged > height {
		return errors.New("lastHeightChanged cannot be greater than ValidatorsInfo height")
	}
	valInfo := &tmstate.ValidatorsInfo{
		LastHeightChanged: lastHeightChanged,
	}
	// Only persist validator set if it was updated or checkpoint height (see
	// valSetCheckpointInterval) is reached.
	if height == lastHeightChanged || height%valSetCheckpointInterval == 0 {
		pv, err := valSet.ToProto()
		if err != nil {
			return err
		}
		valInfo.ValidatorSet = pv
	}

	bz, err := valInfo.Marshal()
	if err != nil {
		return err
	}

	//fmt.Println("stompesi-init-Set-calcValidatorsKey", valInfo)
	err = store.db.Set(calcValidatorsKey(height), bz)
	if err != nil {
		return err
	}

	return nil
}

func (store dbStore) saveStandingMembersInfo(height, lastHeightChanged int64, valSet *types.StandingMemberSet) error {
	//fmt.Println("stompesi-jkjkjkdjfkjdkf")

	if lastHeightChanged > height {
		return errors.New("lastHeightChanged cannot be greater than StandingMembersInfo height")
	}
	smInfo := &tmstate.StandingMembersInfo{
		LastHeightChanged: lastHeightChanged,
	}
	// Only persist validator set if it was updated or checkpoint height (see
	// valSetCheckpointInterval) is reached.
	if height == lastHeightChanged || height%valSetCheckpointInterval == 0 {
		pv, err := valSet.ToProto()
		if err != nil {
			return err
		}
		smInfo.StandingMemberSet = pv
	}

	bz, err := smInfo.Marshal()
	if err != nil {
		return err
	}

	//fmt.Println("stompesi-saveStandingMembersInfo", smInfo)
	err = store.db.Set(calcStandingMembersKey(height), bz)
	if err != nil {
		return err
	}

	return nil
}

func (store dbStore) saveQrnsInfo(height, lastHeightChanged int64, valSet *types.QrnSet) error {
	//fmt.Println("stompesi-jkjkjkdjfkjdkf")

	if lastHeightChanged > height {
		return errors.New("lastHeightChanged cannot be greater than QrnsInfo height")
	}
	smInfo := &tmstate.QrnsInfo{
		LastHeightChanged: lastHeightChanged,
	}
	// Only persist validator set if it was updated or checkpoint height (see
	// valSetCheckpointInterval) is reached.
	if height == lastHeightChanged || height%valSetCheckpointInterval == 0 {
		pv, err := valSet.ToProto()
		if err != nil {
			return err
		}
		smInfo.QrnSet = pv
	}

	bz, err := smInfo.Marshal()
	if err != nil {
		return err
	}

	//fmt.Println("stompesi-saveQrnsInfo", smInfo)
	err = store.db.Set(calcQrnsKey(height), bz)
	if err != nil {
		return err
	}

	return nil
}

//-----------------------------------------------------------------------------

// ConsensusParamsInfo represents the latest consensus params, or the last height it changed

// LoadConsensusParams loads the ConsensusParams for a given height.
func (store dbStore) LoadConsensusParams(height int64) (tmproto.ConsensusParams, error) {
	empty := tmproto.ConsensusParams{}

	paramsInfo, err := store.loadConsensusParamsInfo(height)
	if err != nil {
		return empty, fmt.Errorf("could not find consensus params for height #%d: %w", height, err)
	}

	if paramsInfo.ConsensusParams.Equal(&empty) {
		paramsInfo2, err := store.loadConsensusParamsInfo(paramsInfo.LastHeightChanged)
		if err != nil {
			return empty, fmt.Errorf(
				"couldn't find consensus params at height %d as last changed from height %d: %w",
				paramsInfo.LastHeightChanged,
				height,
				err,
			)
		}

		paramsInfo = paramsInfo2
	}

	return paramsInfo.ConsensusParams, nil
}

func (store dbStore) loadConsensusParamsInfo(height int64) (*tmstate.ConsensusParamsInfo, error) {
	buf, err := store.db.Get(calcConsensusParamsKey(height))
	if err != nil {
		return nil, err
	}
	if len(buf) == 0 {
		return nil, errors.New("value retrieved from db is empty")
	}

	paramsInfo := new(tmstate.ConsensusParamsInfo)
	if err = paramsInfo.Unmarshal(buf); err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadConsensusParams: Data has been corrupted or its spec has changed:
                %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return paramsInfo, nil
}

// saveConsensusParamsInfo persists the consensus params for the next block to disk.
// It should be called from s.Save(), right before the state itself is persisted.
// If the consensus params did not change after processing the latest block,
// only the last height for which they changed is persisted.
func (store dbStore) saveConsensusParamsInfo(nextHeight, changeHeight int64, params tmproto.ConsensusParams) error {
	paramsInfo := &tmstate.ConsensusParamsInfo{
		LastHeightChanged: changeHeight,
	}

	if changeHeight == nextHeight {
		paramsInfo.ConsensusParams = params
	}
	bz, err := paramsInfo.Marshal()
	if err != nil {
		return err
	}

	err = store.db.Set(calcConsensusParamsKey(nextHeight), bz)
	if err != nil {
		return err
	}

	return nil
}

func (store dbStore) saveConsensusRoundInfo(nextHeight, changeHeight int64, params tmproto.ConsensusRound) error {
	paramsInfo := &tmstate.ConsensusRoundInfo{
		LastHeightChanged: changeHeight,
	}

	if changeHeight == nextHeight {
		paramsInfo.ConsensusRoundInfo = &params
	}
	bz, err := paramsInfo.Marshal()
	if err != nil {
		return err
	}

	err = store.db.Set(calcConsensusRoundKey(nextHeight), bz)
	if err != nil {
		return err
	}

	return nil
}

func (store dbStore) LoadStandingMembers(height int64) (*types.StandingMemberSet, error) {
	smInfo, err := loadStandingMembersInfo(store.db, height)
	if err != nil {
		return nil, ErrNoValSetForHeight{height}
	}
	if smInfo.StandingMemberSet == nil {
		lastStoredHeight := lastStoredHeightFor(height, smInfo.LastHeightChanged)
		smInfo2, err := loadStandingMembersInfo(store.db, lastStoredHeight)
		if err != nil || smInfo2.StandingMemberSet == nil {
			return nil,
				fmt.Errorf("couldn't find standing members at height %d (height %d was originally requested): %w",
					lastStoredHeight,
					height,
					err,
				)
		}

		vs, err := types.StandingMemberSetFromProto(smInfo2.StandingMemberSet)
		if err != nil {
			return nil, err
		}

		vi2, err := vs.ToProto()
		if err != nil {
			return nil, err
		}

		smInfo2.StandingMemberSet = vi2
		smInfo = smInfo2
	}

	vip, err := types.StandingMemberSetFromProto(smInfo.StandingMemberSet)
	if err != nil {
		return nil, err
	}

	return vip, nil
}

// CONTRACT: Returned StandingMemberInfo can be mutated.
func loadStandingMembersInfo(db dbm.DB, height int64) (*tmstate.StandingMembersInfo, error) {
	buf, err := db.Get(calcStandingMembersKey(height))
	if err != nil {
		return nil, err
	}

	if len(buf) == 0 {
		return nil, errors.New("value retrieved from db is empty")
	}

	v := new(tmstate.StandingMembersInfo)
	err = v.Unmarshal(buf)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadStandingMembers: Data has been corrupted or its spec has changed:
                %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return v, nil
}

func calcStandingMembersKey(height int64) []byte {
	return []byte(fmt.Sprintf("standingMembersKey:%v", height))
}

func (store dbStore) LoadQrns(height int64) (*types.QrnSet, error) {
	smInfo, err := loadQrnsInfo(store.db, height)
	if err != nil {
		return nil, ErrNoValSetForHeight{height}
	}
	if smInfo.QrnSet == nil {
		lastStoredHeight := lastStoredHeightFor(height, smInfo.LastHeightChanged)
		smInfo2, err := loadQrnsInfo(store.db, lastStoredHeight)
		if err != nil || smInfo2.QrnSet == nil {
			return nil,
				fmt.Errorf("couldn't find qrns at height %d (height %d was originally requested): %w",
					lastStoredHeight,
					height,
					err,
				)
		}

		vs, err := types.QrnSetFromProto(smInfo2.QrnSet)
		if err != nil {
			return nil, err
		}

		vi2, err := vs.ToProto()
		if err != nil {
			return nil, err
		}

		smInfo2.QrnSet = vi2
		smInfo = smInfo2
	}

	vip, err := types.QrnSetFromProto(smInfo.QrnSet)
	if err != nil {
		return nil, err
	}

	return vip, nil
}

// CONTRACT: Returned QrnInfo can be mutated.
func loadQrnsInfo(db dbm.DB, height int64) (*tmstate.QrnsInfo, error) {
	buf, err := db.Get(calcQrnsKey(height))
	if err != nil {
		return nil, err
	}

	if len(buf) == 0 {
		return nil, errors.New("value retrieved from db is empty")
	}

	v := new(tmstate.QrnsInfo)
	err = v.Unmarshal(buf)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadQrns: Data has been corrupted or its spec has changed:
                %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return v, nil
}

func calcQrnsKey(height int64) []byte {
	return []byte(fmt.Sprintf("qrnsKey:%v", height))
}
