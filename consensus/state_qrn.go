package consensus

import (
	"github.com/reapchain/reapchain-core/p2p"
	"github.com/reapchain/reapchain-core/types"
)

func (cs *State) tryAddQrn(qrn *types.Qrn, peerID p2p.ID) (bool, error) {
	if err := cs.state.NextQrnSet.AddQrn(qrn); err != nil {
		return false, err
	}

	return true, nil
}
