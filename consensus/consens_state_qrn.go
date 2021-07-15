package consensus

import (
	"fmt"

	tmrand "github.com/reapchain/reapchain-core/libs/rand"
	"github.com/reapchain/reapchain-core/p2p"
	"github.com/reapchain/reapchain-core/types"
	tmtime "github.com/reapchain/reapchain-core/types/time"
)

func (cs *State) addQrn(qrn *types.Qrn, peerID p2p.ID) (added bool, err error) {
	if qrn.Height+1 != cs.Height {
		return
	}

	standingMemberIndex, standingMember := cs.StandingMemberSet.GetByAddress(qrn.StandingMemberPubKey.Address())
	if standingMember == nil {
		return false, fmt.Errorf(
			"cannot find standing member in standing member")
	}

	qrn.StandingMemberIndex = standingMemberIndex

	added, err = cs.HeightQrnSet.AddQrn(qrn, peerID)
	if !added {
		return
	}

	if err := cs.eventBus.PublishEventQrn(types.EventDataQrn{Qrn: qrn}); err != nil {
		return added, err
	}
	cs.evsw.FireEvent(types.EventQrn, qrn)

	return added, err
}

func (cs *State) signAddQrn() *types.Qrn {
	if cs.privValidator == nil { // the node does not have a key
		return nil
	}

	if cs.privValidatorPubKey == nil {
		// Qrn won't be signed, but it's not critical.
		cs.Logger.Error(fmt.Sprintf("signAddQrn: %v", errPubKeyIsNotSet))
		return nil
	}

	// If the node not in the validator set, do nothing.
	if !cs.StandingMemberSet.HasAddress(cs.privValidatorPubKey.Address()) {
		return nil
	}

	qrn, err := cs.signQrn()
	if err == nil {
		cs.sendInternalMessage(MsgInfo{&QrnMessage{qrn}, ""})

		return qrn
	}
	return nil
}

func (cs *State) signQrn() (*types.Qrn, error) {
	// Flush the WAL. Otherwise, we may not recompute the same vote to sign,
	// and the privValidator will refuse to sign anything.
	if err := cs.wal.FlushAndSync(); err != nil {
		return nil, err
	}

	if cs.privValidatorPubKey == nil {
		return nil, errPubKeyIsNotSet
	}

	addr := cs.privValidatorPubKey.Address()
	index, _ := cs.StandingMemberSet.GetByAddress(addr)

	rand := tmrand.Uint64()

	qrn := &types.Qrn{
		Height:               cs.Height,
		Timestamp:            tmtime.Now(),
		StandingMemberIndex:  index,
		StandingMemberPubKey: cs.privValidatorPubKey,
		Value:                rand,
	}

	qrnProto, err := qrn.ToProto()
	if err != nil {
		fmt.Println("can't get qrn proto", "err", err)
	}

	err = cs.privValidator.SignQrn(cs.state.ChainID, qrnProto)
	qrn.Signature = qrnProto.Signature

	return qrn, err

}
