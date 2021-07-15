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
	Height int64

	mtx          tmsync.Mutex
	Qrns         []*Qrn
	QrnsBitArray *bits.BitArray
}

func NewQrnSet(height int64, standingMemberSet *StandingMemberSet, qrns []*Qrn) *QrnSet {

	fmt.Println("NewQrnSet", len(standingMemberSet.StandingMembers))

	if qrns == nil {
		qrns = make([]*Qrn, standingMemberSet.Size())
		for i, standingMember := range standingMemberSet.StandingMembers {
			qrns[i] = NewQrnAsEmpty(standingMember.PubKey)
		}
	}

	fmt.Println("NewQrnSet-length(qrns)", len(qrns))
	return &QrnSet{
		Height:       height,
		QrnsBitArray: bits.NewBitArray(standingMemberSet.Size()),
		Qrns:         qrns,
	}
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

func (qrnSet *QrnSet) Size() int {
	if qrnSet == nil {
		return 0
	}
	return len(qrnSet.Qrns)
}

func (qrnSet *QrnSet) GetHeight() int64 {
	if qrnSet == nil {
		return 0
	}
	return qrnSet.Height
}

func (qrnSet *QrnSet) AddQrn(qrn *Qrn) (added bool, err error) {

	if qrnSet == nil {
		panic("AddQrn() on nil QrnSet")
	}
	qrnSet.mtx.Lock()
	defer qrnSet.mtx.Unlock()

	if qrn == nil {
		return false, ErrQrnNil
	}

	// Make sure the step matches.
	// TODO: stompesi
	// if qrn.Height+1 != qrnSet.Height {
	// 	return false, fmt.Errorf("expected %d, but got %d: %w", qrnSet.Height, qrn.Height, ErrQrnUnexpectedStep)
	// }

	// fmt.Println("------------------------------")
	// fmt.Println(qrn.Height)
	// fmt.Println(qrn.StandingMemberPubKey.Address())
	// fmt.Println(qrn.StandingMemberIndex)
	// fmt.Println(qrn.Value)
	// fmt.Println(qrn.Signature)
	// fmt.Println("------------------------------")
	// fmt.Println()

	if qrn.Height+1 != qrnSet.Height {
		return false, fmt.Errorf("expected")
	}

	fmt.Println("AddQrn",
		"qrnHeight", qrn.Height,
		"qrnSetHeight", qrnSet.Height,
		"qrn.StandingMemberIndex", qrn.StandingMemberIndex,
		"qrn.StandingMemberPubKey.Address", qrn.StandingMemberPubKey.Address(),
		"qrnSet.Qrns[qrn.StandingMemberIndex].StandingMemberPubKey.Address", qrnSet.Qrns[qrn.StandingMemberIndex].StandingMemberPubKey.Address())
	// Check signature.
	if err := qrn.Verify(qrn.StandingMemberPubKey); err != nil {
		return false, fmt.Errorf("failed to verify qrn and PubKey %s: %w", qrn.StandingMemberPubKey, err)
	}

	// Add qrn
	qrnSet.Qrns[qrn.StandingMemberIndex] = qrn
	qrnSet.QrnsBitArray.SetIndex(int(qrn.StandingMemberIndex), true)

	return true, nil
}

func (qrnSet *QrnSet) GetQrn(standingMemberPubKey crypto.PubKey) (qrn *Qrn) {
	for _, qrn := range qrnSet.Qrns {
		if qrn != nil {
			if standingMemberPubKey.Equals(qrn.StandingMemberPubKey) == true {
				return qrn
			}
		}
	}

	return nil
}

func (qrnSet *QrnSet) SortQrnSet() error {
	sort.Slice(qrnSet.Qrns, func(i, j int) bool {
		return bytes.Compare(qrnSet.Qrns[i].StandingMemberPubKey.Address(), qrnSet.Qrns[j].StandingMemberPubKey.Address()) > 0
	})

	sort.Slice(qrnSet.Qrns, func(i, j int) bool {
		return qrnSet.Qrns[i].Value < qrnSet.Qrns[j].Value
	})

	return nil
}

func (qrnSet *QrnSet) ValidateBasic() error {

	if qrnSet.IsNilOrEmpty() {
		return errors.New("qrn set is nil or empty")
	}

	for idx, qrn := range qrnSet.Qrns {
		fmt.Println("QrnSet-ValidateBasic",
			"qrn.Height", qrn.Height,
			"qrn.Value", qrn.Value,
			"qrn.qrn.StandingMemberPubKey.Address", qrn.StandingMemberPubKey.Address())
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

	return qrnSetProto, nil
}

func (qrns *QrnSet) Copy() *QrnSet {
	return &QrnSet{
		Height: qrns.Height,
		Qrns:   qrnListCopy(qrns.Qrns),
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

func (qrnSet *QrnSet) GetByIndex(standingMemberIndex int32) *Qrn {
	if qrnSet == nil {
		return nil
	}
	qrnSet.mtx.Lock()
	defer qrnSet.mtx.Unlock()
	return qrnSet.Qrns[standingMemberIndex]
}

func (qrnSet *QrnSet) BitArray() *bits.BitArray {
	if qrnSet == nil {
		return nil
	}
	qrnSet.mtx.Lock()
	defer qrnSet.mtx.Unlock()
	return qrnSet.QrnsBitArray.Copy()
}

const nilQrnSetString = "nil-QrnSet"

func (qrnSet *QrnSet) StringShort() string {
	if qrnSet == nil {
		return nilQrnSetString
	}
	qrnSet.mtx.Lock()
	defer qrnSet.mtx.Unlock()

	return fmt.Sprintf(`QrnSet{H:%v %v}`,
		qrnSet.Height, qrnSet.QrnsBitArray)
}

func (qrnSet *QrnSet) QrnStrings() []string {
	qrnSet.mtx.Lock()
	defer qrnSet.mtx.Unlock()
	return qrnSet.qrnStrings()
}

func (qrnSet *QrnSet) qrnStrings() []string {
	qrnStrings := make([]string, len(qrnSet.Qrns))
	for i, qrn := range qrnSet.Qrns {
		if qrn == nil {
			qrnStrings[i] = nilVoteStr
		} else {
			qrnStrings[i] = qrn.String()
		}
	}
	return qrnStrings
}

func (qrnSet *QrnSet) BitArrayString() string {
	qrnSet.mtx.Lock()
	defer qrnSet.mtx.Unlock()
	return qrnSet.bitArrayString()
}

func (qrnSet *QrnSet) bitArrayString() string {
	bAString := qrnSet.QrnsBitArray.String()
	return fmt.Sprintf("%s", bAString)
}
