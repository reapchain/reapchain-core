package types

import (
	"bytes"
	encoding_binary "encoding/binary"
	"errors"
	"fmt"

	"github.com/reapchain/reapchain-core/crypto"
	"github.com/reapchain/reapchain-core/crypto/merkle"
	"github.com/reapchain/reapchain-core/libs/bits"
	tmsync "github.com/reapchain/reapchain-core/libs/sync"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain-core/types"
)

type QrnSet struct {
	Height            int64
	mtx               tmsync.Mutex
	Qrns              []*Qrn
	QrnsBitArray      *bits.BitArray
}

func NewQrnSet(height int64, standingMemberSet *StandingMemberSet, qrns []*Qrn) *QrnSet {
	if qrns == nil || len(qrns) == 0 {
		qrns = make([]*Qrn, standingMemberSet.Size())
		for i, standingMember := range standingMemberSet.StandingMembers {
			qrns[i] = NewQrnAsEmpty(height, standingMember.PubKey)
			qrns[i].QrnIndex = int32(i)
		}
	}

	return &QrnSet{
		Height:            height,
		QrnsBitArray:      bits.NewBitArray(standingMemberSet.Size()),
		Qrns:              qrns,
	}
}

func QrnSetFromProto(qrnSetProto *tmproto.QrnSet) (*QrnSet, error) {
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
	qrnSet.QrnsBitArray = bits.NewBitArray(len(qrns))
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

// Add qrn in the set
func (qrnSet *QrnSet) AddQrn(chainID string, qrn *Qrn) bool {
	qrnSet.mtx.Lock()
	defer qrnSet.mtx.Unlock()

	if qrn == nil {
		return false
	}
	
	if qrn.Value != 0 {
		if qrn.VerifySign(chainID) == false {
			return false
		}
	
		if qrnSet.Height != qrn.Height {
			return false
		}
	}

	qrnIndex := qrnSet.GetQrnIndexByAddress(qrn.StandingMemberPubKey.Address())

	if qrnIndex == -1 {
		return false
	}

	if qrnSet.QrnsBitArray.GetIndex(int(qrnIndex)) == false {
		qrn.QrnIndex = qrnIndex
		qrnSet.Qrns[qrnIndex] = qrn.Copy()
		qrnSet.QrnsBitArray.SetIndex(int(qrnIndex), true)
		return true
	}
	return false
}


func (qrnSet *QrnSet) GetQrn(standingMemberPubKey crypto.PubKey) (qrn *Qrn) {
	qrnIndex := qrnSet.GetQrnIndexByAddress(standingMemberPubKey.Address())

	if qrnIndex != -1 {
		return qrnSet.Qrns[qrnIndex]
	}

	return nil
}

func (qrnSet *QrnSet) Hash() []byte {
	qrnBytesArray := make([][]byte, len(qrnSet.Qrns))
	for i, qrn := range qrnSet.Qrns {
		if qrn != nil && qrn.Signature != nil {
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
	if qrnSet == nil {
		return nil
	}

	qrnsCopy := make([]*Qrn, len(qrnSet.Qrns))
	for i, qrn := range qrnSet.Qrns {
		if qrn != nil {
			qrnsCopy[i] = qrn.Copy()
		}
	}

	return &QrnSet{
		Height:            qrnSet.Height,
		Qrns:              qrnsCopy,
		QrnsBitArray:      qrnSet.QrnsBitArray.Copy(),
	}
}

func (qrnSet *QrnSet) GetByIndex(qrnIndex int32) *Qrn {
	if qrnSet == nil {
		return nil
	}

	return qrnSet.Qrns[qrnIndex]
}

func (qrnSet *QrnSet) BitArray() *bits.BitArray {
	if qrnSet == nil {
		return nil
	}

	qrnSet.mtx.Lock()
	defer qrnSet.mtx.Unlock()
	return qrnSet.QrnsBitArray.Copy()
}

type QrnSetReader interface {
	GetHeight() int64
	Size() int
	BitArray() *bits.BitArray
	GetByIndex(int32) *Qrn
}

func (qrnSet *QrnSet) UpdateWithChangeSet(standingMemberSet *StandingMemberSet) error {
	qrnSet.mtx.Lock()
	defer qrnSet.mtx.Unlock()

	qrns := make([]*Qrn, 0, len(standingMemberSet.StandingMembers))

	for _, standingMember := range standingMemberSet.StandingMembers {
		qrn := qrnSet.GetQrn(standingMember.PubKey)

		if qrn != nil {
			qrns = append(qrns, qrn.Copy())
		}
	}

	qrnsBitArray := bits.NewBitArray(len(qrns))
	
	for i, qrn := range qrns {
		qrnIndex := qrnSet.GetQrnIndexByAddress(qrn.StandingMemberPubKey.Address())
		qrnsBitArray.SetIndex(i, qrnSet.QrnsBitArray.GetIndex(int(qrnIndex)))
	}

	qrnSet.Qrns = qrns[:]
	qrnSet.QrnsBitArray = qrnsBitArray.Copy()

	return nil
}

func (qrnSet *QrnSet) GetQrnIndexByAddress(address []byte) (int32) {
	if qrnSet == nil {
		return -1
	}
	for idx, qrn := range qrnSet.Qrns {
		if bytes.Equal(qrn.StandingMemberPubKey.Address(), address) {
			return int32(idx)
		}
	}
	return -1
}

func (qrnSet *QrnSet) HasAddress(address []byte) (bool) {
	if qrnSet == nil {
		return false
	}
	for _, qrn := range qrnSet.Qrns {
		if bytes.Equal(qrn.StandingMemberPubKey.Address(), address) {
			return true
		}
	}
	return false
}

func QrnSetFromExistingQrns(qrns []*Qrn) (*QrnSet, error) {
	if len(qrns) == 0 {
		return nil, errors.New("qrn set is empty")
	}
	for _, val := range qrns {
		err := val.ValidateBasic()
		if err != nil {
			return nil, fmt.Errorf("can't create qrn set: %w", err)
		}
	}

	qrnSet := &QrnSet{
		Qrns: qrns,
	}
	return qrnSet, nil
}