package consensus

import (
	bc "github.com/reapchain/reapchain-core/blockchain"
	"github.com/reapchain/reapchain-core/p2p"
	bcproto "github.com/reapchain/reapchain-core/proto/reapchain-core/blockchain"
	sm "github.com/reapchain/reapchain-core/state"
)

// the reactor send the states of settingSteeringMember, nextVrfSet, nextQrnSet for sync.
// these states are related to next consensus round's validators.
func (conR *Reactor) respondStateToPeer(height int64,
	peer p2p.Peer) (queued bool) {
	settingSteeringMember, err := conR.stateStore.LoadSettingSteeringMember(height+1)
	if err != nil {
		return false
	}

	nextVrfSet, err := conR.stateStore.LoadNextVrfSet(height)
	if err != nil {
		return false
	}

	nextQrnSet, err := conR.stateStore.LoadNextQrnSet(height)
	if err != nil {
		return false
	}

	state := &sm.State{
		LastBlockHeight:            height,
		SettingSteeringMember:      settingSteeringMember,
		NextVrfSet:                 nextVrfSet,
		NextQrnSet:                 nextQrnSet,
	}
	
	stateProto, err := state.ToProto()
	if err != nil {
		return false
	}

	msgBytes, err := bc.EncodeMsg(&bcproto.StateResponse{State: stateProto})

	if err != nil {
		conR.Logger.Error("could not marshal msg", "err", err)
		return false
	}

	return peer.TrySend(CatchUpChannel, msgBytes)
}
