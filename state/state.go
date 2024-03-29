package state

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gogo/protobuf/proto"

	tmstate "github.com/reapchain/reapchain-core/proto/podc/state"
	tmproto "github.com/reapchain/reapchain-core/proto/podc/types"
	tmversion "github.com/reapchain/reapchain-core/proto/podc/version"
	"github.com/reapchain/reapchain-core/types"
	tmtime "github.com/reapchain/reapchain-core/types/time"
	"github.com/reapchain/reapchain-core/version"
)

// database keys
var (
	stateKey = []byte("stateKey")
)

// rollback state key
var (
	previousConsensusRound *tmproto.ConsensusRound = nil
	rollbackStateKey = []byte("rollbackStateKey")
)

//-----------------------------------------------------------------------------

// InitStateVersion sets the Consensus.Block and Software versions,
// but leaves the Consensus.App version blank.
// The Consensus.App version will be set during the Handshake, once
// we hear from the app what protocol version it is running.
var InitStateVersion = tmstate.Version{
	Consensus: tmversion.Consensus{
		Block: version.BlockProtocol,
		App:   0,
	},
	Software: version.TMCoreSemVer,
}

//-----------------------------------------------------------------------------

// State is a short description of the latest committed block of the ReapchainCore consensus.
// It keeps all information necessary to validate new blocks,
// including the last validator set and the consensus params.
// All fields are exposed so the struct can be easily serialized,
// but none of them should be mutated directly.
// Instead, use state.Copy() or state.NextState(...).
// NOTE: not goroutine-safe.
type State struct {
	Version tmstate.Version

	// immutable
	ChainID       string
	InitialHeight int64 // should be 1, not 0, when starting from height 1

	// LastBlockHeight=0 at genesis (ie. block(H=0) does not exist)
	LastBlockHeight int64
	LastBlockID     types.BlockID
	LastBlockTime   time.Time

	// LastValidators is used to validate block.LastCommit.
	// Validators are persisted to the database separately every time they change,
	// so we can query for historical validator sets.
	// Note that if s.LastBlockHeight causes a valset change,
	// we set s.LastHeightValidatorsChanged = s.LastBlockHeight + 1 + 1
	// Extra +1 due to nextValSet delay.
	NextValidators              *types.ValidatorSet
	Validators                  *types.ValidatorSet
	LastValidators              *types.ValidatorSet
	LastHeightValidatorsChanged int64

	// Consensus parameters used for validating blocks.
	// Changes returned by EndBlock and updated after Commit.
	ConsensusParams                  tmproto.ConsensusParams
	LastHeightConsensusParamsChanged int64

	// Merkle root of the results from executing prev block
	LastResultsHash []byte

	// the latest AppHash we've received from calling abci.Commit()
	AppHash []byte

	StandingMemberSet                *types.StandingMemberSet
	LastHeightStandingMembersChanged int64

	SteeringMemberCandidateSet                *types.SteeringMemberCandidateSet
	LastHeightSteeringMemberCandidatesChanged int64

	ConsensusRound                  tmproto.ConsensusRound
	LastHeightConsensusRoundChanged int64

	QrnSet     *types.QrnSet
	LastHeightQrnChanged int64

	NextQrnSet *types.QrnSet
	LastHeightNextQrnChanged int64

	VrfSet     *types.VrfSet
	LastHeightVrfChanged int64

	NextVrfSet *types.VrfSet
	LastHeightNextVrfChanged int64

	SettingSteeringMember *types.SettingSteeringMember
	LastHeightSettingSteeringMemberChanged int64
	
	IsSetSteeringMember   bool
}

