package consensus

import (
	"fmt"
	"time"

	"github.com/reapchain/reapchain-core/p2p"
)

func (conR *Reactor) gossipSettingSteeringMemberRoutine(peer p2p.Peer, ps *PeerState) {
	logger := conR.Logger.With("peer", peer)

OUTER_LOOP:
	for {
		// Manage disconnects from self or peer.
		if !peer.IsRunning() || !conR.IsRunning() {
			logger.Info("Stopping gossipSettingSteeringMemberRoutine for peer")
			return
		}
		rs := conR.conS.GetRoundState()
		prs := ps.GetRoundState()

		if rs.Height == prs.Height {
			if ps.PickSendSettingSteeringMember(conR.conS.state.SettingSteeringMember) {
				time.Sleep(conR.conS.config.PeerGossipSleepDuration)
				continue OUTER_LOOP
			}
		}

		time.Sleep(conR.conS.config.PeerGossipSleepDuration)
	}
}

func (conR *Reactor) tryAddCatchupSettingSteeringMemberMessage(settingSteeringMemberMessage *SettingSteeringMemberMessage) (error) {
	if settingSteeringMemberMessage.SettingSteeringMember.VerifySign() == false {
		return fmt.Errorf("Invalid seeting steering member sign")
	}

	
	if conR.CatchupSettingSteeringMemberMessage == nil || 
		 conR.CatchupSettingSteeringMemberMessage.SettingSteeringMember.Height < settingSteeringMemberMessage.SettingSteeringMember.Height {
		conR.CatchupSettingSteeringMemberMessage = &SettingSteeringMemberMessage{settingSteeringMemberMessage.SettingSteeringMember}
	}

	return nil
}