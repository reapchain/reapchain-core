package consensus

import (
	"github.com/reapchain/reapchain-core/p2p"
	"github.com/reapchain/reapchain-core/types"
)

func (cs *State) tryAddVrf(vrf *types.Vrf, peerID p2p.ID) (bool, error) {
	if cs.state.NextVrfSet.AddVrf(vrf) {
		if err := cs.eventBus.PublishEventVrf(types.EventDataVrf{Vrf: vrf}); err != nil {
			return true, err
		}
	
		cs.evsw.FireEvent(types.EventVrf, vrf)
	}
	
	return true, nil
}