// Copy makes a copy of the State for mutating.
func (state State) Copy() State {

	return State{
		Version:       state.Version,
		ChainID:       state.ChainID,
		InitialHeight: state.InitialHeight,

		LastBlockHeight: state.LastBlockHeight,
		LastBlockID:     state.LastBlockID,
		LastBlockTime:   state.LastBlockTime,

		NextValidators:              state.NextValidators.Copy(),
		Validators:                  state.Validators.Copy(),
		LastValidators:              state.LastValidators.Copy(),
		LastHeightValidatorsChanged: state.LastHeightValidatorsChanged,

		ConsensusParams:                  state.ConsensusParams,
		LastHeightConsensusParamsChanged: state.LastHeightConsensusParamsChanged,

		LastResultsHash: state.LastResultsHash,
		AppHash: state.AppHash,

		StandingMemberSet:                state.StandingMemberSet.Copy(),
		LastHeightStandingMembersChanged: state.LastHeightStandingMembersChanged,
		SteeringMemberCandidateSet:                state.SteeringMemberCandidateSet.Copy(),
		LastHeightSteeringMemberCandidatesChanged: state.LastHeightSteeringMemberCandidatesChanged,
		ConsensusRound:                  state.ConsensusRound,
		LastHeightConsensusRoundChanged: state.LastHeightConsensusRoundChanged,
		
		QrnSet:     state.QrnSet.Copy(),
		LastHeightQrnChanged: state.LastHeightQrnChanged,

		NextQrnSet: state.NextQrnSet.Copy(),
		LastHeightNextQrnChanged: state.LastHeightNextQrnChanged,

		VrfSet:     state.VrfSet.Copy(),
		LastHeightVrfChanged: state.LastHeightVrfChanged,
		
		NextVrfSet: state.NextVrfSet.Copy(),
		LastHeightNextVrfChanged: state.LastHeightNextVrfChanged,
		
		SettingSteeringMember: state.SettingSteeringMember.Copy(),
		LastHeightSettingSteeringMemberChanged: state.LastHeightSettingSteeringMemberChanged,
		
		IsSetSteeringMember:   state.IsSetSteeringMember,
	}
}

// Equals returns true if the States are identical.
func (state State) Equals(state2 State) bool {
	sbz, s2bz := state.Bytes(), state2.Bytes()
	return bytes.Equal(sbz, s2bz)
}

// Bytes serializes the State using protobuf.
// It panics if either casting to protobuf or serialization fails.
func (state State) Bytes() []byte {
	sm, err := state.ToProto()
	if err != nil {
		panic(err)
	}
	bz, err := proto.Marshal(sm)
	if err != nil {
		panic(err)
	}
	return bz
}

// IsEmpty returns true if the State is equal to the empty State.
func (state State) IsEmpty() bool {
	return state.Validators == nil // XXX can't compare to Empty
}

