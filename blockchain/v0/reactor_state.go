package v0

import (
	bc "github.com/reapchain/reapchain-core/blockchain"
	"github.com/reapchain/reapchain-core/p2p"
	bcproto "github.com/reapchain/reapchain-core/proto/reapchain-core/blockchain"
	sm "github.com/reapchain/reapchain-core/state"
)

func (bcR *BlockchainReactor) respondStateToPeer(msg *bcproto.StateRequest,
	src p2p.Peer) (queued bool) {
	settingSteeringMember, err := bcR.stateStore.LoadSettingSteeringMember(msg.Height+1)
	if err != nil {
		return false
	}

	nextVrfSet, err := bcR.stateStore.LoadNextVrfSet(msg.Height)
	if err != nil {
		return false
	}

	nextQrnSet, err := bcR.stateStore.LoadNextQrnSet(msg.Height)
	if err != nil {
		return false
	}

	state := &sm.State{
		LastBlockHeight:            msg.Height,
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
		bcR.Logger.Error("could not marshal msg", "err", err)
		return false
	}

	return src.TrySend(BlockchainChannel, msgBytes)
}
