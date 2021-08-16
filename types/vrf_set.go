package types

import (
	"errors"
	"fmt"

	"github.com/reapchain/reapchain-core/crypto"
	"github.com/reapchain/reapchain-core/crypto/merkle"
	"github.com/reapchain/reapchain-core/libs/bits"
	tmsync "github.com/reapchain/reapchain-core/libs/sync"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

type VrfSet struct {
	Height int64

	SteeringMemberCandidateSet *SteeringMemberCandidateSet

	mtx          tmsync.Mutex
	Vrfs         []*Vrf
	VrfsBitArray *bits.BitArray
}

func NewVrfSet(height int64, steeringMemberCandidateSet *SteeringMemberCandidateSet, vrfs []*Vrf) *VrfSet {
	if vrfs == nil {
		vrfs = make([]*Vrf, steeringMemberCandidateSet.Size())
		for i, steeringMemberCandidate := range steeringMemberCandidateSet.SteeringMemberCandidates {
			vrfs[i] = NewVrfAsEmpty(height, steeringMemberCandidate.PubKey)
		}
	}

	return &VrfSet{
		Height:                     height,
		SteeringMemberCandidateSet: steeringMemberCandidateSet,
		VrfsBitArray:               bits.NewBitArray(steeringMemberCandidateSet.Size()),
		Vrfs:                       vrfs,
	}
}

func VrfSetFromProto(vrfSetProto *tmproto.VrfSet) (*VrfSet, error) {
	if vrfSetProto == nil {
		return nil, errors.New("nil vrf set")
	}

	vrfSet := new(VrfSet)
	vrfs := make([]*Vrf, len(vrfSetProto.Vrfs))

	for i, vrfProto := range vrfSetProto.Vrfs {
		vrf := VrfFromProto(vrfProto)
		if vrf == nil {
			return nil, errors.New(nilVrfStr)
		}
		vrfs[i] = vrf
	}
	vrfSet.Height = vrfSetProto.Height
	vrfSet.Vrfs = vrfs

	return vrfSet, vrfSet.ValidateBasic()
}

func (vrfSet *VrfSet) ValidateBasic() error {
	if vrfSet == nil || len(vrfSet.Vrfs) == 0 {
		return errors.New("vrf set is nil or empty")
	}

	for idx, vrf := range vrfSet.Vrfs {
		if err := vrf.ValidateBasic(); err != nil {
			return fmt.Errorf("Invalid vrf #%d: %w", idx, err)
		}
	}

	return nil
}

func (vrfSet *VrfSet) GetHeight() int64 {
	if vrfSet == nil {
		return 0
	}
	return vrfSet.Height
}

func (vrfSet *VrfSet) Size() int {
	if vrfSet == nil {
		return 0
	}

	return len(vrfSet.Vrfs)
}

func (vrfSet *VrfSet) IsNilOrEmpty() bool {
	return vrfSet == nil || len(vrfSet.Vrfs) == 0
}

func (vrfSet *VrfSet) AddVrf(vrf *Vrf) error {
	vrfSet.mtx.Lock()
	defer vrfSet.mtx.Unlock()

	if vrf == nil {
		return fmt.Errorf("Vrf is nil")
	}

	if vrf.Verify() == false {
		return fmt.Errorf("Invalid vrf sign")
	}

	steeringMemberCandidateIndex, _ := vrfSet.SteeringMemberCandidateSet.GetSteeringMemberCandidateByAddress(vrf.SteeringMemberCandidatePubKey.Address())

	if steeringMemberCandidateIndex == -1 {
		return fmt.Errorf("Not exist standing member of vrf: %v", vrf.SteeringMemberCandidatePubKey.Address())
	}

	vrf.SteeringMemberCandidateIndex = steeringMemberCandidateIndex
	vrfSet.Vrfs[steeringMemberCandidateIndex] = vrf.Copy()

	vrfSet.VrfsBitArray.SetIndex(int(steeringMemberCandidateIndex), true)

	return nil
}

func (vrfSet *VrfSet) GetVrf(steeringMemberCandidatePubKey crypto.PubKey) (vrf *Vrf) {
	steeringMemberCandidateIndex, _ := vrfSet.SteeringMemberCandidateSet.GetSteeringMemberCandidateByAddress(steeringMemberCandidatePubKey.Address())

	if steeringMemberCandidateIndex != -1 {
		return vrfSet.Vrfs[steeringMemberCandidateIndex]
	}

	return nil
}

func (vrfSet *VrfSet) Hash() []byte {
	vrfBytesArray := make([][]byte, len(vrfSet.Vrfs))
	for i, vrf := range vrfSet.Vrfs {
		if vrf != nil {
			vrfBytesArray[i] = vrf.GetVrfBytes()
		}
	}
	return merkle.HashFromByteSlices(vrfBytesArray)
}

func (vrfSet *VrfSet) ToProto() (*tmproto.VrfSet, error) {
	if vrfSet.IsNilOrEmpty() {
		return &tmproto.VrfSet{}, nil
	}

	vrfSetProto := new(tmproto.VrfSet)
	vrfsProto := make([]*tmproto.Vrf, len(vrfSet.Vrfs))
	for i, vrf := range vrfSet.Vrfs {
		vrfProto := vrf.ToProto()
		if vrfProto != nil {
			vrfsProto[i] = vrfProto
		}
	}
	vrfSetProto.Height = vrfSet.Height
	vrfSetProto.Vrfs = vrfsProto

	return vrfSetProto, nil
}

func (vrfSet *VrfSet) Copy() *VrfSet {
	vrfSet.mtx.Lock()
	defer vrfSet.mtx.Unlock()

	vrfsCopy := make([]*Vrf, len(vrfSet.Vrfs))
	for i, vrf := range vrfSet.Vrfs {
		if vrf != nil {
			vrfsCopy[i] = vrf.Copy()
		}
	}
	return &VrfSet{
		Height:                     vrfSet.Height,
		SteeringMemberCandidateSet: vrfSet.SteeringMemberCandidateSet,
		Vrfs:                       vrfsCopy,
		VrfsBitArray:               vrfSet.VrfsBitArray,
	}
}

func (vrfSet *VrfSet) GetByIndex(steeringMemberCandidateIndex int32) *Vrf {
	if vrfSet == nil {
		return nil
	}

	vrfSet.mtx.Lock()
	defer vrfSet.mtx.Unlock()
	return vrfSet.Vrfs[steeringMemberCandidateIndex]
}

func (vrfSet *VrfSet) BitArray() *bits.BitArray {
	if vrfSet == nil {
		return nil
	}

	vrfSet.mtx.Lock()
	defer vrfSet.mtx.Unlock()
	return vrfSet.VrfsBitArray.Copy()
}

// -----
const nilVrfSetString = "nil-VrfSet"

func (vrfSet *VrfSet) StringShort() string {
	if vrfSet == nil {
		return nilVrfSetString
	}
	vrfSet.mtx.Lock()
	defer vrfSet.mtx.Unlock()

	return fmt.Sprintf(`VrfSet{H:%v %v}`,
		vrfSet.Height, vrfSet.VrfsBitArray)
}

func (vrfSet *VrfSet) VrfStrings() []string {
	vrfSet.mtx.Lock()
	defer vrfSet.mtx.Unlock()

	vrfStrings := make([]string, len(vrfSet.Vrfs))
	for i, vrf := range vrfSet.Vrfs {
		if vrf == nil {
			vrfStrings[i] = "nil-Vrf"
		} else {
			vrfStrings[i] = vrf.String()
		}
	}
	return vrfStrings
}

func (vrfSet *VrfSet) BitArrayString() string {
	vrfSet.mtx.Lock()
	defer vrfSet.mtx.Unlock()
	return vrfSet.bitArrayString()
}

func (vrfSet *VrfSet) bitArrayString() string {
	bAString := vrfSet.VrfsBitArray.String()
	return fmt.Sprintf("%s", bAString)
}

type VrfSetReader interface {
	GetHeight() int64
	Size() int
	BitArray() *bits.BitArray
	GetByIndex(int32) *Vrf
}

func (vrfSet *VrfSet) UpdateWithChangeSet(vrfs []*Vrf) error {
	return vrfSet.updateWithChangeSet(vrfs)
}

func (vrfSet *VrfSet) updateWithChangeSet(vrfs []*Vrf) error {
	if len(vrfs) != 0 {
		vrfSet.Vrfs = vrfs[:]
	}

	return nil
}
