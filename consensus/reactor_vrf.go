package consensus

import (
	"bytes"
	"fmt"
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
			roundPeriod := rs.LockedBlock.ConsensusRound.Period

			if consensusStartBlockHeight <= prs.Height && prs.Height < consensusStartBlockHeight + int64(roundPeriod) {
				if ps.PickSendVrf(conR.conS.state.NextVrfSet) {
					continue OUTER_LOOP
				}
			}
		}

		time.Sleep(conR.conS.config.PeerGossipSleepDuration)
	}
}

func (conR *Reactor) tryAddCatchupVrfMessage(vrfMessage *VrfMessage) (error) {
	if vrfMessage == nil {
		return fmt.Errorf("VrfMessage is nil")
	}

	if vrfMessage.Vrf.Verify() == false {
		return fmt.Errorf("Invalid vrf sign")
	}

	for idx, currentVrfMessage := range conR.CatchupVrfMessages {
		if bytes.Equal(currentVrfMessage.Vrf.SteeringMemberCandidatePubKey.Address(), vrfMessage.Vrf.SteeringMemberCandidatePubKey.Address()) {
			if currentVrfMessage.Vrf.Height < vrfMessage.Vrf.Height {
				conR.CatchupVrfMessages[idx] = &VrfMessage{Vrf: vrfMessage.Vrf.Copy()}
			}
			return nil
		}
	}

	conR.CatchupVrfMessages = append(conR.CatchupVrfMessages, &VrfMessage{Vrf: vrfMessage.Vrf.Copy()})
	return nil
}