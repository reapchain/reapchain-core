package consensus

import (
	"fmt"

	"github.com/reapchain/reapchain-core/p2p"
	"github.com/reapchain/reapchain-core/types"
)

func (cs *State) tryAddVrf(vrf *types.Vrf, peerID p2p.ID) (bool, error) {
	if err := cs.state.NextVrfSet.AddVrf(vrf); err != nil {
		fmt.Println("stompesi - tryAddVrf", err)
		return false, err
	}
	
	return true, nil
}
