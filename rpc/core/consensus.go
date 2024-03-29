package core

import (
	"time"
	
	cm "github.com/reapchain/reapchain-core/consensus"
	tmmath "github.com/reapchain/reapchain-core/libs/math"
	ctypes "github.com/reapchain/reapchain-core/rpc/core/types"
	rpctypes "github.com/reapchain/reapchain-core/rpc/jsonrpc/types"
	"github.com/reapchain/reapchain-core/types"
)

// Validators gets the validator set at the given block height.
//
// If no height is provided, it will fetch the latest validator set. Note the
// validators are sorted by their voting power - this is the canonical order
// for the validators in the set as used in computing their Merkle root.
//
// More: https://docs.tendermint.com/master/rpc/#/Info/validators
func Validators(ctx *rpctypes.Context, heightPtr *int64, pagePtr, perPagePtr *int) (*ctypes.ResultValidators, error) {
	// The latest validator that we know is the NextValidator of the last block.
	height, err := getHeight(latestUncommittedHeight(), heightPtr)
	if err != nil {
		return nil, err
	}

	validators, err := env.StateStore.LoadValidators(height)
	if err != nil {
		return nil, err
	}

	totalCount := len(validators.Validators)
	perPage := validatePerPage(perPagePtr)
	page, err := validatePage(pagePtr, perPage, totalCount)
	if err != nil {
		return nil, err
	}

	skipCount := validateSkipCount(page, perPage)

	v := validators.Validators[skipCount : skipCount+tmmath.MinInt(perPage, totalCount-skipCount)]

	return &ctypes.ResultValidators{
		BlockHeight: height,
		Validators:  v,
		Count:       len(v),
		Total:       totalCount}, nil
}

// StandingMembers gets the standing member set at the given block height.
//
// If no height is provided, it will fetch the latest standing member set.

func StandingMembers(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.ResultStandingMembers, error) {
	height, err := getHeight(latestUncommittedHeight(), heightPtr)
	if err != nil {
		return nil, err
	}

	standingMemberSet, err := env.StateStore.LoadStandingMemberSet(height)
	if err != nil {
		return nil, err
	}

	return &ctypes.ResultStandingMembers{
		BlockHeight:               height,
		StandingMembers:           standingMemberSet.StandingMembers[:],
		CurrentCoordinatorRanking: standingMemberSet.CurrentCoordinatorRanking,
		Count:                     standingMemberSet.Size(),
	}, nil
}

// SteeringMemberCandidates gets the steering memeber candiate set at the given block height.
//
// If no height is provided, it will fetch the latest steering memeber candiate set.
func SteeringMemberCandidates(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.ResultSteeringMemberCandidates, error) {
	height, err := getHeight(latestUncommittedHeight(), heightPtr)
	if err != nil {
		return nil, err
	}

	standingMemberSet, err := env.StateStore.LoadSteeringMemberCandidateSet(height)
	if err != nil {
		return nil, err
	}

	return &ctypes.ResultSteeringMemberCandidates{
		BlockHeight:              height,
		SteeringMemberCandidates: standingMemberSet.SteeringMemberCandidates[:],
		Count:                    standingMemberSet.Size(),
	}, nil
}

// Qrns gets the qrn set at the given block height.
func Qrns(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.ResultQrns, error) {
	height, err := getHeight(latestUncommittedHeight(), heightPtr)
	if err != nil {
		return nil, err
	}

	qrnSet, err := env.StateStore.LoadQrnSet(height)
	if err != nil {
		return nil, err
	}

	return &ctypes.ResultQrns{
		BlockHeight: height,
		Qrns:        qrnSet.Qrns[:],
		Count:       qrnSet.Size(),
		QrnHash:       qrnSet.Hash(),
	}, nil
}

// NextQrns gets the next consensus round qrn set at the given block height.
func NextQrns(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.ResultQrns, error) {
	height, err := getHeight(latestUncommittedHeight(), heightPtr)
	if err != nil {
		return nil, err
	}

	qrnSet, err := env.StateStore.LoadNextQrnSet(height)
	if err != nil {
		return nil, err
	}

	return &ctypes.ResultQrns{
		BlockHeight: height,
		Qrns:        qrnSet.Qrns[:],
		Count:       qrnSet.Size(),
		QrnHash: 		 qrnSet.Hash(),
	}, nil
}

