package types

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/reapchain/reapchain-core/crypto/merkle"
	tmsync "github.com/reapchain/reapchain-core/libs/sync"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

var (
	ErrGotQrnFromUnwantedCosensusRound = errors.New(
		"peer has sent a qrn that does not match our height for more than one round",
	)
)

var (
	ErrQrnUnexpectedStep               = errors.New("unexpected step")
	ErrQrnInvalidStandingMemberIndex   = errors.New("invalid standing member index")
	ErrQrnInvalidStandingMemberAddress = errors.New("invalid standing member address")
	ErrQrnInvalidSignature             = errors.New("invalid signature")
	ErrQrnInvalidBlockHash             = errors.New("invalid block hash")
	ErrQrnNonDeterministicSignature    = errors.New("non-deterministic signature")
	ErrQrnNil                          = errors.New("nil qrn")
)

type QrnSet struct {
	Height            int64
	StandingMemberSet *StandingMemberSet

	Qrns []*Qrn

	mtx tmsync.Mutex
}

func NewQrnSet(height int64, standingMemberSet *StandingMemberSet, qrns []*Qrn) *QrnSet {
	return &QrnSet{
		Height:            height,
		StandingMemberSet: standingMemberSet,
		Qrns:              qrns,
	}
}

func (qrnSet *QrnSet) Size() int {
	if qrnSet == nil {
		return 0
	}
	return len(qrnSet.Qrns)
}

func (qrnSet *QrnSet) UpdateQrn(qrn *Qrn) (added bool, err error) {
	fmt.Println("한종빈4")
	if qrnSet == nil {
		panic("UpdateQrn() on nil QrnSet")
	}
	qrnSet.mtx.Lock()
	defer qrnSet.mtx.Unlock()

	return qrnSet.updateQrn(qrn)
}

func (qrnSet *QrnSet) updateQrn(qrn *Qrn) (added bool, err error) {
	fmt.Println("한종빈6")
	if qrn == nil {
		return false, ErrQrnNil
	}
	standingMemberIndex := qrn.StandingMemberIndex
	standingMemberAddr := qrn.StandingMemberPubKey.Address()

	fmt.Println("한종빈7")
	// Ensure that standing member index was set
	if standingMemberIndex < 0 {
		return false, fmt.Errorf("index < 0: %w", ErrQrnInvalidStandingMemberIndex)
	} else if len(standingMemberAddr) == 0 {
		return false, fmt.Errorf("empty address: %w", ErrQrnInvalidStandingMemberAddress)
	}

	// Make sure the step matches.
	fmt.Println("한종빈8", qrn.Height, qrnSet.Height)
	if qrn.Height != qrnSet.Height {
		return false, fmt.Errorf("expected %d, but got %d: %w", qrnSet.Height, qrn.Height, ErrQrnUnexpectedStep)
	}

	// Ensure that signer is a standing member.

	fmt.Println("한종빈9", qrnSet.StandingMemberSet, standingMemberIndex)
	lookupAddr, standingMember := qrnSet.StandingMemberSet.GetByIndex(standingMemberIndex)
	fmt.Println("한종빈9", lookupAddr, standingMember)
	if standingMember == nil {
		return false, fmt.Errorf(
			"cannot find standing member %d in standing member set of size %d: %w",
			standingMemberIndex, qrnSet.StandingMemberSet.Size(), ErrQrnInvalidStandingMemberIndex)
	}

	fmt.Println("한종빈10")
	// Ensure that the signer has the right address.
	if !bytes.Equal(standingMemberAddr, lookupAddr) {
		return false, fmt.Errorf(
			"qrn.StandingMemberAddress (%X) does not match address (%X) for qrn.StandingMemberIndex (%d)\n"+
				"Ensure the genesis file is correct across all standing members: %w",
			standingMemberAddr, lookupAddr, standingMemberIndex, ErrQrnInvalidStandingMemberAddress)
	}

	// Check signature.
	fmt.Println("한종빈11")
	if err := qrn.Verify(standingMember.PubKey); err != nil {
		return false, fmt.Errorf("failed to verify qrn and PubKey %s: %w", standingMember.PubKey, err)
	}

	// Add qrn
	fmt.Println("한종빈12")
	qrnSet.Qrns[standingMemberIndex] = qrn

	return added, nil
}

