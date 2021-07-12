package types

import (
	"sync"

	"github.com/reapchain/reapchain-core/p2p"
	"github.com/reapchain/reapchain-core/types"
)

type HeightQrnSet struct {
	chainID           string
	height            int64
	standingMemberSet *types.StandingMemberSet

	mtx    sync.Mutex
	qrnSet *types.QrnSet // keys: [0...round]
	// peerCatchupRounds map[p2p.ID][]int64 // keys: peer.ID; standingMemberues: at most 2 rounds
}

func NewHeightQrnSet(height int64, standingMemberSet *types.StandingMemberSet, qrnSet *types.QrnSet) *HeightQrnSet {
	heightQrnSet := &HeightQrnSet{}
	heightQrnSet.Reset(height, standingMemberSet, qrnSet)
	return heightQrnSet
}

func (heightQrnSet *HeightQrnSet) Reset(height int64, standingMemberSet *types.StandingMemberSet, qrnSet *types.QrnSet) {
	heightQrnSet.mtx.Lock()
	defer heightQrnSet.mtx.Unlock()

	heightQrnSet.height = height
	heightQrnSet.standingMemberSet = standingMemberSet
	heightQrnSet.qrnSet = qrnSet
}

func (heightQrnSet *HeightQrnSet) Height() int64 {
	heightQrnSet.mtx.Lock()
	defer heightQrnSet.mtx.Unlock()
	return heightQrnSet.height
}

func (hieghtQrnSet *HeightQrnSet) AddQrn(qrn *types.Qrn, peerID p2p.ID) (added bool, err error) {
	hieghtQrnSet.mtx.Lock()
	defer hieghtQrnSet.mtx.Unlock()

	added, err = hieghtQrnSet.qrnSet.AddQrn(qrn)
	return
}