// SettingSteeringMember gets the steerin member list to be applied next consensus round at the given block height.
func SettingSteeringMember(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.ResultSettingSteeringMember, error) {
	height, err := getHeight(latestUncommittedHeight(), heightPtr)
	if err != nil {
		return nil, err
	}

	settingSteeringMember, err := env.StateStore.LoadSettingSteeringMember(height)
	if err != nil {
		return nil, err
	}

	if settingSteeringMember == nil {
		return &ctypes.ResultSettingSteeringMember{
			BlockHeight:           		height,
			Height:                		0,
			SteeringMemberAddresses: 	nil,
			Timestamp:             		time.Unix(0, 0),
			Address:               		[]byte(""),
		}, nil
	}

	return &ctypes.ResultSettingSteeringMember{
		BlockHeight:           height,
		Height:                settingSteeringMember.Height,
		SteeringMemberAddresses: settingSteeringMember.SteeringMemberAddresses,
		Timestamp:             settingSteeringMember.Timestamp,
		Address:               settingSteeringMember.CoordinatorPubKey.Address(),
	}, nil
}

// Vrf gets the vrf set at the given block height.
func Vrfs(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.ResultVrfs, error) {
	height, err := getHeight(latestUncommittedHeight(), heightPtr)
	if err != nil {
		return nil, err
	}

	vrfSet, err := env.StateStore.LoadVrfSet(height)
	if err != nil {
		return nil, err
	}

	return &ctypes.ResultVrfs{
		BlockHeight: height,
		Vrfs:        vrfSet.Vrfs[:],
		Count:       vrfSet.Size(),
		VrfHash:       vrfSet.Hash(),
	}, nil
}

// NextVrfs gets the next consensus round vrf set at the given block height.
func NextVrfs(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.ResultVrfs, error) {
	height, err := getHeight(latestUncommittedHeight(), heightPtr)
	if err != nil {
		return nil, err
	}

	vrfSet, err := env.StateStore.LoadNextVrfSet(height)
	if err != nil {
		return nil, err
	}

	return &ctypes.ResultVrfs{
		BlockHeight: height,
		Vrfs:        vrfSet.Vrfs[:],
		Count:       vrfSet.Size(),
		VrfHash:       vrfSet.Hash(),
	}, nil
}



// DumpConsensusState dumps consensus state.
// UNSTABLE
// More: https://docs.tendermint.com/master/rpc/#/Info/dump_consensus_state
func DumpConsensusState(ctx *rpctypes.Context) (*ctypes.ResultDumpConsensusState, error) {
	// Get Peer consensus states.
	peers := env.P2PPeers.Peers().List()
	peerStates := make([]ctypes.PeerStateInfo, len(peers))
	for i, peer := range peers {
		peerState, ok := peer.Get(types.PeerStateKey).(*cm.PeerState)
		if !ok { // peer does not have a state yet
			continue
		}
		peerStateJSON, err := peerState.ToJSON()
		if err != nil {
			return nil, err
		}
		peerStates[i] = ctypes.PeerStateInfo{
			// Peer basic info.
			NodeAddress: peer.SocketAddr().String(),
			// Peer consensus state.
			PeerState: peerStateJSON,
		}
	}
	// Get self round state.
	roundState, err := env.ConsensusState.GetRoundStateJSON()
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultDumpConsensusState{
		RoundState: roundState,
		Peers:      peerStates}, nil
}

// ConsensusState returns a concise summary of the consensus state.
// UNSTABLE
// More: https://docs.tendermint.com/master/rpc/#/Info/consensus_state
func ConsensusState(ctx *rpctypes.Context) (*ctypes.ResultConsensusState, error) {
	// Get self round state.
	bz, err := env.ConsensusState.GetRoundStateSimpleJSON()
	return &ctypes.ResultConsensusState{RoundState: bz}, err
}

// ConsensusParams gets the consensus parameters at the given block height.
// If no height is provided, it will fetch the latest consensus params.
// More: https://docs.tendermint.com/master/rpc/#/Info/consensus_params
func ConsensusParams(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.ResultConsensusParams, error) {
	// The latest consensus params that we know is the consensus params after the
	// last block.
	height, err := getHeight(latestUncommittedHeight(), heightPtr)
	if err != nil {
		return nil, err
	}

	consensusParams, err := env.StateStore.LoadConsensusParams(height)
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultConsensusParams{
		BlockHeight:     height,
		ConsensusParams: consensusParams}, nil
}
