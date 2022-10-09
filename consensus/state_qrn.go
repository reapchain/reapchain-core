package consensus

import (
	"github.com/reapchain/reapchain-core/p2p"
	"github.com/reapchain/reapchain-core/types"
)

func (cs *State) tryAddQrn(qrn *types.Qrn, peerID p2p.ID) (bool, error) {
	if cs.state.NextQrnSet.AddQrn(qrn) {
		if err := cs.eventBus.PublishEventQrn(types.EventDataQrn{Qrn: qrn}); err != nil {
			return true, err
		}
	
		cs.evsw.FireEvent(types.EventQrn, qrn)
	}


	return true, nil
}

