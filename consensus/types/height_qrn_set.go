package types

import (
	"fmt"
	"strings"
	"sync"

	tmjson "github.com/reapchain/reapchain-core/libs/json"
	"github.com/reapchain/reapchain-core/p2p"
	"github.com/reapchain/reapchain-core/types"
)

type HeightQrnSet struct {
	chainID           string //??
	height            int64
	standingMemberSet *types.StandingMemberSet

	mtx    sync.Mutex
	qrnSet *types.QrnSet
}

func NewHeightQrnSet(chainID string, height int64, standingMemberSet *types.StandingMemberSet) *HeightQrnSet {
	heightQrnSet := &HeightQrnSet{
		chainID: chainID,
	}
	heightQrnSet.Reset(height, standingMemberSet)
	return heightQrnSet
}

func (heightQrnSet *HeightQrnSet) Reset(height int64, standingMemberSet *types.StandingMemberSet) {
	heightQrnSet.mtx.Lock()
	defer heightQrnSet.mtx.Unlock()

	heightQrnSet.height = height
	heightQrnSet.standingMemberSet = standingMemberSet
}

func (heightQrnSet *HeightQrnSet) Height() int64 {
	heightQrnSet.mtx.Lock()
	defer heightQrnSet.mtx.Unlock()
	return heightQrnSet.height
}

// Duplicate qrns return added=false, err=nil.
// By convention, peerID is "" if origin is self.
func (heightQrnSet *HeightQrnSet) AddQrn(qrn *types.Qrn, peerID p2p.ID) (added bool, err error) {
	heightQrnSet.mtx.Lock()
	defer heightQrnSet.mtx.Unlock()

	if heightQrnSet.qrnSet == nil {
		//TODO:
		// if rndz := heightQrnSet.peerCatchupRounds[peerID]; len(rndz) < 2 {
		// 	heightQrnSet.addRound(qrn.Round)
		// 	qrnSet = heightQrnSet.getQrnSet(qrn.Round, qrn.Type)
		// 	heightQrnSet.peerCatchupRounds[peerID] = append(rndz, qrn.Round)
		// } else {
		// 	// punish peer
		// 	err = ErrGotQrnFromUnwantedRound
		// 	return
		// }
	}
	added, err = heightQrnSet.qrnSet.AddQrn(qrn)
	return
}

func (heightQrnSet *HeightQrnSet) Preqrns(round int32) *types.QrnSet {
	heightQrnSet.mtx.Lock()
	defer heightQrnSet.mtx.Unlock()
	return heightQrnSet.qrnSet
}

//---------------------------------------------------------
// string and json

func (heightQrnSet *HeightQrnSet) String() string {
	return heightQrnSet.StringIndented("")
}

func (heightQrnSet *HeightQrnSet) StringIndented(indent string) string {
	heightQrnSet.mtx.Lock()
	defer heightQrnSet.mtx.Unlock()
	//TODO:
	vsStrings := make([]string, 0, 1)
	// rounds 0 ~ heightQrnSet.round inclusive
	qrnSetString := heightQrnSet.qrnSet.StringShort()
	vsStrings = append(vsStrings, qrnSetString)
	return fmt.Sprintf(`HeightQrnSet{H:%v %s  %v
%s}`,
		heightQrnSet.height, indent, strings.Join(vsStrings, "\n"+indent+"  "), indent)
}

func (heightQrnSet *HeightQrnSet) MarshalJSON() ([]byte, error) {
	heightQrnSet.mtx.Lock()
	defer heightQrnSet.mtx.Unlock()
	return tmjson.Marshal(heightQrnSet.toAllRoundQrns())
}

func (heightQrnSet *HeightQrnSet) toAllRoundQrns() []roundQrns {

	allQrns := make([]roundQrns, 1)
	// rounds 0 ~ heightQrnSet.round inclusive
	allQrns[0] = roundQrns{
		Qrns:         heightQrnSet.qrnSet.QrnStrings(),
		QrnsBitArray: heightQrnSet.qrnSet.BitArrayString(),
	}
	// TODO: all other peer catchup rounds
	return allQrns
}

type roundQrns struct {
	Qrns         []string `json:"qrns"`
	QrnsBitArray string   `json:"qrns_bit_array"`
}
