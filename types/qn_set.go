package types

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	ce "gitlab.reappay.net/sucs-lab//reapchain/crypto/encoding"
	"gitlab.reappay.net/sucs-lab//reapchain/crypto/merkle"
	tmsync "gitlab.reappay.net/sucs-lab//reapchain/libs/sync"
	tmproto "gitlab.reappay.net/sucs-lab//reapchain/proto/reapchain/types"
)

var (
	ErrGotQnFromUnwantedCosensusRound = errors.New(
		"peer has sent a qn that does not match our height for more than one round",
	)
)

var (
	ErrQnUnexpectedStep               = errors.New("unexpected step")
	ErrQnInvalidStandingMemberIndex   = errors.New("invalid standing member index")
	ErrQnInvalidStandingMemberAddress = errors.New("invalid standing member address")
	ErrQnInvalidSignature             = errors.New("invalid signature")
	ErrQnInvalidBlockHash             = errors.New("invalid block hash")
	ErrQnNonDeterministicSignature    = errors.New("non-deterministic signature")
	ErrQnNil                          = errors.New("nil vote")
)

type QnSet struct {
	height            int64
	standingMemberSet *StandingMemberSet

	qns []*Qn

	mtx tmsync.Mutex
}

func NewQnSet(height int64, standingMemberSet *StandingMemberSet) *QnSet {
	if height == 0 {
		panic("Cannot make QnSet for height == 0, doesn't make sense.")
	}
	return &QnSet{
		height:            height,
		standingMemberSet: standingMemberSet,
		qns:               make([]*Qn, standingMemberSet.Size()),
	}
}

func (qnSet *QnSet) GetHeight() int64 {
	if qnSet == nil {
		return 0
	}
	return qnSet.height
}

func (qnSet *QnSet) Size() int {
	if qnSet == nil {
		return 0
	}
	return qnSet.standingMemberSet.Size()
}

func (qnSet *QnSet) AddQn(qn *Qn) (added bool, err error) {
	if qnSet == nil {
		panic("AddQn() on nil QnSet")
	}
	qnSet.mtx.Lock()
	defer qnSet.mtx.Unlock()

	return qnSet.addQn(qn)
}

func (qnSet *QnSet) addQn(qn *Qn) (added bool, err error) {
	if qn == nil {
		return false, ErrQnNil
	}
	standingMemberIndex := qn.StandingMemberIndex
	standingMemberAddr := qn.StandingMemberAddress

	// Ensure that standing member index was set
	if standingMemberIndex < 0 {
		return false, fmt.Errorf("index < 0: %w", ErrQnInvalidStandingMemberIndex)
	} else if len(standingMemberAddr) == 0 {
		return false, fmt.Errorf("empty address: %w", ErrQnInvalidStandingMemberAddress)
	}

	// Make sure the step matches.
	if qn.Height != qnSet.height {
		return false, fmt.Errorf("expected %d, but got %d: %w", qnSet.height, qn.Height)
	}

	// Ensure that signer is a standing member.
	lookupAddr, standingMember := qnSet.standingMemberSet.GetByIndex(standingMemberIndex)
	if standingMember == nil {
		return false, fmt.Errorf(
			"cannot find standing member %d in standing member set of size %d: %w",
			standingMemberIndex, qnSet.standingMemberSet.Size(), ErrQnInvalidStandingMemberIndex)
	}

	// Ensure that the signer has the right address.
	if !bytes.Equal(standingMemberAddr, lookupAddr) {
		return false, fmt.Errorf(
			"qn.StandingMemberAddress (%X) does not match address (%X) for qn.StandingMemberIndex (%d)\n"+
				"Ensure the genesis file is correct across all standing members: %w",
			standingMemberAddr, lookupAddr, standingMemberIndex, ErrQnInvalidStandingMemberAddress)
	}

	// Check signature.
	if err := qn.Verify(standingMember.PubKey); err != nil {
		return false, fmt.Errorf("failed to verify qn and PubKey %s: %w", standingMember.PubKey, err)
	}

	// Add qn
	qnSet.qns[standingMemberIndex] = qn

	return added, nil
}

