package state

import (
	"errors"
	"fmt"

	tmstate "github.com/reapchain/reapchain-core/proto/podc/state"
	tmversion "github.com/reapchain/reapchain-core/proto/podc/version"
	"github.com/reapchain/reapchain-core/version"
)

// Rollback overwrites the current ReapchainCore state (height n) with the most
// recent previous state (height n - 1).
// Note that this function does not affect application state.
func Rollback(bs BlockStore, ss Store) (int64, []byte, error) {
	invalidState, err := ss.LoadRollback()
	if err != nil {
		return -1, nil, err
	}
	if invalidState.IsEmpty() {
		return -1, nil, errors.New("no state found")
	}

	height := bs.Height()

	// NOTE: persistence of state and blocks don't happen atomically. Therefore it is possible that
	// when the user stopped the node the state wasn't updated but the blockstore was. In this situation
	// we don't need to rollback any state and can just return early
	if height == invalidState.LastBlockHeight + 1 {
		return invalidState.LastBlockHeight, invalidState.AppHash, nil
	}

	// state store height is equal to blockstore height. We're good to proceed with rolling back state
	rollbackHeight := invalidState.LastBlockHeight - 1
	rollbackBlock := bs.LoadBlockMeta(rollbackHeight)
	if rollbackBlock == nil {
		return -1, nil, fmt.Errorf("block at height %d not found", rollbackHeight)
	}
	
	// We also need to retrieve the latest block because the app hash and last
	// results hash is only agreed upon in the following block.
	latestBlock := bs.LoadBlockMeta(invalidState.LastBlockHeight)
	if latestBlock == nil {
		return -1, nil, fmt.Errorf("block at height %d not found", invalidState.LastBlockHeight)
	}

	previousLastValidatorSet, err := ss.LoadValidators(rollbackHeight)
	if err != nil {
		return -1, nil, err
	}

	previousParams, err := ss.LoadConsensusParams(rollbackHeight + 1)
	if err != nil {
		return -1, nil, err
	}

	valChangeHeight := invalidState.LastHeightValidatorsChanged
	// this can only happen if the validator set changed since the last block
	if valChangeHeight > rollbackHeight {
		valChangeHeight = rollbackHeight + 1
	}

	paramsChangeHeight := invalidState.LastHeightConsensusParamsChanged
	// this can only happen if params changed from the last block
	if paramsChangeHeight > rollbackHeight {
		paramsChangeHeight = rollbackHeight + 1
	}

	previousStandingMemberSet, err := ss.LoadStandingMemberSet(rollbackHeight)
	if err != nil {
		return -1, nil, err
	}

	standingMemberChangeHeight := invalidState.LastHeightStandingMembersChanged
	// this can only happen if the validator set changed since the last block
	if standingMemberChangeHeight > rollbackHeight {
		standingMemberChangeHeight = rollbackHeight + 1
	}

	previousSteeringMemberCandidateSet, err := ss.LoadSteeringMemberCandidateSet(rollbackHeight)
	if err != nil {
		return -1, nil, err
	}

	steeringMemberCandidateChangeHeight := invalidState.LastHeightSteeringMemberCandidatesChanged
	// this can only happen if the validator set changed since the last block
	if steeringMemberCandidateChangeHeight > rollbackHeight {
		steeringMemberCandidateChangeHeight = rollbackHeight + 1
	}

	previousConsensusRound, err := ss.LoadConsensusRound(rollbackHeight)
	if err != nil {
		return -1, nil, err
	}

	consensusRoundChangeHeight := invalidState.LastHeightConsensusRoundChanged
	// this can only happen if the validator set changed since the last block
	if consensusRoundChangeHeight > rollbackHeight {
		consensusRoundChangeHeight = rollbackHeight + 1
	}

	previousQrnSet, err := ss.LoadQrnSet(rollbackHeight)
	if err != nil {
		return -1, nil, err
	}

	qrnSetChangeHeight := invalidState.LastHeightQrnChanged
	// this can only happen if the validator set changed since the last block
	if qrnSetChangeHeight > rollbackHeight {
		qrnSetChangeHeight = rollbackHeight + 1
	}

	previousNextQrnSet, err := ss.LoadNextQrnSet(rollbackHeight)
	if err != nil {
		return -1, nil, err
	}

	nextQrnSetChangeHeight := invalidState.LastHeightNextQrnChanged
	// this can only happen if the validator set changed since the last block
	if nextQrnSetChangeHeight > rollbackHeight {
		nextQrnSetChangeHeight = rollbackHeight + 1
	}

	previousVrfSet, err := ss.LoadVrfSet(rollbackHeight)
	if err != nil {
		return -1, nil, err
	}

	vrfSetChangeHeight := invalidState.LastHeightVrfChanged
	// this can only happen if the validator set changed since the last block
	if vrfSetChangeHeight > rollbackHeight {
		vrfSetChangeHeight = rollbackHeight + 1
	}

	previousNextVrfSet, err := ss.LoadNextVrfSet(rollbackHeight)
	if err != nil {
		return -1, nil, err
	}

	nextVrfSetChangeHeight := invalidState.LastHeightNextVrfChanged
	// this can only happen if the validator set changed since the last block
	if nextVrfSetChangeHeight > rollbackHeight {
		nextVrfSetChangeHeight = rollbackHeight + 1
	}


	previousSettingSteeringMember, err := ss.LoadSettingSteeringMember(rollbackHeight)
	if err != nil {
		return -1, nil, err
	}

	settingSteeringMemberChangeHeight := invalidState.LastHeightSettingSteeringMemberChanged
	// this can only happen if the validator set changed since the last block
	if settingSteeringMemberChangeHeight > rollbackHeight {
		settingSteeringMemberChangeHeight = rollbackHeight + 1
	}

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

		LastBlockHeight: rollbackBlock.Header.Height,
		LastBlockID:     rollbackBlock.BlockID,
		LastBlockTime:   rollbackBlock.Header.Time,

		NextValidators:              invalidState.Validators,
		Validators:                  invalidState.LastValidators,
		LastValidators:              previousLastValidatorSet,
		LastHeightValidatorsChanged: valChangeHeight,

		ConsensusParams:                  previousParams,
		LastHeightConsensusParamsChanged: paramsChangeHeight,

		LastResultsHash: latestBlock.Header.LastResultsHash,
		AppHash:         latestBlock.Header.AppHash,

		StandingMemberSet: previousStandingMemberSet,
		LastHeightStandingMembersChanged: standingMemberChangeHeight,

		SteeringMemberCandidateSet: previousSteeringMemberCandidateSet,
		LastHeightSteeringMemberCandidatesChanged: steeringMemberCandidateChangeHeight,

		ConsensusRound: previousConsensusRound,
		LastHeightConsensusRoundChanged: consensusRoundChangeHeight,

		QrnSet: previousQrnSet,
		LastHeightQrnChanged: qrnSetChangeHeight,

		NextQrnSet: previousNextQrnSet,
		LastHeightNextQrnChanged: nextQrnSetChangeHeight,

		VrfSet: previousVrfSet,
		LastHeightVrfChanged: vrfSetChangeHeight,

		NextVrfSet: previousNextVrfSet,
		LastHeightNextVrfChanged: nextVrfSetChangeHeight,

		SettingSteeringMember: previousSettingSteeringMember,
		LastHeightSettingSteeringMemberChanged: settingSteeringMemberChangeHeight,

		IsSetSteeringMember: false,
	}

	// persist the new state. This overrides the invalid one. NOTE: this will also
	// persist the validator set and consensus params over the existing structures,
	// but both should be the same
	if err := ss.Save(rolledBackState); err != nil {
		return -1, nil, fmt.Errorf("failed to save rolled back state: %w", err)
	}

	bs.SaveRollbackBlock(rolledBackState.LastBlockHeight)

	return rolledBackState.LastBlockHeight, rolledBackState.AppHash, nil
}
