package types

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/reapchain/reapchain-core/crypto"
	"github.com/reapchain/reapchain-core/crypto/merkle"
	"github.com/reapchain/reapchain-core/libs/bits"
	tmsync "github.com/reapchain/reapchain-core/libs/sync"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain-core/types"
)

const MAXIMUM_STEERING_MEMBER_CANDIDATES = 30

type VrfSet struct {
	Height int64

	mtx          tmsync.Mutex
	Vrfs         []*Vrf
	VrfsBitArray *bits.BitArray
}

func NewVrfSet(height int64, steeringMemberCandidateSet *SteeringMemberCandidateSet, vrfs []*Vrf) *VrfSet {

	if vrfs == nil || len(vrfs) == 0 {
		vrfs = make([]*Vrf, steeringMemberCandidateSet.Size())
		for i, steeringMemberCandidate := range steeringMemberCandidateSet.SteeringMemberCandidates {
			vrfs[i] = NewVrfAsEmpty(height, steeringMemberCandidate.PubKey)
			vrfs[i].VrfIndex = int32(i)
		}
	}

	return &VrfSet{
		Height:                     height,
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
	vrfSet.VrfsBitArray = bits.NewBitArray(len(vrfs))
	vrfSet.Vrfs = vrfs

	return vrfSet, vrfSet.ValidateBasic()
}

func (vrfSet *VrfSet) ValidateBasic() error {
	// if vrfSet == nil || len(vrfSet.Vrfs) == 0 {
	// 	return errors.New("vrf set is nil or empty")
	// }

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

	if vrf.Value != nil {
		if vrf.Verify() == false {
			return fmt.Errorf("Invalid vrf sign")
		}
	
		if vrfSet.Height != vrf.Height {
			fmt.Println("stompesi - Invalid vrf height", "vrfSet.Height", vrfSet.Height, "vrfHeight", vrf.Height, vrf.SteeringMemberCandidatePubKey.Address(), vrf)
			return fmt.Errorf("Invalid vrf height")
		}
	}


	vrfIndex := vrfSet.GetVrfIndexByAddress(vrf.SteeringMemberCandidatePubKey.Address())
	if vrfIndex == -1 {
		return fmt.Errorf("Not exist steering member candidate of vrf: %v", vrf.SteeringMemberCandidatePubKey.Address())
	}

	if vrfSet.VrfsBitArray.GetIndex(int(vrfIndex)) == false {
		vrf.VrfIndex = vrfIndex
		vrfSet.Vrfs[vrfIndex] = vrf.Copy()
		vrfSet.VrfsBitArray.SetIndex(int(vrfIndex), true)
	}

	return nil
}

func (vrfSet *VrfSet) GetVrf(steeringMemberCandidatePubKey crypto.PubKey) (vrf *Vrf) {
	vrfIndex := vrfSet.GetVrfIndexByAddress(steeringMemberCandidatePubKey.Address())

	if vrfIndex != -1 {
		return vrfSet.Vrfs[vrfIndex]
	}

	return nil
}

func (vrfSet *VrfSet) Hash() []byte {
	vrfBytesArray := make([][]byte, len(vrfSet.Vrfs))
	for i, vrf := range vrfSet.Vrfs {
		if vrf != nil && vrf.Proof != nil {
			// fmt.Println("stompesi - vrf hash", vrf)
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
	if vrfSet == nil {
		return nil
	}

	vrfsCopy := make([]*Vrf, len(vrfSet.Vrfs))
	for i, vrf := range vrfSet.Vrfs {
		if vrf != nil {
			vrfsCopy[i] = vrf.Copy()
		}
	}

	return &VrfSet{
		Height:                     vrfSet.Height,
		Vrfs:                       vrfsCopy,
		VrfsBitArray:               vrfSet.VrfsBitArray.Copy(),
	}
}

func (vrfSet *VrfSet) GetByIndex(vrfIndex int32) *Vrf {
	if vrfSet == nil {
		return nil
	}

	return vrfSet.Vrfs[vrfIndex]
}

func (vrfSet *VrfSet) BitArray() *bits.BitArray {
	if vrfSet == nil {
		return nil
	}

	vrfSet.mtx.Lock()
	defer vrfSet.mtx.Unlock()
	return vrfSet.VrfsBitArray.Copy()
}

type VrfSetReader interface {
	GetHeight() int64
	Size() int
	BitArray() *bits.BitArray
	GetByIndex(int32) *Vrf
}

func (vrfSet *VrfSet) UpdateWithChangeSet(steeringMemberCandidateSet *SteeringMemberCandidateSet) error {
	vrfSet.mtx.Lock()
	defer vrfSet.mtx.Unlock()

	vrfs := make([]*Vrf, 0, len(steeringMemberCandidateSet.SteeringMemberCandidates))

	for _, steeringMemberCandidate := range steeringMemberCandidateSet.SteeringMemberCandidates {
		vrf := vrfSet.GetVrf(steeringMemberCandidate.PubKey)

		if vrf != nil {
			vrfs = append(vrfs, vrf.Copy())
		}
	}

	vrfsBitArray := bits.NewBitArray(len(vrfs))
	
	for i, vrf := range vrfs {
		vrfIndex := vrfSet.GetVrfIndexByAddress(vrf.SteeringMemberCandidatePubKey.Address())
		vrfsBitArray.SetIndex(i, vrfSet.VrfsBitArray.GetIndex(int(vrfIndex)))
	}

	vrfSet.Vrfs = vrfs[:]
	vrfSet.VrfsBitArray = vrfsBitArray.Copy()

	return nil
}

func (vrfSet *VrfSet) GetVrfIndexByAddress(address []byte) (int32) {
	if vrfSet == nil {
		return -1
	}
	for idx, vrf := range vrfSet.Vrfs {
		if bytes.Equal(vrf.SteeringMemberCandidatePubKey.Address(), address) {
			return int32(idx)
		}
	}
	return -1
}

func (vrfSet *VrfSet) HasAddress(address []byte) (bool) {
	if vrfSet == nil {
		return false
	}
	for _, vrf := range vrfSet.Vrfs {
		if bytes.Equal(vrf.SteeringMemberCandidatePubKey.Address(), address) {
			return true
		}
	}
	return false
}

func (vrfSet *VrfSet) GetSteeringMemberAddresses() *SettingSteeringMember {
	if len(vrfSet.Vrfs) != 0 {
		sort.Sort(SortedVrfs(vrfSet.Vrfs))
		var steeringMemberSize int
		if len(vrfSet.Vrfs) < MAXIMUM_STEERING_MEMBER_CANDIDATES {
			steeringMemberSize = len(vrfSet.Vrfs)
		} else {
			steeringMemberSize = MAXIMUM_STEERING_MEMBER_CANDIDATES
		}

		settingSteeringMember := NewSettingSteeringMember(steeringMemberSize)

		for i := 0; i < steeringMemberSize; i++ {
			if (vrfSet.Vrfs[i].Value != nil) {
				settingSteeringMember.SteeringMemberAddresses = append(settingSteeringMember.SteeringMemberAddresses, vrfSet.Vrfs[i].SteeringMemberCandidatePubKey.Address())
			}
		}

		fmt.Println("stompesi - settingSteeringMember", settingSteeringMember)
		return settingSteeringMember
	}
	return nil
}