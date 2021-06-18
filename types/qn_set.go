package types

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	ce "github.com/reapchain/reapchain/crypto/encoding"
	"github.com/reapchain/reapchain/crypto/merkle"
	tmproto "github.com/reapchain/reapchain/proto/reapchain/types"
)

type QnSet struct {
	Qns []*Qn `json:"qns"`
}

func NewQnSet(qnz []*Qn) *QnSet {
	qns := &QnSet{}

	err := qns.updateWithChangeSet(qnz)
	if err != nil {
		panic(fmt.Sprintf("Cannot create qn set: %v", err))
	}

	return qns
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
		valp, err := qns.Qns[i].ToProto()
		if err != nil {
			return nil, err
		}
		qnsProto[i] = valp
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

func (qns *QnSet) Size() int {
	return len(qns.Qns)
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

	return v, nil
}

func (qn *Qn) ToProto() (*tmproto.Qn, error) {
	if qn == nil {
		return nil, errors.New("nil qn")
	}

	pk, err := ce.PubKeyToProto(qn.PubKey)
	if err != nil {
		return nil, err
	}

	vp := tmproto.Qn{
		Address: qn.Address,
		PubKey:  pk,
		Value:   qn.Value,
	}

	return &vp, nil
}

func (qn *Qn) Copy() *Qn {
	qnCopy := *qn
	return &qnCopy
}
