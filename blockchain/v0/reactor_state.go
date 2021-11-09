package v0

// import (
// 	"fmt"

// 	bc "github.com/reapchain/reapchain-core/blockchain"
// 	"github.com/reapchain/reapchain-core/p2p"
// 	bcproto "github.com/reapchain/reapchain-core/proto/reapchain/blockchain"
// 	sm "github.com/reapchain/reapchain-core/state"
// )

// func (bcR *BlockchainReactor) respondStateToPeer(msg *bcproto.StateRequest,
// 	src p2p.Peer) (queued bool) {
// 	qrnSet, err := bcR.stateStore.LoadQrnSet(msg.Height)
// 	if err != nil {
// 		return false
// 	}

// 	if qrnSet != nil {

// 		vrfSet, err := bcR.stateStore.LoadVrfSet(msg.Height)
// 		if err != nil {
// 			return false
// 		}

// 		settingSteeringMember, err := bcR.stateStore.LoadSettingSteeringMember(msg.Height)
// 		if err != nil {
// 			return false
// 		}

// 		state := &sm.State{
// 			LastBlockHeight:       msg.Height,
// 			SettingSteeringMember: settingSteeringMember,
// 			VrfSet:                vrfSet,
// 			QrnSet:                qrnSet,
// 		}

// 		stateProto, err := state.ToProto()
// 		if err != nil {
// 			return false
// 		}

// 		fmt.Println("jw - 2", stateProto.LastBlockHeight)

// 		msgBytes, err := bc.EncodeMsg(&bcproto.StateResponse{State: stateProto})

// 		if err != nil {
// 			bcR.Logger.Error("could not marshal msg", "err", err)
// 			return false
// 		}

// 		return src.TrySend(StateChannel, msgBytes)
// 	}

// 	msgBytes, err := bc.EncodeMsg(&bcproto.NoStateResponse{Height: msg.Height})
// 	if err != nil {
// 		bcR.Logger.Error("could not convert msg to protobuf", "err", err)
// 		return false
// 	}

// 	return src.TrySend(StateChannel, msgBytes)
// }
