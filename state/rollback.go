package state

import (
	"errors"
	"fmt"

	tmstate "github.com/reapchain/reapchain-core/proto/reapchain-core/state"
	tmversion "github.com/reapchain/reapchain-core/proto/reapchain-core/version"
	"github.com/reapchain/reapchain-core/version"
)

// Rollback overwrites the current Tendermint state (height n) with the most
// recent previous state (height n - 1).
// Note that this function does not affect application state.
func Rollback(bs BlockStore, ss Store) (int64, []byte, error) {
	invalidState, err := ss.Load()
	if err != nil {
		return -1, nil, err
	}
	if invalidState.IsEmpty() {
		return -1, nil, errors.New("no state found")
	}

	rollbackHeight := invalidState.LastBlockHeight
	rollbackBlock := bs.LoadBlockMeta(rollbackHeight)
	if rollbackBlock == nil {
		return -1, nil, fmt.Errorf("block at height %d not found", rollbackHeight)
	}

	previousValidatorSet, err := ss.LoadValidators(rollbackHeight - 1)
	if err != nil {
		return -1, nil, err
	}

	previousParams, err := ss.LoadConsensusParams(rollbackHeight)
	if err != nil {
		return -1, nil, err
	}

	valChangeHeight := invalidState.LastHeightValidatorsChanged
	// this can only happen if the validator set changed since the last block
	if valChangeHeight > rollbackHeight {
		valChangeHeight = rollbackHeight
	}

	paramsChangeHeight := invalidState.LastHeightConsensusParamsChanged
	// this can only happen if params changed from the last block
	if paramsChangeHeight > rollbackHeight {
		paramsChangeHeight = rollbackHeight
	}

	standingMembersChangeHeight := invalidState.LastHeightStandingMembersChanged
	if standingMembersChangeHeight > rollbackHeight {
		standingMembersChangeHeight = rollbackHeight
	}

	steeringMemberCandidatesChangeHeight := invalidState.LastHeightSteeringMemberCandidatesChanged
	if steeringMemberCandidatesChangeHeight > rollbackHeight {
		steeringMemberCandidatesChangeHeight = rollbackHeight
	}

	consensusRoundChangeHeight := invalidState.LastHeightConsensusRoundChanged
	if consensusRoundChangeHeight > rollbackHeight {
		consensusRoundChangeHeight = rollbackHeight
	}

	qrnChangeHeight := invalidState.LastHeightQrnChanged
	if qrnChangeHeight > rollbackHeight {
		qrnChangeHeight = rollbackHeight
	}

	nextQrnChangeHeight := invalidState.LastHeightNextQrnChanged
	if nextQrnChangeHeight > rollbackHeight {
		nextQrnChangeHeight = rollbackHeight
	}

	vrfChangeHeight := invalidState.LastHeightVrfChanged
	if vrfChangeHeight > rollbackHeight {
		vrfChangeHeight = rollbackHeight
	}

	nextVrfChangeHeight := invalidState.LastHeightNextVrfChanged
	if nextVrfChangeHeight > rollbackHeight {
		nextVrfChangeHeight = rollbackHeight
	}

	settingSteeringMemberChangeHeight := invalidState.LastHeightSettingSteeringMemberChanged
	if settingSteeringMemberChangeHeight > rollbackHeight {
		settingSteeringMemberChangeHeight = rollbackHeight
	}

	// build the new state from the old state and the prior block
	rolledBackState := State{
		Version: tmstate.Version{
			Consensus: tmversion.Consensus{
				Block: version.BlockProtocol,
				App:   previousParams.Version.AppVersion,
			},
			Software: version.TMCoreSemVer,
		},
		// immutable fields
		ChainID:       invalidState.ChainID,
		InitialHeight: invalidState.InitialHeight,

		LastBlockHeight: invalidState.LastBlockHeight - 1,
		LastBlockID:     rollbackBlock.Header.LastBlockID,
		LastBlockTime:   rollbackBlock.Header.Time,

		NextValidators:              invalidState.Validators,
		Validators:                  invalidState.LastValidators,
		LastValidators:              previousValidatorSet,
		LastHeightValidatorsChanged: valChangeHeight,

		ConsensusParams:                  previousParams,
		LastHeightConsensusParamsChanged: paramsChangeHeight,

		LastResultsHash: rollbackBlock.Header.LastResultsHash,
		AppHash:         rollbackBlock.Header.AppHash,
		
		StandingMemberSet: invalidState.StandingMemberSet,
		LastHeightStandingMembersChanged: standingMembersChangeHeight,

		SteeringMemberCandidateSet: invalidState.SteeringMemberCandidateSet,
		LastHeightSteeringMemberCandidatesChanged: steeringMemberCandidatesChangeHeight,
		
		ConsensusRound: invalidState.ConsensusRound,
		LastHeightConsensusRoundChanged: consensusRoundChangeHeight,
		
		QrnSet: invalidState.QrnSet,
		LastHeightQrnChanged: qrnChangeHeight,

		NextQrnSet: invalidState.NextQrnSet,
		LastHeightNextQrnChanged: nextQrnChangeHeight,

		VrfSet: invalidState.VrfSet,
		LastHeightVrfChanged: vrfChangeHeight,

		NextVrfSet: invalidState.NextVrfSet,
		LastHeightNextVrfChanged: nextVrfChangeHeight, 

		SettingSteeringMember: invalidState.SettingSteeringMember,
		LastHeightSettingSteeringMemberChanged: settingSteeringMemberChangeHeight,

		IsSetSteeringMember: invalidState.IsSetSteeringMember,
	}

	// persist the new state. This overrides the invalid one. NOTE: this will also
	// persist the validator set and consensus params over the existing structures,
	// but both should be the same
	if err := ss.Save(rolledBackState); err != nil {
		return -1, nil, fmt.Errorf("failed to save rolled back state: %w", err)
	}

	return rolledBackState.LastBlockHeight, rolledBackState.AppHash, nil
}
