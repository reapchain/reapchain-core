package consensus

import (
	"time"

	"github.com/reapchain/reapchain-core/p2p"
)


func (conR *Reactor) gossipVrfsRoutine(peer p2p.Peer, ps *PeerState) {
	logger := conR.Logger.With("peer", peer)

OUTER_LOOP:
	for {
		// Manage disconnects from self or peer.
		if !peer.IsRunning() || !conR.IsRunning() {
			logger.Info("Stopping gossipVrfRoutine for peer")
			return
		}
		rs := conR.conS.GetRoundState()
		prs := ps.GetRoundState()

		if rs.Height == prs.Height {
			if ps.PickSendVrf(conR.conS.state.NextVrfSet) {
				continue OUTER_LOOP
			}
		}

		time.Sleep(conR.conS.config.PeerGossipSleepDuration)
	}
}
