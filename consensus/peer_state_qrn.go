package consensus

import (
	"github.com/reapchain/reapchain-core/libs/bits"
	"github.com/reapchain/reapchain-core/types"
)

// Set a flag that the peer has a qrn
func (ps *PeerState) SetHasQrn(qrn *types.Qrn) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	ps.setHasQrn(qrn.Height, qrn.QrnIndex)
}

func (ps *PeerState) setHasQrn(height int64, index int32) {
	psQrns := ps.getQrnBitArray(height)
	if psQrns != nil {
		psQrns.SetIndex(int(index), true)
	}
}

// Select a random qrn which is not sent to the peer, and send the qrn to the peer
func (ps *PeerState) PickSendQrn(qrnSet types.QrnSetReader) bool {
	if qrn, ok := ps.pickQrnToSend(qrnSet); ok {
		msg := &QrnMessage{qrn}

		if ps.peer.Send(QrnChannel, MustEncode(msg)) {
			ps.SetHasQrn(qrn)
			return true
		}
		return false
	}
	return false
}

// Select a random qrn which is not sent to the peer.
func (ps *PeerState) pickQrnToSend(qrnSet types.QrnSetReader) (qrn *types.Qrn, ok bool) {
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
		if qrn == nil {
			return nil, false
		} else if qrn.Value == 0 {
			return nil, false
		}
		return qrn, true
	}
	return nil, false
}

// Check a bit array whether the bit array exists or not exists in the next consensus round of the peer. 
// After that, if the bit array does not exist, generate the bit array with a size of standing member count.
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

// Get a bit array of the consensus round hight.
func (ps *PeerState) getQrnBitArray(height int64) *bits.BitArray {
	if ps.NextConsensusStartBlockHeight == height {
		return ps.QrnsBitArray
	}

	return nil
}

func (ps *PeerState) ApplyHasQrnMessage(msg *HasQrnMessage) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	if ps.PRS.Height != msg.Height {
		return
	}

	ps.setHasQrn(msg.Height, msg.Index)
}