// ToProto takes the local state type and returns the equivalent proto type
func (state *State) ToProto() (*tmstate.State, error) {
	if state == nil {
		return nil, errors.New("state is nil")
	}

	sm := new(tmstate.State)

	sm.Version = state.Version
	sm.ChainID = state.ChainID
	sm.InitialHeight = state.InitialHeight
	sm.LastBlockHeight = state.LastBlockHeight

	sm.LastBlockID = state.LastBlockID.ToProto()
	sm.LastBlockTime = state.LastBlockTime
	vals, err := state.Validators.ToProto()
	if err != nil {
		return nil, err
	}
	sm.Validators = vals

	nVals, err := state.NextValidators.ToProto()
	if err != nil {
		return nil, err
	}
	sm.NextValidators = nVals

	if state.LastBlockHeight >= 1 { // At Block 1 LastValidators is nil
		lVals, err := state.LastValidators.ToProto()
		if err != nil {
			return nil, err
		}
		sm.LastValidators = lVals
	}

	sm.LastHeightValidatorsChanged = state.LastHeightValidatorsChanged
	sm.ConsensusParams = state.ConsensusParams
	sm.LastHeightConsensusParamsChanged = state.LastHeightConsensusParamsChanged
	sm.LastResultsHash = state.LastResultsHash
	sm.AppHash = state.AppHash

	standingMemberSetProto, err := state.StandingMemberSet.ToProto()
	if err != nil {
		return nil, err
	}
	sm.StandingMemberSet = standingMemberSetProto
	sm.LastHeightStandingMembersChanged = state.LastHeightStandingMembersChanged

	steeringMemberCandidateProto, err := state.SteeringMemberCandidateSet.ToProto()
	if err != nil {
		return nil, err
	}
	sm.SteeringMemberCandidateSet = steeringMemberCandidateProto
	sm.LastHeightSteeringMemberCandidatesChanged = state.LastHeightSteeringMemberCandidatesChanged

	sm.ConsensusRound = state.ConsensusRound
	sm.LastHeightConsensusRoundChanged = state.LastHeightConsensusRoundChanged

	qrnSetProto, err := state.QrnSet.ToProto()
	if err != nil {
		return nil, err
	}
	sm.QrnSet = qrnSetProto
	sm.LastHeightQrnChanged = state.LastHeightQrnChanged

	nextQrnSetProto, err := state.NextQrnSet.ToProto()
	if err != nil {
		return nil, err
	}
	sm.NextQrnSet = nextQrnSetProto
	sm.LastHeightNextQrnChanged = state.LastHeightNextQrnChanged

	vrfSetProto, err := state.VrfSet.ToProto()
	if err != nil {
		return nil, err
	}
	sm.VrfSet = vrfSetProto
	sm.LastHeightVrfChanged = state.LastHeightVrfChanged

	nextVrfSetProto, err := state.NextVrfSet.ToProto()
	if err != nil {
		return nil, err
	}
	sm.NextVrfSet = nextVrfSetProto
	sm.LastHeightNextVrfChanged = state.LastHeightNextVrfChanged

	sm.SettingSteeringMember = state.SettingSteeringMember.ToProto()
	sm.LastHeightSettingSteeringMemberChanged = state.LastHeightSettingSteeringMemberChanged
	
	sm.IsSetSteeringMember = state.IsSetSteeringMember

	return sm, nil
}

// StateFromProto takes a state proto message & returns the local state type
func StateFromProto(pb *tmstate.State) (*State, error) { //nolint:golint
	if pb == nil {
		return nil, errors.New("nil State")
	}

	state := new(State)

	state.Version = pb.Version
	state.ChainID = pb.ChainID
	state.InitialHeight = pb.InitialHeight

	bi, err := types.BlockIDFromProto(&pb.LastBlockID)
	if err != nil {
		return nil, err
	}
	state.LastBlockID = *bi
	state.LastBlockHeight = pb.LastBlockHeight
	state.LastBlockTime = pb.LastBlockTime

	vals, err := types.ValidatorSetFromProto(pb.Validators)
	if err != nil {
		return nil, err
	}
	state.Validators = vals

	nVals, err := types.ValidatorSetFromProto(pb.NextValidators)
	if err != nil {
		return nil, err
	}
	state.NextValidators = nVals

	if state.LastBlockHeight >= 1 { // At Block 1 LastValidators is nil
		lVals, err := types.ValidatorSetFromProto(pb.LastValidators)
		if err != nil {
			return nil, err
		}
		state.LastValidators = lVals
	} else {
		state.LastValidators = types.NewValidatorSet(nil)
	}

	state.LastHeightValidatorsChanged = pb.LastHeightValidatorsChanged
	state.ConsensusParams = pb.ConsensusParams
	state.LastHeightConsensusParamsChanged = pb.LastHeightConsensusParamsChanged
	state.LastResultsHash = pb.LastResultsHash
	state.AppHash = pb.AppHash

	standingMemberSet, err := types.StandingMemberSetFromProto(pb.StandingMemberSet)
	if err != nil {
		return nil, err
	}
	state.StandingMemberSet = standingMemberSet
	state.LastHeightStandingMembersChanged = pb.LastHeightStandingMembersChanged

	steeringMemberCandidateSet, err := types.SteeringMemberCandidateSetFromProto(pb.SteeringMemberCandidateSet)
	if err != nil {
		return nil, err
	}
	state.SteeringMemberCandidateSet = steeringMemberCandidateSet
	state.LastHeightSteeringMemberCandidatesChanged = pb.LastHeightSteeringMemberCandidatesChanged

	state.ConsensusRound = pb.ConsensusRound
	state.LastHeightConsensusRoundChanged = pb.LastHeightConsensusRoundChanged

	qrnSet, err := types.QrnSetFromProto(pb.QrnSet)
	if err != nil {
		return nil, err
	}
	state.QrnSet = qrnSet
	state.LastHeightQrnChanged = pb.LastHeightQrnChanged

	nextQrnSet, err := types.QrnSetFromProto(pb.NextQrnSet)
	if err != nil {
		return nil, err
	}
	state.NextQrnSet = nextQrnSet
	state.LastHeightNextQrnChanged = pb.LastHeightNextQrnChanged

	vrfSet, err := types.VrfSetFromProto(pb.VrfSet)
	if err != nil {
		return nil, err
	}
	state.VrfSet = vrfSet
	state.LastHeightVrfChanged = pb.LastHeightVrfChanged

	nextVrfSet, err := types.VrfSetFromProto(pb.NextVrfSet)
	if err != nil {
		return nil, err
	}
	state.NextVrfSet = nextVrfSet
	state.LastHeightNextVrfChanged = pb.LastHeightNextVrfChanged

	state.SettingSteeringMember = types.SettingSteeringMemberFromProto(pb.SettingSteeringMember)
	state.LastHeightSettingSteeringMemberChanged = pb.LastHeightSettingSteeringMemberChanged

	state.IsSetSteeringMember = pb.IsSetSteeringMember

	return state, nil
}

