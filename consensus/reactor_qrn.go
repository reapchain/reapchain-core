package consensus

import (
	"time"

	"github.com/reapchain/reapchain-core/p2p"
)

func (conR *Reactor) gossipQrnsRoutine(peer p2p.Peer, ps *PeerState) {
	logger := conR.Logger.With("peer", peer)

OUTER_LOOP:
	for {
		// Manage disconnects from self or peer.
		if !peer.IsRunning() || !conR.IsRunning() {
			logger.Info("Stopping gossipQrnRoutine for peer")
			return
		}
		rs := conR.conS.GetRoundState()
		// prs := ps.GetRoundState()
		
		if rs.LockedBlock != nil {
			// consensusStartBlockHeight := rs.LockedBlock.ConsensusRound.ConsensusStartBlockHeight
			// roundPeriod := rs.LockedBlock.ConsensusRound.Period

			// if consensusStartBlockHeight <= prs.Height && prs.Height < consensusStartBlockHeight + int64(roundPeriod) {
				if ps.PickSendQrn(conR.conS.state.NextQrnSet) {
					continue OUTER_LOOP
				}
			// }
		}

		time.Sleep(conR.conS.config.PeerGossipSleepDuration)
	}
}