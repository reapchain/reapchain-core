package state

import (
	"errors"
	"fmt"

	"github.com/gogo/protobuf/proto"
	dbm "github.com/tendermint/tm-db"

	abci "github.com/reapchain/reapchain-core/abci/types"
	cfg "github.com/reapchain/reapchain-core/config"
	tmmath "github.com/reapchain/reapchain-core/libs/math"
	tmos "github.com/reapchain/reapchain-core/libs/os"
	tmstate "github.com/reapchain/reapchain-core/proto/podc/state"
	tmproto "github.com/reapchain/reapchain-core/proto/podc/types"
	"github.com/reapchain/reapchain-core/types"
)

const (
	// persist validators every valSetCheckpointInterval blocks to avoid
	// LoadValidators taking too much time.
	// https://github.com/reapchain/reapchain-core/pull/3438
	// 100000 results in ~ 100ms to get 100 validators (see BenchmarkLoadValidators)
	valSetCheckpointInterval = 100000
)

// ------------------------------------------------------------------------
func calcConsensusRoundKey(height int64) []byte {
	return []byte(fmt.Sprintf("consensusRoundKey:%v", height))
}

func calcQrnsKey(height int64) []byte {
	return []byte(fmt.Sprintf("qrnsKey:%v", height))
}

func calcNextQrnsKey(height int64) []byte {
	return []byte(fmt.Sprintf("nextQrnsKey:%v", height))
}

func calcVrfsKey(height int64) []byte {
	return []byte(fmt.Sprintf("vrfsKey:%v", height))
}

func calcNextVrfsKey(height int64) []byte {
	return []byte(fmt.Sprintf("nextVrfsKey:%v", height))
}

func calcSettingSteeringMemberKey(height int64) []byte {
	return []byte(fmt.Sprintf("settingSteeringMemberKey:%v", height))
}

func calcStandingMembersKey(height int64) []byte {
	return []byte(fmt.Sprintf("standingMembersKey:%v", height))
}

func calcSteeringMemberCandidatesKey(height int64) []byte {
	return []byte(fmt.Sprintf("steeringMemberCandidatesKey:%v", height))
}

func calcValidatorsKey(height int64) []byte {
	return []byte(fmt.Sprintf("validatorsKey:%v", height))
}

