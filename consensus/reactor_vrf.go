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

		if rs.LockedBlock != nil {
			consensusStartBlockHeight := rs.LockedBlock.ConsensusRound.ConsensusStartBlockHeight
			qrnPeriod := rs.LockedBlock.ConsensusRound.QrnPeriod
			vrfPeriod := rs.LockedBlock.ConsensusRound.VrfPeriod

			if prs.Height >= consensusStartBlockHeight + int64(qrnPeriod) && prs.Height < consensusStartBlockHeight + int64(qrnPeriod) + int64(vrfPeriod) {
				if ps.PickSendVrf(conR.conS.state.NextVrfSet) {
					time.Sleep(conR.conS.config.PeerGossipSleepDuration)
					continue OUTER_LOOP
				}
			}
		}

		time.Sleep(conR.conS.config.PeerGossipSleepDuration)
	}
}