func (qrns *QrnSet) UpdateWithChangeSet(qrnz []*Qrn) error {
	return qrns.updateWithChangeSet(qrnz)
}

func (qrnSet *QrnSet) ValidateBasic() error {

	if qrnSet.IsNilOrEmpty() {
		return errors.New("qrn set is nil or empty")
	}

	for idx, qrn := range qrnSet.Qrns {
		if err := qrn.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid qrn #%d: %w", idx, err)
		}
	}

	return nil
}

func (qrns *QrnSet) IsNilOrEmpty() bool {
	return qrns == nil || len(qrns.Qrns) == 0
}

func (qrnSet *QrnSet) Hash() []byte {
	qrnBytesArray := make([][]byte, len(qrnSet.Qrns))
	for i, qrn := range qrnSet.Qrns {
		if qrn != nil {
			qrnBytesArray[i] = qrn.Bytes()
		}
	}
	return merkle.HashFromByteSlices(qrnBytesArray)
}

func QrnSetFromProto(vp *tmproto.QrnSet) (*QrnSet, error) {
	if vp == nil {
		return nil, errors.New("nil qrn set")
	}

	qrns := new(QrnSet)
	qrnsProto := make([]*Qrn, len(vp.Qrns))

	for i := 0; i < len(vp.Qrns); i++ {
		v, err := QrnFromProto(vp.Qrns[i])
		if err != nil {
			return nil, err
		}
		qrnsProto[i] = v
	}
	qrns.Qrns = qrnsProto

	return qrns, qrns.ValidateBasic()
}

func (qrnSet *QrnSet) ToProto() (*tmproto.QrnSet, error) {
	if qrnSet.IsNilOrEmpty() {
		return &tmproto.QrnSet{}, nil
	}

	qrnSetProto := new(tmproto.QrnSet)
	qrnsProto := make([]*tmproto.Qrn, len(qrnSet.Qrns))
	for i := 0; i < len(qrnSet.Qrns); i++ {
		if qrnSet.Qrns[i] != nil {
			valp, err := qrnSet.Qrns[i].ToProto()
			if err == nil {
				qrnsProto[i] = valp
			}
		}
	}
	qrnSetProto.Qrns = qrnsProto

	fmt.Println("qrnSetProto.Qrns", qrnSetProto.Qrns)
	return qrnSetProto, nil
}

func (qrns *QrnSet) Copy() *QrnSet {
	return &QrnSet{
		Height:            qrns.Height,
		StandingMemberSet: qrns.StandingMemberSet,
		Qrns:              qrnListCopy(qrns.Qrns),
	}
}

func qrnListCopy(qrns []*Qrn) []*Qrn {
	if qrns == nil {
		return nil
	}
	qrnsCopy := make([]*Qrn, len(qrns))
	for i, qrn := range qrns {
		if qrn != nil {
			qrnsCopy[i] = qrn.Copy()
		}
	}
	return qrnsCopy
}

// 정렬하기 위한 구조체
type QrnsByAddress []*Qrn

func (qrns QrnsByAddress) Len() int { return len(qrns) }

func (qrns QrnsByAddress) Less(i, j int) bool {
	return qrns[i].Value < qrns[j].Value
}

func (qrns QrnsByAddress) Swap(i, j int) {
	qrns[i], qrns[j] = qrns[j], qrns[i]
}

func (qrns *QrnSet) updateWithChangeSet(qrnz []*Qrn) error {
	if len(qrnz) != 0 {
		sort.Sort(QrnsByAddress(qrnz))
		qrns.Qrns = qrnz[:]
	}

	return nil
}