func calcConsensusParamsKey(height int64) []byte {
	return []byte(fmt.Sprintf("consensusParamsKey:%v", height))
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
	// Load loads the rollback state of the blockchain
	LoadRollback() (State, error)
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
	// Close closes the connection with the database
	Close() error

	LoadConsensusRound(int64) (tmproto.ConsensusRound, error)
	LoadStandingMemberSet(int64) (*types.StandingMemberSet, error)
	LoadSteeringMemberCandidateSet(int64) (*types.SteeringMemberCandidateSet, error)
	LoadQrnSet(int64) (*types.QrnSet, error)
	LoadVrfSet(int64) (*types.VrfSet, error)

	LoadNextQrnSet(int64) (*types.QrnSet, error)
	LoadNextVrfSet(int64) (*types.VrfSet, error)
	LoadSettingSteeringMember(height int64) (*types.SettingSteeringMember, error)
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

func ExportState(config *cfg.Config, height int64) (*State, error) {
	dbType := dbm.BackendType(config.DBBackend)
	stateDB, err := dbm.NewDB("state", dbType, config.DBDir())
	if err != nil {
		tmos.Exit(err.Error())
	}
	
	stateStore := NewStore(stateDB)

	nextVrfSet, err := stateStore.LoadNextVrfSet(height)
	if err != nil {
		return nil, fmt.Errorf("couldn't read nextVrfSet: %w", err)
	}


	nextQrnSet, err := stateStore.LoadNextQrnSet(height)
	if err != nil {
		return nil, fmt.Errorf("couldn't read nextQrnSet: %w", err)
	}

	vrfSet, err := stateStore.LoadVrfSet(height)
	if err != nil {
		return nil, fmt.Errorf("couldn't read vrfSet: %w", err)
	}

	qrnSet, err := stateStore.LoadQrnSet(height)
	if err != nil {
		return nil, fmt.Errorf("couldn't read qrnSet: %w", err)
	}

	consensusRound, err := stateStore.LoadConsensusRound(height)
	if err != nil {
		return nil, fmt.Errorf("couldn't read consensusRound: %w", err)
	}

	standingMemberSet, err := stateStore.LoadStandingMemberSet(height)
	if err != nil {
		return nil, fmt.Errorf("couldn't read standingMemberSet: %w", err)
	}

	steeringMemberCandidateSet, err := stateStore.LoadSteeringMemberCandidateSet(height)
	if err != nil {
		return nil, fmt.Errorf("couldn't read steeringMemberCandidateSet: %w", err)
	}

	return &State{
		NextVrfSet: nextVrfSet,
		NextQrnSet: nextQrnSet,
		VrfSet: vrfSet,
		QrnSet: qrnSet,
		ConsensusRound: consensusRound,
		StandingMemberSet: standingMemberSet,
		SteeringMemberCandidateSet: steeringMemberCandidateSet,
	}, nil
}

// LoadState loads the State from the database.
func (store dbStore) Load() (State, error) {
	return store.loadState(stateKey)
}

func (store dbStore) loadState(key []byte) (state State, err error) {
	buf, err := store.db.Get(key)
	if err != nil {
		return state, err
	}

	if len(buf) == 0 {
		return state, nil
	}

	sp := new(tmstate.State)

	err = proto.Unmarshal(buf, sp)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadState: Data has been corrupted or its spec has changed:
		%v\n`, err))
	}

	sm, err := StateFromProto(sp)
	if err != nil {
		return state, err
	}

	return *sm, nil
}

// LoadState loads the State from the database.
func (store dbStore) LoadRollback() (State, error) {
	return store.loadState(rollbackStateKey)
}

func (store dbStore) loadRollbackState(key []byte) (state State, err error) {
	buf, err := store.db.Get(key)
	if err != nil {
		return state, err
	}

	if len(buf) == 0 {
		return state, nil
	}

	sp := new(tmstate.State)

	err = proto.Unmarshal(buf, sp)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadState: Data has been corrupted or its spec has changed:
		%v\n`, err))
	}

	sm, err := StateFromProto(sp)
	if err != nil {
		return state, err
	}

	return *sm, nil
}

// Save persists the State, the ValidatorsInfo, and the ConsensusParamsInfo to the database.
// This flushes the writes (e.g. calls SetSync).
func (store dbStore) Save(state State) error {
	if previousConsensusRound == nil {
		previousConsensusRound = &state.ConsensusRound
	}
	
	if previousConsensusRound.ConsensusStartBlockHeight + int64(previousConsensusRound.Period) - 1 == state.LastBlockHeight {
		previousConsensusRound = &state.ConsensusRound
		store.save(state, rollbackStateKey)
	}
	return store.save(state, stateKey)
}

