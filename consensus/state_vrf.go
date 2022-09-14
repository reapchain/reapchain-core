package consensus

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/reapchain/reapchain-core/p2p"
	"github.com/reapchain/reapchain-core/types"
)

func (cs *State) tryAddVrf(vrf *types.Vrf, peerID p2p.ID) (bool, error) {
	fmt.Println("stompesi - tryAddVrf", "vrf", vrf, "time", time.Now().UTC())

	if err := cs.state.NextVrfSet.AddVrf(vrf); err != nil {
		cs.Logger.Error("tryAddVrf FAILURE!!!", "err", err, "stack", string(debug.Stack()))
		return false, err
	}
	fmt.Println("stompesi - tryAddVrf - success", )
	
	return true, nil
}