func SyncStateFromProto(pb *tmstate.State) (*State, error) { //nolint:golint
	if pb == nil {
		return nil, errors.New("nil State")
	}

	state := new(State)
	state.LastBlockHeight = pb.LastBlockHeight

	if pb.NextQrnSet != nil {
		nextQrnSet, err := types.QrnSetFromProto(pb.NextQrnSet)
		if err != nil {
			return nil, err
		}
		state.NextQrnSet = nextQrnSet
	}
	state.LastHeightNextQrnChanged = pb.LastHeightNextQrnChanged

	if pb.NextVrfSet != nil {
		nextVrfSet, err := types.VrfSetFromProto(pb.NextVrfSet)
		if err != nil {
			return nil, err
		}
		state.NextVrfSet = nextVrfSet
	}
	state.LastHeightNextVrfChanged = pb.LastHeightNextVrfChanged

	state.SettingSteeringMember = types.SettingSteeringMemberFromProto(pb.SettingSteeringMember)
	state.LastHeightSettingSteeringMemberChanged = pb.LastHeightSettingSteeringMemberChanged

	return state, nil
}

//------------------------------------------------------------------------
// Create a block from the latest state

// MakeBlock builds a block from the current state with the given txs, commit,
// and evidence. Note it also takes a proposerAddress because the state does not
// track rounds, and hence does not know the correct proposer. TODO: fix this!
func (state State) MakeBlock(
	height int64,
	txs []types.Tx,
	commit *types.Commit,
	evidence []types.Evidence,
	proposerAddress []byte,
) (*types.Block, *types.PartSet) {
	// Build base block with block data.
	block := types.MakeBlock(height, txs, commit, evidence)

	// Set time.
	var timestamp time.Time
	if height == state.InitialHeight {
		timestamp = state.LastBlockTime // genesis time
	} else {
		timestamp = MedianTime(commit, state.LastValidators)
	}

	consensusRound, _ := types.ConsensusRoundFromProto(state.ConsensusRound)

	// Fill rest of header with state data.
	block.Header.Populate(
		state.Version.Consensus, state.ChainID,
		timestamp, state.LastBlockID,
		state.Validators.Hash(), 
		state.NextValidators.Hash(),
		types.HashConsensusParams(state.ConsensusParams), state.AppHash, state.LastResultsHash,
		proposerAddress,
		state.StandingMemberSet.Hash(),
		state.SteeringMemberCandidateSet.Hash(),
		consensusRound,
		state.QrnSet.Hash(),
		state.VrfSet.Hash(),
	)

	return block, block.MakePartSet(types.BlockPartSizeBytes)
}

