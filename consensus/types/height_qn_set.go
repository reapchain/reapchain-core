package types

import (
	"sync"

	"gitlab.reappay.net/sucs-lab/reapchain/p2p"
	"gitlab.reappay.net/sucs-lab/reapchain/types"
)

type HeightQnSet struct {
	chainID           string
	height            int64
	standingMemberSet *types.StandingMemberSet

	mtx   sync.Mutex
	qnSet *types.QnSet // keys: [0...round]
	// peerCatchupRounds map[p2p.ID][]int64 // keys: peer.ID; standingMemberues: at most 2 rounds
}

func NewHeightQnSet(height int64, standingMemberSet *types.StandingMemberSet) *HeightQnSet {
	heightQnSet := &HeightQnSet{}
	heightQnSet.Reset(height, standingMemberSet)
	return heightQnSet
}

func (heightQnSet *HeightQnSet) Reset(height int64, standingMemberSet *types.StandingMemberSet) {
	heightQnSet.mtx.Lock()
	defer heightQnSet.mtx.Unlock()

	heightQnSet.height = height
	heightQnSet.standingMemberSet = standingMemberSet
}

func (heightQnSet *HeightQnSet) Height() int64 {
	heightQnSet.mtx.Lock()
	defer heightQnSet.mtx.Unlock()
	return heightQnSet.height
}

func (hieghtQnSet *HeightQnSet) AddQn(qn *types.Qn, peerID p2p.ID) (added bool, err error) {
	hieghtQnSet.mtx.Lock()
	defer hieghtQnSet.mtx.Unlock()

	added, err = hieghtQnSet.qnSet.AddQn(qn)
	return
}