func (store dbStore) save(state State, key []byte) error {
	nextHeight := state.LastBlockHeight + 1
	// If first block, save validators for the block.
	if nextHeight == 1 {
		nextHeight = state.InitialHeight
		// This extra logic due to Reapchain validator set changes being delayed 1 block.
		// It may get overwritten due to InitChain validator updates.

		if err := store.saveValidatorsInfo(nextHeight, nextHeight, state.Validators); err != nil {
			return err
		}

		if err := store.saveQrnsInfo(nextHeight, nextHeight, state.QrnSet); err != nil {
			return err
		}

		if err := store.saveNextQrnsInfo(nextHeight, nextHeight, state.NextQrnSet); err != nil {
			return err
		}

		if err := store.saveVrfsInfo(nextHeight, nextHeight, state.VrfSet); err != nil {
			return err
		}

		if err := store.saveNextVrfsInfo(nextHeight, nextHeight, state.NextVrfSet); err != nil {
			return err
		}

		if err := store.saveStandingMembersInfo(nextHeight, nextHeight, state.StandingMemberSet); err != nil {
			return err
		}

		if err := store.saveSteeringMemberCandidatesInfo(nextHeight, nextHeight, state.SteeringMemberCandidateSet); err != nil {
			return err
		}
	}
	if err := store.saveConsensusRoundInfo(nextHeight, state.LastHeightConsensusRoundChanged, state.ConsensusRound); err != nil {
		return err
	}

	if err := store.saveQrnsInfo(nextHeight, state.LastHeightQrnChanged, state.QrnSet); err != nil {
		return err
	}

	if err := store.saveNextQrnsInfo(nextHeight, state.LastHeightNextQrnChanged, state.NextQrnSet); err != nil {
		return err
	}

	if err := store.saveVrfsInfo(nextHeight, state.LastHeightVrfChanged, state.VrfSet); err != nil {
		return err
	}

	if err := store.saveNextVrfsInfo(nextHeight, state.LastHeightNextVrfChanged, state.NextVrfSet); err != nil {
		return err
	}

	if err := store.saveSettingSteeringMemberInfo(nextHeight, state.LastHeightSettingSteeringMemberChanged, state.SettingSteeringMember); err != nil {
		return err
	}

	if err := store.saveStandingMembersInfo(nextHeight, state.LastHeightStandingMembersChanged, state.StandingMemberSet); err != nil {
		return err
	}

	if err := store.saveSteeringMemberCandidatesInfo(nextHeight, state.LastHeightSteeringMemberCandidatesChanged, state.SteeringMemberCandidateSet); err != nil {
		return err
	}

	// Save next validators.
	if err := store.saveValidatorsInfo(nextHeight+1, state.LastHeightValidatorsChanged, state.NextValidators); err != nil {
		return err
	}

	// Save next consensus params.
	if err := store.saveConsensusParamsInfo(nextHeight, state.LastHeightConsensusParamsChanged, state.ConsensusParams); err != nil {
		return err
	}

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

	if err := store.saveConsensusRoundInfo(height, height, state.ConsensusRound); err != nil {
		return err
	}

	if err := store.saveQrnsInfo(height, height, state.QrnSet); err != nil {
		return err
	}

	if err := store.saveNextQrnsInfo(height, height, state.NextQrnSet); err != nil {
		return err
	}

	if err := store.saveVrfsInfo(height, height, state.VrfSet); err != nil {
		return err
	}

	if err := store.saveNextVrfsInfo(height, height, state.NextVrfSet); err != nil {
		return err
	}

	if err := store.saveSettingSteeringMemberInfo(height, height, state.SettingSteeringMember); err != nil {
		return err
	}

	if err := store.saveStandingMembersInfo(height, height, state.StandingMemberSet); err != nil {
		return err
	}

	if err := store.saveSteeringMemberCandidatesInfo(height, height, state.SteeringMemberCandidateSet); err != nil {
		return err
	}

	if err := store.saveValidatorsInfo(height, height, state.Validators); err != nil {
		return err
	}

	if err := store.saveValidatorsInfo(height+1, height+1, state.NextValidators); err != nil {
		return err
	}

	if err := store.saveConsensusParamsInfo(height, state.LastHeightConsensusParamsChanged, state.ConsensusParams); err != nil {
		return err
	}

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

	smInfo, err := loadStandingMemberSetInfo(store.db, to)
	if err != nil {
		return fmt.Errorf("validators at height %v not found: %w", to, err)
	}

	qrnsInfo, err := loadQrnSetInfo(store.db, to)
	if err != nil {
		return fmt.Errorf("qrns at height %v not found: %w", to, err)
	}

	SettingSteeringMemberInfo, err := loadSettingSteeringMemberInfo(store.db, to)
	if err != nil {
		return fmt.Errorf("qrns at height %v not found: %w", to, err)
	}

	keepSettingSteeringMember := make(map[int64]bool)
	if SettingSteeringMemberInfo.SettingSteeringMember == nil {
		keepSettingSteeringMember[SettingSteeringMemberInfo.LastHeightChanged] = true
		keepSettingSteeringMember[tmmath.MaxInt64(to, SettingSteeringMemberInfo.LastHeightChanged)] = true
	}

	keepVals := make(map[int64]bool)
	if valInfo.ValidatorSet == nil {
		keepVals[valInfo.LastHeightChanged] = true
		keepVals[lastStoredHeightFor(to, valInfo.LastHeightChanged)] = true // keep last checkpoint too
	}

	keepParams := make(map[int64]bool)
	if paramsInfo.ConsensusParams.Equal(&tmproto.ConsensusParams{}) {
		keepParams[paramsInfo.LastHeightChanged] = true
	}

	keepSms := make(map[int64]bool)
	if smInfo.StandingMemberSet == nil {
		keepSms[smInfo.LastHeightChanged] = true
		keepSms[lastStoredHeightFor(to, smInfo.LastHeightChanged)] = true // keep last checkpoint too
	}

	keepQrnSet := make(map[int64]bool)
	if qrnsInfo.QrnSet == nil {
		keepQrnSet[qrnsInfo.LastHeightChanged] = true
		keepQrnSet[tmmath.MaxInt64(to, qrnsInfo.LastHeightChanged)] = true // keep last checkpoint too
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

	err = store.db.Set(calcValidatorsKey(height), bz)
	if err != nil {
		return err
	}

	return nil
}

func (store dbStore) saveQrnsInfo(height, lastHeightChanged int64, qrnSet *types.QrnSet) error {
	if lastHeightChanged > height {
		return errors.New("lastHeightChanged cannot be greater than QrnsInfo height")
	}

	qrnSetInfo := &tmstate.QrnsInfo{
		LastHeightChanged: lastHeightChanged,
	}

	if height == lastHeightChanged {
		qrnSetProto, err := qrnSet.ToProto()
		if err != nil {
			return err
		}
		qrnSetInfo.QrnSet = qrnSetProto
	}

	bz, err := qrnSetInfo.Marshal()
	if err != nil {
		return err
	}

	err = store.db.Set(calcQrnsKey(height), bz)
	if err != nil {
		return err
	}

	return nil
}

func (store dbStore) saveNextQrnsInfo(height, lastHeightChanged int64, qrnSet *types.QrnSet) error {
	if lastHeightChanged > height {
		return errors.New("lastHeightChanged cannot be greater than QrnsInfo height")
	}

	qrnSetInfo := &tmstate.QrnsInfo{
		LastHeightChanged: lastHeightChanged,
	}
	
	if height == lastHeightChanged {
		qrnSetProto, err := qrnSet.ToProto()
		if err != nil {
			return err
		}
		qrnSetInfo.QrnSet = qrnSetProto
	}

	bz, err := qrnSetInfo.Marshal()
	if err != nil {
		return err
	}

	err = store.db.Set(calcNextQrnsKey(height), bz)
	if err != nil {
		return err
	}

	return nil
}

func (store dbStore) saveVrfsInfo(height, lastHeightChanged int64, vrfSet *types.VrfSet) error {
	if lastHeightChanged > height {
		return errors.New("lastHeightChanged cannot be greater than VrfsInfo height")
	}

	vrfSetInfo := &tmstate.VrfsInfo{
		LastHeightChanged: lastHeightChanged,
	}

	if height == lastHeightChanged {
		vrfSetProto, err := vrfSet.ToProto()
		if err != nil {
			return err
		}
		vrfSetInfo.VrfSet = vrfSetProto
	}

	bz, err := vrfSetInfo.Marshal()
	if err != nil {
		return err
	}

	err = store.db.Set(calcVrfsKey(height), bz)
	if err != nil {
		return err
	}

	return nil
}

func (store dbStore) saveNextVrfsInfo(height, lastHeightChanged int64, vrfSet *types.VrfSet) error {
	if lastHeightChanged > height {
		return errors.New("lastHeightChanged cannot be greater than VrfsInfo height")
	}

	vrfSetInfo := &tmstate.VrfsInfo{
		LastHeightChanged: lastHeightChanged,
	}

	if height == lastHeightChanged {
		vrfSetProto, err := vrfSet.ToProto()
		if err != nil {
			return err
		}
		vrfSetInfo.VrfSet = vrfSetProto
	}

	bz, err := vrfSetInfo.Marshal()
	if err != nil {
		return err
	}

	err = store.db.Set(calcNextVrfsKey(height), bz)
	if err != nil {
		return err
	}

	return nil
}

func (store dbStore) saveSettingSteeringMemberInfo(height, lastHeightChanged int64, settingSteeringMember *types.SettingSteeringMember) error {
	if lastHeightChanged > height {
		return errors.New("lastHeightChanged cannot be greater than SettingSteeringMemberInfo height")
	}

	SettingSteeringMemberInfo := &tmstate.SettingSteeringMemberInfo{
		LastHeightChanged: height,
	}

	if height == lastHeightChanged {
		SettingSteeringMemberInfo.SettingSteeringMember = settingSteeringMember.ToProto()
	}

	bz, err := SettingSteeringMemberInfo.Marshal()
	if err != nil {
		return err
	}

	err = store.db.Set(calcSettingSteeringMemberKey(height), bz)
	if err != nil {
		return err
	}

	return nil
}

func (store dbStore) saveStandingMembersInfo(height, lastHeightChanged int64, standingMemberSet *types.StandingMemberSet) error {
	smInfo := &tmstate.StandingMembersInfo{
		LastHeightChanged:         lastHeightChanged,
		CurrentCoordinatorRanking: standingMemberSet.CurrentCoordinatorRanking,
	}
	if height == lastHeightChanged || height%valSetCheckpointInterval == 0 {
		standingMemberSetProto, err := standingMemberSet.ToProto()
		if err != nil {
			return err
		}
		smInfo.StandingMemberSet = standingMemberSetProto
	}

	bz, err := smInfo.Marshal()
	if err != nil {
		return err
	}

	err = store.db.Set(calcStandingMembersKey(height), bz)
	if err != nil {
		return err
	}

	return nil
}

func (store dbStore) saveSteeringMemberCandidatesInfo(height, lastHeightChanged int64, steeringMemberCandidateSet *types.SteeringMemberCandidateSet) error {
	smInfo := &tmstate.SteeringMemberCandidateSetInfo{
		LastHeightChanged: lastHeightChanged,
	}

	if height == lastHeightChanged || height%valSetCheckpointInterval == 0 {
		steeringMemberCandidateSetProto, err := steeringMemberCandidateSet.ToProto()
		if err != nil {
			return err
		}
		smInfo.SteeringMemberCandidateSet = steeringMemberCandidateSetProto
	}

	bz, err := smInfo.Marshal()
	if err != nil {
		return err
	}

	err = store.db.Set(calcSteeringMemberCandidatesKey(height), bz)
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

func (store dbStore) loadConsensusRoundInfo(height int64) (*tmstate.ConsensusRoundInfo, error) {
	buf, err := store.db.Get(calcConsensusRoundKey(height))
	if err != nil {
		return nil, err
	}
	if len(buf) == 0 {
		return nil, errors.New("value retrieved from db is empty")
	}

	paramsInfo := new(tmstate.ConsensusRoundInfo)
	if err = paramsInfo.Unmarshal(buf); err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadConsensusRound: Data has been corrupted or its spec has changed:
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

func (store dbStore) saveConsensusRoundInfo(nextHeight, changeHeight int64, consensusRoundProto tmproto.ConsensusRound) error {
	paramsInfo := &tmstate.ConsensusRoundInfo{
		LastHeightChanged: changeHeight,
	}

	if changeHeight == nextHeight {
		paramsInfo.ConsensusRound = consensusRoundProto
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

func (store dbStore) LoadStandingMemberSet(height int64) (*types.StandingMemberSet, error) {
	standingMemberSetInfo, err := loadStandingMemberSetInfo(store.db, height)
	currentCoordinatorRanking := standingMemberSetInfo.CurrentCoordinatorRanking

	if err != nil {
		return nil, ErrNoStandingMemberSetForHeight{height}
	}

	if standingMemberSetInfo.StandingMemberSet == nil {
		lastStoredHeight := lastStoredHeightFor(height, standingMemberSetInfo.LastHeightChanged)

		standingMemberSetInfo2, err := loadStandingMemberSetInfo(store.db, lastStoredHeight)
		if err != nil || standingMemberSetInfo2.StandingMemberSet == nil {
			return nil,
				fmt.Errorf("couldn't find standing members at height %d (height %d was originally requested): %w",
					lastStoredHeight,
					height,
					err,
				)
		}
		standingMemberSetInfo = standingMemberSetInfo2
	}

	standingMemberSet, err := types.StandingMemberSetFromProto(standingMemberSetInfo.StandingMemberSet)
	if err != nil {
		return nil, err
	}
	standingMemberSet.CurrentCoordinatorRanking = currentCoordinatorRanking

	return standingMemberSet, nil
}

func (store dbStore) LoadSteeringMemberCandidateSet(height int64) (*types.SteeringMemberCandidateSet, error) {
	steeringMemberCandidateSetInfo, err := loadSteeringMemberCandidateSetInfo(store.db, height)
	if err != nil {
		return nil, ErrNoSteeringMemberSetForHeight{height}
	}
	if steeringMemberCandidateSetInfo.SteeringMemberCandidateSet == nil {
		lastStoredHeight := lastStoredHeightFor(height, steeringMemberCandidateSetInfo.LastHeightChanged)
		steeringMemberCandidateSetInfo2, err := loadSteeringMemberCandidateSetInfo(store.db, lastStoredHeight)
		if err != nil || steeringMemberCandidateSetInfo2.SteeringMemberCandidateSet == nil {
			return nil,
				fmt.Errorf("couldn't find steering member candidates at height %d (height %d was originally requested): %w",
					lastStoredHeight,
					height,
					err,
				)
		}
		steeringMemberCandidateSetInfo = steeringMemberCandidateSetInfo2
	}

	steeringMemberCandidateSet, err := types.SteeringMemberCandidateSetFromProto(steeringMemberCandidateSetInfo.SteeringMemberCandidateSet)
	if err != nil {
		return nil, err
	}

	return steeringMemberCandidateSet, nil
}

func (store dbStore) LoadQrnSet(height int64) (*types.QrnSet, error) {
	qrnSetInfo, err := loadQrnSetInfo(store.db, height)
	if err != nil {
		return nil, ErrNoQrnSetForHeight{height}
	}

	if qrnSetInfo.QrnSet == nil {
		lastStoredHeight := qrnSetInfo.LastHeightChanged
		qrnSetInfo2, err := loadQrnSetInfo(store.db, lastStoredHeight)
		if err != nil || qrnSetInfo2.QrnSet == nil {
			return nil,
				fmt.Errorf("couldn't find qrns at height %d (height %d was originally requested): %w",
					lastStoredHeight,
					height,
					err,
				)
		}
		qrnSetInfo = qrnSetInfo2
	}

	qrnSet, err := types.QrnSetFromProto(qrnSetInfo.QrnSet)
	if err != nil {
		return nil, err
	}

	return qrnSet, nil
}

func (store dbStore) LoadNextQrnSet(height int64) (*types.QrnSet, error) {
	qrnSetInfo, err := loadNextQrnSetInfo(store.db, height)
	if err != nil {
		return nil, ErrNoQrnSetForHeight{height}
	}

	if qrnSetInfo.QrnSet == nil {
		lastStoredHeight := qrnSetInfo.LastHeightChanged
		qrnSetInfo2, err := loadNextQrnSetInfo(store.db, lastStoredHeight)
		if err != nil || qrnSetInfo2.QrnSet == nil {
			return nil,
				fmt.Errorf("couldn't find qrns at height %d (height %d was originally requested): %w",
					lastStoredHeight,
					height,
					err,
				)
		}
		qrnSetInfo = qrnSetInfo2
	}

	qrnSet, err := types.QrnSetFromProto(qrnSetInfo.QrnSet)
	if err != nil {
		return nil, err
	}

	return qrnSet, nil
}

func (store dbStore) LoadSettingSteeringMember(height int64) (*types.SettingSteeringMember, error) {
	SettingSteeringMemberInfo, err := loadSettingSteeringMemberInfo(store.db, height)
	if err != nil {
		return nil, ErrNoSettingSteeringMemberForHeight{height}
	}

	settingSteeringMember := types.SettingSteeringMemberFromProto(SettingSteeringMemberInfo.SettingSteeringMember)

	return settingSteeringMember, nil
}

func (store dbStore) LoadVrfSet(height int64) (*types.VrfSet, error) {
	vrfSetInfo, err := loadVrfSetInfo(store.db, height)
	if err != nil {
		return nil, ErrNoVrfSetForHeight{height}
	}

	if vrfSetInfo.VrfSet == nil {
		lastStoredHeight := vrfSetInfo.LastHeightChanged
		vrfSetInfo2, err := loadVrfSetInfo(store.db, lastStoredHeight)
		if err != nil || vrfSetInfo2.VrfSet == nil {
			return nil,
				fmt.Errorf("couldn't find vrfs at height %d (height %d was originally requested): %w",
					lastStoredHeight,
					height,
					err,
				)
		}
		vrfSetInfo = vrfSetInfo2
	}

	vrfSet, err := types.VrfSetFromProto(vrfSetInfo.VrfSet)
	if err != nil {
		return nil, err
	}

	return vrfSet, nil
}

func (store dbStore) LoadNextVrfSet(height int64) (*types.VrfSet, error) {
	vrfSetInfo, err := loadNextVrfSetInfo(store.db, height)
	if err != nil {
		return nil, ErrNoVrfSetForHeight{height}
	}

	if vrfSetInfo.VrfSet == nil {
		lastStoredHeight := vrfSetInfo.LastHeightChanged
		vrfSetInfo2, err := loadNextVrfSetInfo(store.db, lastStoredHeight)
		if err != nil || vrfSetInfo2.VrfSet == nil {
			return nil,
				fmt.Errorf("couldn't find vrfs at height %d (height %d was originally requested): %w",
					lastStoredHeight,
					height,
					err,
				)
		}
		vrfSetInfo = vrfSetInfo2
	}

	vrfSet, err := types.VrfSetFromProto(vrfSetInfo.VrfSet)
	if err != nil {
		return nil, err
	}

	return vrfSet, nil
}

func (store dbStore) LoadConsensusRound(height int64) (tmproto.ConsensusRound, error) {
	empty := tmproto.ConsensusRound{}

	consensusRoundInfo, err := loadConsensusRoundInfo(store.db, height)
	if err != nil {
		return empty, ErrNoConsensusRoundForHeight{height}
	}

	if consensusRoundInfo.ConsensusRound.Equal(&empty) {
		lastStoredHeight := consensusRoundInfo.LastHeightChanged
		consensusRoundInfo2, err := loadConsensusRoundInfo(store.db, lastStoredHeight)
		if err != nil {
			return empty,
				fmt.Errorf("couldn't find consensus round at height %d (height %d was originally requested): %w",
					lastStoredHeight,
					height,
					err,
				)
		}
		consensusRoundInfo = consensusRoundInfo2
	}

	return consensusRoundInfo.ConsensusRound, nil
}

// CONTRACT: Returned StandingMemberInfo can be mutated.
func loadStandingMemberSetInfo(db dbm.DB, height int64) (*tmstate.StandingMembersInfo, error) {
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
		tmos.Exit(fmt.Sprintf(`LoadStandingMemberSet: Data has been corrupted or its spec has changed: %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return v, nil
}

func loadSteeringMemberCandidateSetInfo(db dbm.DB, height int64) (*tmstate.SteeringMemberCandidateSetInfo, error) {
	buf, err := db.Get(calcSteeringMemberCandidatesKey(height))
	if err != nil {
		return nil, err
	}

	if len(buf) == 0 {
		return nil, errors.New("value retrieved from db is empty")
	}

	v := new(tmstate.SteeringMemberCandidateSetInfo)
	err = v.Unmarshal(buf)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadSteeringMemberCandidateSet: Data has been corrupted or its spec has changed: %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return v, nil
}

func loadQrnSetInfo(db dbm.DB, height int64) (*tmstate.QrnsInfo, error) {
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
		tmos.Exit(fmt.Sprintf(`LoadQrnSet: Data has been corrupted or its spec has changed: %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return v, nil
}

func loadNextQrnSetInfo(db dbm.DB, height int64) (*tmstate.QrnsInfo, error) {
	buf, err := db.Get(calcNextQrnsKey(height))
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
		tmos.Exit(fmt.Sprintf(`LoadQrnSet: Data has been corrupted or its spec has changed: %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return v, nil
}

func loadSettingSteeringMemberInfo(db dbm.DB, height int64) (*tmstate.SettingSteeringMemberInfo, error) {
	buf, err := db.Get(calcSettingSteeringMemberKey(height))
	if err != nil {
		return nil, err
	}

	if len(buf) == 0 {
		return nil, errors.New("value retrieved from db is empty")
	}

	v := new(tmstate.SettingSteeringMemberInfo)
	err = v.Unmarshal(buf)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadSettingSteeringMemberInfo: Data has been corrupted or its spec has changed: %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return v, nil
}

func loadVrfSetInfo(db dbm.DB, height int64) (*tmstate.VrfsInfo, error) {
	buf, err := db.Get(calcVrfsKey(height))
	if err != nil {
		return nil, err
	}

	if len(buf) == 0 {
		return nil, errors.New("value retrieved from db is empty")
	}

	v := new(tmstate.VrfsInfo)
	err = v.Unmarshal(buf)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadVrfSet: Data has been corrupted or its spec has changed: %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return v, nil
}

func loadNextVrfSetInfo(db dbm.DB, height int64) (*tmstate.VrfsInfo, error) {
	buf, err := db.Get(calcNextVrfsKey(height))
	if err != nil {
		return nil, err
	}

	if len(buf) == 0 {
		return nil, errors.New("value retrieved from db is empty")
	}

	v := new(tmstate.VrfsInfo)
	err = v.Unmarshal(buf)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadVrfSet: Data has been corrupted or its spec has changed: %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return v, nil
}

func loadConsensusRoundInfo(db dbm.DB, height int64) (*tmstate.ConsensusRoundInfo, error) {
	buf, err := db.Get(calcConsensusRoundKey(height))
	if err != nil {
		return nil, err
	}

	if len(buf) == 0 {
		return nil, errors.New("value retrieved from db is empty")
	}

	consensusRound := new(tmstate.ConsensusRoundInfo)
	err = consensusRound.Unmarshal(buf)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadConsensusRoundInfo: Data has been corrupted or its spec has changed: %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return consensusRound, nil
}

func (store dbStore) Close() error {
	return store.db.Close()
}
