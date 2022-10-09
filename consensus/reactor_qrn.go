package consensus

import (
	"bytes"
	"fmt"
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
		prs := ps.GetRoundState()
		
		if rs.LockedBlock != nil {
			consensusStartBlockHeight := rs.LockedBlock.ConsensusRound.ConsensusStartBlockHeight
			roundPeriod := rs.LockedBlock.ConsensusRound.Period

			if consensusStartBlockHeight <= prs.Height && prs.Height < consensusStartBlockHeight + int64(roundPeriod) {
				if ps.PickSendQrn(conR.conS.state.NextQrnSet) {
					continue OUTER_LOOP
				}
			}
		}

		time.Sleep(conR.conS.config.PeerGossipSleepDuration)
	}
}

func (conR *Reactor) tryAddCatchupQrnMessage(qrnMessage *QrnMessage) (error) {
	if qrnMessage == nil {
		return fmt.Errorf("QrnMessage is nil")
	}

	if qrnMessage.Qrn.VerifySign() == false {
		return fmt.Errorf("Invalid qrn sign")
	}

	for idx, currentQrnMessage := range conR.CatchupQrnMessages {
		if bytes.Equal(currentQrnMessage.Qrn.StandingMemberPubKey.Address(), qrnMessage.Qrn.StandingMemberPubKey.Address()) {
			if currentQrnMessage.Qrn.Height < qrnMessage.Qrn.Height {
				conR.CatchupQrnMessages[idx] = &QrnMessage{Qrn: qrnMessage.Qrn.Copy()}
			}
			return nil
		}
	}

	conR.CatchupQrnMessages = append(conR.CatchupQrnMessages, &QrnMessage{Qrn: qrnMessage.Qrn.Copy()})
	return nil
}