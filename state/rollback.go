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

	height := bs.Height()

	// NOTE: persistence of state and blocks don't happen atomically. Therefore it is possible that
	// when the user stopped the node the state wasn't updated but the blockstore was. In this situation
	// we don't need to rollback any state and can just return early
	if height == invalidState.LastBlockHeight+1 {
		return invalidState.LastBlockHeight, invalidState.AppHash, nil
	}

	// If the state store isn't one below nor equal to the blockstore height than this violates the
	// invariant
	if height != invalidState.LastBlockHeight {
		return -1, nil, fmt.Errorf("statestore height (%d) is not one below or equal to blockstore height (%d)",
			invalidState.LastBlockHeight, height)
	}

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

	standingMembersChangeHeight := invalidState.LastHeightStandingMembersChanged
	if standingMembersChangeHeight > rollbackHeight {
		standingMembersChangeHeight = rollbackHeight + 1
	}

	steeringMemberCandidatesChangeHeight := invalidState.LastHeightSteeringMemberCandidatesChanged
	if steeringMemberCandidatesChangeHeight > rollbackHeight {
		steeringMemberCandidatesChangeHeight = rollbackHeight + 1
	}

	consensusRoundChangeHeight := invalidState.LastHeightConsensusRoundChanged
	if consensusRoundChangeHeight > rollbackHeight {
		consensusRoundChangeHeight = rollbackHeight + 1
	}

	qrnChangeHeight := invalidState.LastHeightQrnChanged
	if qrnChangeHeight > rollbackHeight {
		qrnChangeHeight = rollbackHeight + 1
	}

	nextQrnChangeHeight := invalidState.LastHeightNextQrnChanged
	if nextQrnChangeHeight > rollbackHeight {
		nextQrnChangeHeight = rollbackHeight + 1
	}

	vrfChangeHeight := invalidState.LastHeightVrfChanged
	if vrfChangeHeight > rollbackHeight {
		vrfChangeHeight = rollbackHeight + 1
	}

	nextVrfChangeHeight := invalidState.LastHeightNextVrfChanged
	if nextVrfChangeHeight > rollbackHeight {
		nextVrfChangeHeight = rollbackHeight + 1
	}

	settingSteeringMemberChangeHeight := invalidState.LastHeightSettingSteeringMemberChanged
	if settingSteeringMemberChangeHeight > rollbackHeight {
		settingSteeringMemberChangeHeight = rollbackHeight + 1
	}

	fmt.Println("LastBlockHeight", rollbackBlock.Header.Height)
	fmt.Println("LastBlockID", rollbackBlock.BlockID)
	fmt.Println("LastBlockTime", rollbackBlock.Header.Time)

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
