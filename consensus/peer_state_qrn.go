package consensus

import (
	"fmt"

	"github.com/reapchain/reapchain-core/libs/bits"
	"github.com/reapchain/reapchain-core/types"
)

func (ps *PeerState) setHasQrn(height int64, index int32) {
	psQrns := ps.getQrnBitArray(height)
	if psQrns != nil {
		psQrns.SetIndex(int(index), true)
	}
}

func (ps *PeerState) SetHasQrn(qrn *types.Qrn) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	ps.setHasQrn(qrn.Height, qrn.StandingMemberIndex)
}

func (ps *PeerState) PickSendQrn(qrnSet types.QrnSetReader) bool {
	if qrn, ok := ps.PickQrnToSend(qrnSet); ok {
		msg := &QrnMessage{qrn}

		if ps.peer.Send(QrnChannel, MustEncode(msg)) {
			ps.SetHasQrn(qrn)
			return true
		} else {
			fmt.Println("Faile to send")
		}
		return false
	}
	return false
}

func (ps *PeerState) PickQrnToSend(qrnSet types.QrnSetReader) (qrn *types.Qrn, ok bool) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	if qrnSet.Size() == 0 {
		return nil, false
	}

	height, size := qrnSet.GetHeight(), qrnSet.Size()

	ps.ensureQrnBitArrays(height, size)

	psQrnBitArray := ps.getQrnBitArray(height)

	if psQrnBitArray == nil {
		return nil, false // Not something worth sending
	}

	if index, ok := qrnSet.BitArray().Sub(psQrnBitArray).PickRandom(); ok {
		qrn := qrnSet.GetByIndex(int32(index))
		fmt.Println("pick send qrn", qrn.StandingMemberIndex, int32(index), qrn.Value, height, ps.peer.ID())
		return qrn, true
	}
	return nil, false
}

func (ps *PeerState) EnsureQrnBitArrays(height int64, numStandingMembers int) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()
	ps.ensureQrnBitArrays(height, numStandingMembers)
}

func (ps *PeerState) ensureQrnBitArrays(height int64, numStandingMembers int) {
	if ps.NextConsensusStartBlockHeight == height {
		if ps.QrnsBitArray == nil {
			ps.QrnsBitArray = bits.NewBitArray(numStandingMembers)
		}
	}
}

func (ps *PeerState) getQrnBitArray(height int64) *bits.BitArray {
	if ps.NextConsensusStartBlockHeight == height {
		return ps.QrnsBitArray
	}

	return nil
}
