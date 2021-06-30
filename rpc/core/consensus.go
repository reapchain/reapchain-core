package core

import (
	cm "gitlab.reappay.net/sucs-lab//reapchain/consensus"
	tmmath "gitlab.reappay.net/sucs-lab//reapchain/libs/math"
	ctypes "gitlab.reappay.net/sucs-lab//reapchain/rpc/core/types"
	rpctypes "gitlab.reappay.net/sucs-lab//reapchain/rpc/jsonrpc/types"
	"gitlab.reappay.net/sucs-lab//reapchain/types"
)

// Validators gets the validator set at the given block height.
//
// If no height is provided, it will fetch the latest validator set. Note the
// validators are sorted by their voting power - this is the canonical order
// for the validators in the set as used in computing their Merkle root.
//
// More: https://docs.reapchain.com/master/rpc/#/Info/validators
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
//
// More: https://docs.reapchain.com/master/rpc/#/Info/standingMembers
func StandingMembers(ctx *rpctypes.Context, heightPtr *int64, pagePtr, perPagePtr *int) (*ctypes.ResultStandingMembers, error) {
	height, err := getHeight(latestUncommittedHeight(), heightPtr)
	if err != nil {
		return nil, err
	}

	standingMembers, err := env.StateStore.LoadStandingMembers(height)
	if err != nil {
		return nil, err
	}

	totalCount := len(standingMembers.StandingMembers)
	perPage := validatePerPage(perPagePtr)
	page, err := validatePage(pagePtr, perPage, totalCount)
	if err != nil {
		return nil, err
	}

	skipCount := validateSkipCount(page, perPage)

	v := standingMembers.StandingMembers[skipCount : skipCount+tmmath.MinInt(perPage, totalCount-skipCount)]

	return &ctypes.ResultStandingMembers{
		BlockHeight:     height,
		StandingMembers: v,
		Count:           len(v),
		Total:           totalCount}, nil
}

func Qns(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.ResultQns, error) {
	height, err := getHeight(latestUncommittedHeight(), heightPtr)
	if err != nil {
		return nil, err
	}

	qns, err := env.StateStore.LoadQns(height)
	if err != nil {
		return nil, err
	}

	totalCount := len(qns.Qns)
	if err != nil {
		return nil, err
	}

	v := qns.Qns[:]

	return &ctypes.ResultQns{
		BlockHeight: height,
		Qns:         qns.Qns,
		Count:       len(v),
		Total:       totalCount}, nil
}

// DumpConsensusState dumps consensus state.
// UNSTABLE
// More: https://docs.reapchain.com/master/rpc/#/Info/dump_consensus_state
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
// More: https://docs.reapchain.com/master/rpc/#/Info/consensus_state
func ConsensusState(ctx *rpctypes.Context) (*ctypes.ResultConsensusState, error) {
	// Get self round state.
	bz, err := env.ConsensusState.GetRoundStateSimpleJSON()
	return &ctypes.ResultConsensusState{RoundState: bz}, err
}

// ConsensusParams gets the consensus parameters at the given block height.
// If no height is provided, it will fetch the latest consensus params.
// More: https://docs.reapchain.com/master/rpc/#/Info/consensus_params
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