// MedianTime computes a median time for a given Commit (based on Timestamp field of votes messages) and the
// corresponding validator set. The computed time is always between timestamps of
// the votes sent by honest processes, i.e., a faulty processes can not arbitrarily increase or decrease the
// computed value.
func MedianTime(commit *types.Commit, validators *types.ValidatorSet) time.Time {
	weightedTimes := make([]*tmtime.WeightedTime, len(commit.Signatures))
	totalVotingPower := int64(0)

	for i, commitSig := range commit.Signatures {
		if commitSig.Absent() {
			continue
		}
		_, validator := validators.GetByAddress(commitSig.ValidatorAddress)
		// If there's no condition, TestValidateBlockCommit panics; not needed normally.
		if validator != nil {
			totalVotingPower += validator.VotingPower
			weightedTimes[i] = tmtime.NewWeightedTime(commitSig.Timestamp, validator.VotingPower)
		}
	}

	return tmtime.WeightedMedian(weightedTimes, totalVotingPower)
}

//------------------------------------------------------------------------
// Genesis

// MakeGenesisStateFromFile reads and unmarshals state from the given
// file.
//
// Used during replay and in tests.
func MakeGenesisStateFromFile(genDocFile string) (State, error) {
	genDoc, err := MakeGenesisDocFromFile(genDocFile)
	if err != nil {
		return State{}, err
	}
	return MakeGenesisState(genDoc)
}

// MakeGenesisDocFromFile reads and unmarshals genesis doc from the given file.
func MakeGenesisDocFromFile(genDocFile string) (*types.GenesisDoc, error) {
	genDocJSON, err := ioutil.ReadFile(genDocFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't read GenesisDoc file: %v", err)
	}
	genDoc, err := types.GenesisDocFromJSON(genDocJSON)
	if err != nil {
		return nil, fmt.Errorf("error reading GenesisDoc: %v", err)
	}
	return genDoc, nil
}

