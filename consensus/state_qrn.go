package consensus

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/reapchain/reapchain-core/p2p"
	"github.com/reapchain/reapchain-core/types"
)

func (cs *State) tryAddQrn(qrn *types.Qrn, peerID p2p.ID) (bool, error) {
	fmt.Println("stompesi - tryAddQrn", "qrn", qrn, "time", time.Now().UTC())
	
	if err := cs.state.NextQrnSet.AddQrn(qrn); err != nil {
		cs.Logger.Error("tryAddQrn FAILURE!!!", "err", err, "stack", string(debug.Stack()))
		return false, err
	}

	return true, nil
}
