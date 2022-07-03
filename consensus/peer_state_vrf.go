package consensus

import (
	"github.com/reapchain/reapchain-core/libs/bits"
	"github.com/reapchain/reapchain-core/types"
)

func (ps *PeerState) setHasVrf(height int64, index int32) {
	psVrfs := ps.getVrfBitArray(height)
	if psVrfs != nil {
		psVrfs.SetIndex(int(index), true)
	}
}

func (ps *PeerState) SetHasVrf(vrf *types.Vrf) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	ps.setHasVrf(vrf.Height, vrf.VrfIndex)
}

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
		return vrf, true
	}
	return nil, false
}

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