// MakeGenesisState creates state from types.GenesisDoc.
func MakeGenesisState(genDoc *types.GenesisDoc) (State, error) {
	err := genDoc.ValidateAndComplete()
	if err != nil {
		return State{}, fmt.Errorf("error in genesis file: %v", err)
	}

	var validatorSet, nextValidatorSet *types.ValidatorSet
	if genDoc.Validators == nil {
		validatorSet = types.NewValidatorSet(nil)
		nextValidatorSet = types.NewValidatorSet(nil)
	} else {
		validators := make([]*types.Validator, len(genDoc.Validators))
		for i, val := range genDoc.Validators {
			validators[i] = types.NewValidator(val.PubKey, val.Power, val.Type)
		}
		validatorSet = types.NewValidatorSet(validators)
		nextValidatorSet = types.NewValidatorSet(validators).Copy()
	}

	var standingMemberSet *types.StandingMemberSet
	if genDoc.StandingMembers == nil {
		standingMemberSet = types.NewStandingMemberSet(nil)
	} else {
		standingMembers := make([]*types.StandingMember, len(genDoc.StandingMembers))
		for i, standingMember := range genDoc.StandingMembers {
			standingMembers[i] = types.NewStandingMember(standingMember.PubKey, standingMember.Power)
		}
		standingMemberSet = types.NewStandingMemberSet(standingMembers)
	}

	var steeringMemberCandidateSet *types.SteeringMemberCandidateSet
	if genDoc.SteeringMemberCandidates == nil {
		steeringMemberCandidateSet = types.NewSteeringMemberCandidateSet(nil)
	} else {
		steeringMemberCandidates := make([]*types.SteeringMemberCandidate, len(genDoc.SteeringMemberCandidates))
		for i, steeringMemberCandidate := range genDoc.SteeringMemberCandidates {
			steeringMemberCandidates[i] = types.NewSteeringMemberCandidate(steeringMemberCandidate.PubKey, steeringMemberCandidate.Power)
		}
		steeringMemberCandidateSet = types.NewSteeringMemberCandidateSet(steeringMemberCandidates)
	}

	var qrnSet, nextQrnSet *types.QrnSet
	if genDoc.Qrns == nil {
		qrnSet = types.NewQrnSet(genDoc.InitialHeight, standingMemberSet, nil)
	} else {
		qrns := make([]*types.Qrn, len(genDoc.Qrns))
		for i, qrn := range genDoc.Qrns {
			qrns[i] = qrn.Copy()
		}
		qrnSet = types.NewQrnSet(genDoc.InitialHeight, standingMemberSet, qrns)
	}

	if genDoc.NextQrns == nil {
		nextQrnSet = types.NewQrnSet(genDoc.InitialHeight, standingMemberSet, nil)
	} else {
		nextQrns := make([]*types.Qrn, len(genDoc.NextQrns))
		for i, nextQrn := range genDoc.NextQrns {
			nextQrns[i] = nextQrn.Copy()
		}
		nextQrnSet = types.NewQrnSet(genDoc.InitialHeight, standingMemberSet, nextQrns)
	}

	var vrfSet, nextVrfSet *types.VrfSet
	if genDoc.Vrfs == nil {
		vrfSet = types.NewVrfSet(genDoc.InitialHeight, steeringMemberCandidateSet, nil)
	} else {
		vrfs := make([]*types.Vrf, len(genDoc.Vrfs))
		for i, vrf := range genDoc.Vrfs {
			vrfs[i] = vrf.Copy()
		}
		vrfSet = types.NewVrfSet(genDoc.InitialHeight, steeringMemberCandidateSet, vrfs)
	}

	if genDoc.NextVrfs == nil {
		nextVrfSet = types.NewVrfSet(genDoc.InitialHeight, steeringMemberCandidateSet, nil)
	} else {
		nextVrfs := make([]*types.Vrf, len(genDoc.NextVrfs))
		for i, nextVrf := range genDoc.NextVrfs {
			nextVrfs[i] = nextVrf.Copy()
		}
		nextVrfSet = types.NewVrfSet(genDoc.InitialHeight, steeringMemberCandidateSet, nextVrfs)
	}

	standingMemberSet.SetCoordinator(qrnSet)
	_, proposer := validatorSet.GetByAddress(standingMemberSet.Coordinator.PubKey.Address())
	validatorSet.Proposer = proposer

	return State{
		Version:       InitStateVersion,
		ChainID:       genDoc.ChainID,
		InitialHeight: genDoc.InitialHeight,

		LastBlockHeight: 0,
		LastBlockID:     types.BlockID{},
		LastBlockTime:   genDoc.GenesisTime,

		NextValidators:              nextValidatorSet,
		Validators:                  validatorSet,
		LastValidators:              types.NewValidatorSet(nil),
		LastHeightValidatorsChanged: genDoc.InitialHeight,

		ConsensusParams:                  *genDoc.ConsensusParams,
		LastHeightConsensusParamsChanged: genDoc.InitialHeight,

		AppHash: genDoc.AppHash,

		StandingMemberSet:                standingMemberSet,
		LastHeightStandingMembersChanged: genDoc.InitialHeight,

		SteeringMemberCandidateSet:                steeringMemberCandidateSet,
		LastHeightSteeringMemberCandidatesChanged: genDoc.InitialHeight,

		ConsensusRound:                  genDoc.ConsensusRound.ToProto(),
		LastHeightConsensusRoundChanged: genDoc.InitialHeight,

		QrnSet:     qrnSet,
		NextQrnSet: nextQrnSet,

		VrfSet:     vrfSet,
		NextVrfSet: nextVrfSet,

		LastHeightQrnChanged: genDoc.InitialHeight,
		LastHeightNextQrnChanged: genDoc.InitialHeight,
		LastHeightVrfChanged: genDoc.InitialHeight,
		LastHeightNextVrfChanged: genDoc.InitialHeight,
		LastHeightSettingSteeringMemberChanged: genDoc.InitialHeight,
	}, nil
}
