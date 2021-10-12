package types

import (
	encoding_binary "encoding/binary"
	"errors"
	"fmt"

	"github.com/reapchain/reapchain-core/crypto"
	"github.com/reapchain/reapchain-core/crypto/merkle"
	"github.com/reapchain/reapchain-core/libs/bits"
	tmsync "github.com/reapchain/reapchain-core/libs/sync"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

type QrnSet struct {
	Height            int64
	StandingMemberSet *StandingMemberSet
	mtx               tmsync.Mutex
	Qrns              []*Qrn
	QrnsBitArray      *bits.BitArray
}

func NewQrnSet(height int64, standingMemberSet *StandingMemberSet, qrns []*Qrn) *QrnSet {
	if qrns == nil {
		qrns = make([]*Qrn, standingMemberSet.Size())
		for i, standingMember := range standingMemberSet.StandingMembers {
			qrns[i] = NewQrnAsEmpty(height, standingMember.PubKey)
		}
	} else {
		for i, standingMember := range standingMemberSet.StandingMembers {
			standingMemberIndex, _ := standingMemberSet.GetStandingMemberByAddress(standingMember.Address)
			qrns[i].StandingMemberIndex = standingMemberIndex
		}
	}

	return &QrnSet{
		Height:            height,
		StandingMemberSet: standingMemberSet,
		QrnsBitArray:      bits.NewBitArray(standingMemberSet.Size()),
		Qrns:              qrns,
	}
}

func QrnSetFromProto(qrnSetProto *tmproto.QrnSet, standingMemberSetProto *tmproto.StandingMemberSet) (*QrnSet, error) {
	if qrnSetProto == nil {
		return nil, errors.New("nil qrn set")
	}

	qrnSet := new(QrnSet)
	qrns := make([]*Qrn, len(qrnSetProto.Qrns))

	for i, qrnProto := range qrnSetProto.Qrns {
		qrn := QrnFromProto(qrnProto)
		if qrn == nil {
			return nil, errors.New(nilQrnStr)
		}
		qrns[i] = qrn
	}
	qrnSet.Height = qrnSetProto.Height
	standingMemberSet, err := StandingMemberSetFromProto(standingMemberSetProto)
	if err != nil {
		return nil, err
	}
	qrnSet.StandingMemberSet = standingMemberSet
	qrnSet.Qrns = qrns

	return qrnSet, qrnSet.ValidateBasic()
}

func (qrnSet *QrnSet) ValidateBasic() error {
	if qrnSet == nil || len(qrnSet.Qrns) == 0 {
		return errors.New("qrn set is nil or empty")
	}

	for idx, qrn := range qrnSet.Qrns {
		if err := qrn.ValidateBasic(); err != nil {
			return fmt.Errorf("Invalid qrn #%d: %w", idx, err)
		}
	}

	return nil
}

func (qrnSet *QrnSet) GetHeight() int64 {
	if qrnSet == nil {
		return 0
	}
	return qrnSet.Height
}

func (qrnSet *QrnSet) Size() int {
	if qrnSet == nil {
		return 0
	}

	return len(qrnSet.Qrns)
}

func (qrnSet *QrnSet) IsNilOrEmpty() bool {
	return qrnSet == nil || len(qrnSet.Qrns) == 0
}

func (qrnSet *QrnSet) GetMaxValue() uint64 {
	maxValue := uint64(0)
	qrnSetHash := qrnSet.Hash()

	for _, qrn := range qrnSet.Qrns {
		qrnHash := make([][]byte, 2)
		qrnHash[0] = qrnSetHash
		qrnHash[1] = qrn.GetQrnBytes()

		result := merkle.HashFromByteSlices(qrnHash)
		if uint64(maxValue) < encoding_binary.LittleEndian.Uint64(result) {
			maxValue = encoding_binary.LittleEndian.Uint64(result)
		}
	}

	return maxValue
}

func (qrnSet *QrnSet) AddQrn(qrn *Qrn) error {
	qrnSet.mtx.Lock()
	defer qrnSet.mtx.Unlock()

	if qrn == nil {
		return fmt.Errorf("Qrn is nil")
	}

	if qrn.VerifySign() == false {
		return fmt.Errorf("Invalid qrn sign")
	}
	standingMemberIndex, _ := qrnSet.StandingMemberSet.GetStandingMemberByAddress(qrn.StandingMemberPubKey.Address())

	if standingMemberIndex == -1 {
		return fmt.Errorf("Not exist standing member of qrn: %v", qrn.StandingMemberPubKey.Address())
	}

	qrn.StandingMemberIndex = standingMemberIndex
	qrnSet.Qrns[standingMemberIndex] = qrn.Copy()

	qrnSet.QrnsBitArray.SetIndex(int(standingMemberIndex), true)

	return nil
}

func (qrnSet *QrnSet) GetQrn(standingMemberPubKey crypto.PubKey) (qrn *Qrn) {
	standingMemberIndex, _ := qrnSet.StandingMemberSet.GetStandingMemberByAddress(standingMemberPubKey.Address())

	if standingMemberIndex != -1 {
		return qrnSet.Qrns[standingMemberIndex]
	}

	return nil
}

func (qrnSet *QrnSet) Hash() []byte {
	qrnBytesArray := make([][]byte, len(qrnSet.Qrns))
	for i, qrn := range qrnSet.Qrns {
		if qrn != nil {
			qrnBytesArray[i] = qrn.GetQrnBytes()
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
	for i, qrn := range qrnSet.Qrns {
		qrnProto := qrn.ToProto()
		if qrnProto != nil {
			qrnsProto[i] = qrnProto
		}
	}
	qrnSetProto.Height = qrnSet.Height
	qrnSetProto.Qrns = qrnsProto

	return qrnSetProto, nil
}

func (qrnSet *QrnSet) Copy() *QrnSet {
	qrnSet.mtx.Lock()
	defer qrnSet.mtx.Unlock()

	qrnsCopy := make([]*Qrn, len(qrnSet.Qrns))
	for i, qrn := range qrnSet.Qrns {
		if qrn != nil {
			qrnsCopy[i] = qrn.Copy()
		}
	}

	return &QrnSet{
		Height:            qrnSet.Height,
		StandingMemberSet: qrnSet.StandingMemberSet.Copy(),
		Qrns:              qrnsCopy,
		QrnsBitArray:      qrnSet.QrnsBitArray.Copy(),
	}
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

// -----
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

	qrnStrings := make([]string, len(qrnSet.Qrns))
	for i, qrn := range qrnSet.Qrns {
		if qrn == nil {
			qrnStrings[i] = "nil-Qrn"
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

type QrnSetReader interface {
	GetHeight() int64
	Size() int
	BitArray() *bits.BitArray
	GetByIndex(int32) *Qrn
}

func (qrnSet *QrnSet) UpdateWithChangeSet(qrns []*Qrn) error {
	return qrnSet.updateWithChangeSet(qrns)
}

func (qrnSet *QrnSet) updateWithChangeSet(qrns []*Qrn) error {
	if len(qrns) != 0 {
		qrnSet.Qrns = qrns[:]
	}

	return nil
}
