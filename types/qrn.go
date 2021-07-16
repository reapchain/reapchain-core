package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/reapchain/reapchain-core/crypto"
	ce "github.com/reapchain/reapchain-core/crypto/encoding"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

type Qrn struct {
	Height               int64         `json:"height"`
	Timestamp            time.Time     `json:"timestamp"`
	StandingMemberPubKey crypto.PubKey `json:"standing_member_pub_key"`
	Value                uint64        `json:"value"`
	Signature            []byte        `json:"signature"`
}

func NewQrn(height int64, standingMemberPubKey crypto.PubKey, value uint64) Qrn {
	qrn := Qrn{
		Height:               height,
		Timestamp:            time.Now(),
		StandingMemberPubKey: standingMemberPubKey,
		Value:                value,
	}

	return qrn
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

func (qrn *Qrn) GetQrnSignBytes() []byte {
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

func (qrn *Qrn) VerifySign() bool {
	signBytes := qrn.GetQrnSignBytes()
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
	qrn.Value = qrnProto.Value
	qrn.Signature = qrnProto.Signature

	return &Qrn{
		Height:               qrnProto.Height,
		Timestamp:            qrnProto.Timestamp,
		StandingMemberPubKey: pubKey,
		Value:                qrnProto.Value,
		Signature:            qrnProto.Signature,
	}
}