func (qns *QnSet) UpdateWithChangeSet(qnz []*Qn) error {
	return qns.updateWithChangeSet(qnz)
}

func (qns *QnSet) ValidateBasic() error {
	if qns.IsNilOrEmpty() {
		return errors.New("qn set is nil or empty")
	}

	for idx, qn := range qns.Qns {
		if err := qn.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid qn #%d: %w", idx, err)
		}
	}

	return nil
}

func (qns *QnSet) IsNilOrEmpty() bool {
	return qns == nil || len(qns.Qns) == 0
}

func (qns *QnSet) Hash() []byte {
	bzs := make([][]byte, len(qns.Qns))
	for i, qn := range qns.Qns {
		bzs[i] = qn.Bytes()
	}
	return merkle.HashFromByteSlices(bzs)
}

func QnSetFromProto(vp *tmproto.QnSet) (*QnSet, error) {
	if vp == nil {
		return nil, errors.New("nil qn set")
	}
	qns := new(QnSet)

	qnsProto := make([]*Qn, len(vp.Qns))
	for i := 0; i < len(vp.Qns); i++ {
		v, err := QnFromProto(vp.Qns[i])
		if err != nil {
			return nil, err
		}
		qnsProto[i] = v
	}
	qns.Qns = qnsProto

	return qns, qns.ValidateBasic()
}

func (qns *QnSet) ToProto() (*tmproto.QnSet, error) {
	if qns.IsNilOrEmpty() {
		return &tmproto.QnSet{}, nil
	}

	vp := new(tmproto.QnSet)
	qnsProto := make([]*tmproto.Qn, len(qns.Qns))
	for i := 0; i < len(qns.Qns); i++ {
		valp := qns.Qns[i].ToProto()
		if valp != nil {
			qnsProto[i] = valp
		}
	}
	vp.Qns = qnsProto

	return vp, nil
}

func (qns *QnSet) Copy() *QnSet {
	return &QnSet{
		Qns: standingMemberListCopy(qns.Qns),
	}
}

func standingMemberListCopy(qns []*Qn) []*Qn {
	if qns == nil {
		return nil
	}
	qnsCopy := make([]*Qn, len(qns))
	for i, qn := range qns {
		qnsCopy[i] = qn.Copy()
	}
	return qnsCopy
}

// 정렬하기 위한 구조체
type QnsByAddress []*Qn

func (qns QnsByAddress) Len() int { return len(qns) }

func (qns QnsByAddress) Less(i, j int) bool {
	return bytes.Compare(qns[i].Address, qns[j].Address) == -1
}

func (qns QnsByAddress) Swap(i, j int) {
	qns[i], qns[j] = qns[j], qns[i]
}

func (qns *QnSet) updateWithChangeSet(qnz []*Qn) error {
	if len(qnz) != 0 {
		sort.Sort(QnsByAddress(qnz))
		qns.Qns = qnz[:]
	}

	return nil
}

func QnFromProto(vp *tmproto.Qn) (*Qn, error) {
	if vp == nil {
		return nil, errors.New("nil qn")
	}

	pk, err := ce.PubKeyFromProto(vp.PubKey)
	if err != nil {
		return nil, err
	}
	v := new(Qn)
	v.Address = vp.GetAddress()
	v.PubKey = pk
	v.Value = vp.Value
	v.Height = vp.Height

	return v, nil
}

func (qn *Qn) ToProto() *tmproto.Qn {
	if qn == nil {
		return nil
	}

	pk, err := ce.PubKeyToProto(qn.PubKey)
	if err != nil {
		return nil
	}

	qnProto := tmproto.Qn{
		Address: qn.Address,
		PubKey:  pk,
		Value:   qn.Value,
		Height:  qn.Height,
	}

	return &qnProto
}

func (qn *Qn) Copy() *Qn {
	qnCopy := *qn
	return &qnCopy
}
