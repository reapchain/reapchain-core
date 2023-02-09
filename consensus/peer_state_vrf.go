package consensus

import (
	"github.com/reapchain/reapchain-core/libs/bits"
	"github.com/reapchain/reapchain-core/types"
)

// Set a flag that the peer has a vrf
func (ps *PeerState) SetHasVrf(vrf *types.Vrf) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	ps.setHasVrf(vrf.Height, vrf.VrfIndex)
}

func (ps *PeerState) setHasVrf(height int64, index int32) {
	psVrfs := ps.getVrfBitArray(height)
	if psVrfs != nil {
		psVrfs.SetIndex(int(index), true)
	}
}

// Select a random vrf which is not sent to the peer, and send the vrf to the peer
func (ps *PeerState) PickSendVrf(vrfSet types.VrfSetReader) bool {
	if vrf, ok := ps.PickVrfToSend(vrfSet); ok {
		msg := &VrfMessage{vrf}

		if ps.peer.Send(VrfChannel, MustEncode(msg)) {
			ps.SetHasVrf(vrf)
			return true
		}

		return false
	}
	return false
}

// Select a random vrf which is not sent to the peer.
func (ps *PeerState) PickVrfToSend(vrfSet types.VrfSetReader) (vrf *types.Vrf, ok bool) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	if vrfSet.Size() == 0 {
		return nil, false
	}

	height, size := vrfSet.GetHeight(), vrfSet.Size()

	ps.ensureVrfBitArrays(height, size)

	psVrfBitArray := ps.getVrfBitArray(height)

	if psVrfBitArray == nil {
		return nil, false // Not something worth sending
	}

	if index, ok := vrfSet.BitArray().Sub(psVrfBitArray).PickRandom(); ok {
		vrf := vrfSet.GetByIndex(int32(index))
		if vrf == nil {
			return nil, false
		} else if vrf.Value == nil {
			return nil, false
		}
		return vrf, true
	}
	return nil, false
}

// Check a bit array whether the bit array exists or not exists in the next consensus round of the peer. 
// After that, if the bit array does not exist, generate the bit array with a size of steering member candidate count.
func (ps *PeerState) EnsureVrfBitArrays(height int64, numSteeringMemberCandidates int) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()
	ps.ensureVrfBitArrays(height, numSteeringMemberCandidates)
}

func (ps *PeerState) ensureVrfBitArrays(height int64, numSteeringMemberCandidates int) {
	if ps.NextConsensusStartBlockHeight == height {
		if ps.VrfsBitArray == nil {
			ps.VrfsBitArray = bits.NewBitArray(numSteeringMemberCandidates)
		}
	}
}

// Get a bit array of the consensus round hight.
func (ps *PeerState) getVrfBitArray(height int64) *bits.BitArray {
	if ps.NextConsensusStartBlockHeight == height {
		return ps.VrfsBitArray
	}

	return nil
}

func (ps *PeerState) ApplyHasVrfMessage(msg *HasVrfMessage) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	if ps.PRS.Height != msg.Height {
		return
	}

	ps.setHasVrf(msg.Height, msg.Index)
}
