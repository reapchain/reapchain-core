package types

import (
	"errors"
	"fmt"
	"time"

	abci "github.com/reapchain/reapchain-core/abci/types"
	"github.com/reapchain/reapchain-core/crypto"
	ce "github.com/reapchain/reapchain-core/crypto/encoding"
	tmbytes "github.com/reapchain/reapchain-core/libs/bytes"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain-core/types"
)

const (
	nilQrnStr string = "nil-Qrn"
)

type Qrn struct {
	Height               int64         `json:"height"`
	Timestamp            time.Time     `json:"timestamp"`
	StandingMemberPubKey crypto.PubKey `json:"standing_member_pub_key"`
	StandingMemberIndex  int32         `json:"standing_member_index"`
	Value                uint64        `json:"value"`
	Signature            []byte        `json:"signature"`
}

func NewQrn(height int64, standingMemberPubKey crypto.PubKey, value uint64) *Qrn {
	qrn := Qrn{
		Height:               height,
		Timestamp:            time.Now(),
		StandingMemberPubKey: standingMemberPubKey,
		Value:                value,
	}

	return &qrn
}

func NewQrnAsEmpty(height int64, standingMemberPubKey crypto.PubKey) *Qrn {
	qrn := Qrn{
		Height:               height,
		Timestamp:            time.Now(),
		StandingMemberPubKey: standingMemberPubKey,
	}

	return &qrn
}

func (qrn *Qrn) ValidateBasic() error {
	if qrn.Height < 0 {
		return errors.New("negative Height")
	}

	standingMemberAddress := qrn.StandingMemberPubKey.Address()
	if len(standingMemberAddress) != crypto.AddressSize {
		return fmt.Errorf("expected StandingMemberAddress size to be %d bytes, got %d bytes",
			crypto.AddressSize,
			len(standingMemberAddress),
		)
	}

	return nil
}

func (qrn *Qrn) Copy() *Qrn {
	qrnCopy := *qrn
	return &qrnCopy
}

func (qrn *Qrn) String() string {
	if qrn == nil {
		return nilQrnStr
	}
	return fmt.Sprintf("Qrn{%X %v %X @ %s}",
		tmbytes.Fingerprint(qrn.StandingMemberPubKey.Address()),
		qrn.Height,
		tmbytes.Fingerprint(qrn.Signature),
		CanonicalTime(qrn.Timestamp),
	)
}

func (qrn *Qrn) GetQrnBytesForSign() []byte {
	if qrn == nil {
		return nil
	}

	pubKeyProto, err := ce.PubKeyToProto(qrn.StandingMemberPubKey)
	if err != nil {
		panic(err)
	}

	qrnProto := tmproto.Qrn{
		Height:               qrn.Height,
		Timestamp:            qrn.Timestamp,
		StandingMemberPubKey: pubKeyProto,
		Value:                qrn.Value,
	}

	QrnSignBytes, err := qrnProto.Marshal()
	if err != nil {
		panic(err)
	}
	return QrnSignBytes
}

func (qrn *Qrn) GetQrnBytes() []byte {
	if qrn == nil {
		return nil
	}

	pubKeyProto, err := ce.PubKeyToProto(qrn.StandingMemberPubKey)
	if err != nil {
		panic(err)
	}

	qrnProto := tmproto.Qrn{
		Height:               qrn.Height,
		Timestamp:            qrn.Timestamp,
		StandingMemberPubKey: pubKeyProto,
		Value:                qrn.Value,
		Signature:            qrn.Signature,
	}

	qrnSignBytes, err := qrnProto.Marshal()
	if err != nil {
		panic(err)
	}
	return qrnSignBytes
}

func (qrn *Qrn) VerifySign() bool {
	signBytes := qrn.GetQrnBytesForSign()
	if signBytes == nil {
		return false
	}

	return qrn.StandingMemberPubKey.VerifySignature(signBytes, qrn.Signature)
}

func (qrn *Qrn) ToProto() *tmproto.Qrn {
	if qrn == nil {
		return nil
	}

	pubKey, err := ce.PubKeyToProto(qrn.StandingMemberPubKey)
	if err != nil {
		return nil
	}

	qrnProto := tmproto.Qrn{
		Height:               qrn.Height,
		Timestamp:            qrn.Timestamp,
		StandingMemberPubKey: pubKey,
		StandingMemberIndex:  qrn.StandingMemberIndex,
		Value:                qrn.Value,
		Signature:            qrn.Signature,
	}

	return &qrnProto
}

func QrnFromProto(qrnProto *tmproto.Qrn) *Qrn {
	if qrnProto == nil {
		return nil
	}

	pubKey, err := ce.PubKeyFromProto(qrnProto.StandingMemberPubKey)
	if err != nil {
		return nil
	}

	qrn := new(Qrn)
	qrn.Height = qrnProto.Height
	qrn.Timestamp = qrnProto.Timestamp
	qrn.StandingMemberPubKey = pubKey
	qrn.StandingMemberIndex = qrnProto.StandingMemberIndex
	qrn.Value = qrnProto.Value
	qrn.Signature = qrnProto.Signature

	return qrn
}

func QrnFromAbci(qrnUpdate *abci.QrnUpdate) *Qrn {
	if qrnUpdate == nil {
		return nil
	}

	pubKey, err := ce.PubKeyFromProto(qrnUpdate.StandingMemberPubKey)
	if err != nil {
		return nil
	}

	qrn := new(Qrn)
	qrn.Height = qrnUpdate.Height
	qrn.Timestamp = qrnUpdate.Timestamp
	qrn.StandingMemberPubKey = pubKey
	// TODO: qrn.StandingMemberIndex = qrnUpdate.StandingMemberIndex
	qrn.Value = qrnUpdate.Value
	qrn.Signature = qrnUpdate.Signature

	return qrn
}
