package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/reapchain/reapchain-core/crypto"
	ce "github.com/reapchain/reapchain-core/crypto/encoding"
	tmbytes "github.com/reapchain/reapchain-core/libs/bytes"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

// Tx is an arbitrary byte array.
// NOTE: Tx has no types at this level, so when wire encoded it's just length-prefixed.
// Might we want types here ?
// type Qrns []Qrn

// func (qrns Qrns) Hash() []byte {
// 	bzs := make([][]byte, len(qrns))
// 	for i, qrn := range qrns {
// 		bzs[i] = qrn.Bytes()
// 	}
// 	return merkle.HashFromByteSlices(bzs)
// }
type QrnData struct {

	// Txs that will be applied by state @ block.Height+1.
	// NOTE: not all txs here are valid.  We're just agreeing on the order first.
	// This means that block.AppHash does not include these txs.
	Qrns QrnSet `json:"qrns"`

	// Volatile
	hash tmbytes.HexBytes
}

func (qrnData *QrnData) Hash() tmbytes.HexBytes {
	if qrnData == nil {
		return (&QrnSet{}).Hash()
	}
	if qrnData.hash == nil {
		qrnData.hash = qrnData.Qrns.Hash() // NOTE: leaves of merkle tree are TxIDs
	}
	return qrnData.hash
}

// -------------------------------------------------------------------
type Qrn struct {
	Height               int64         `json:"height"`
	Timestamp            time.Time     `json:"timestamp"`
	StandingMemberPubKey crypto.PubKey `json:"standing_member_pub_key"`
	StandingMemberIndex  int32         `json:"standing_member_index"`
	Value                uint64        `json:"value"`
	Signature            []byte        `json:"signature"`
}

func NewQrn(pubKey crypto.PubKey, value uint64, height int64, signature []byte) *Qrn {
	return &Qrn{
		StandingMemberPubKey: pubKey,
		Value:                value,
		Height:               height,
		Signature:            signature,
	}
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

	if qrn.StandingMemberIndex < 0 {
		return errors.New("negative StandingMemberIndex")
	}

	if len(qrn.Signature) == 0 {
		return errors.New("signature is missing")
	}

	if len(qrn.Signature) > MaxSignatureSize {
		return fmt.Errorf("signature is too big (max: %d)", MaxSignatureSize)
	}

	//TODO: stompesi -
	return nil
}

func (qrn *Qrn) Verify(pubKey crypto.PubKey) error {
	if !bytes.Equal(pubKey.Address(), qrn.StandingMemberPubKey.Address()) {
		return ErrQrnInvalidStandingMemberAddress
	}
	qrnProto, err := qrn.ToProto()
	if err != nil {
		return err
	}

	if !pubKey.VerifySignature(QrnValueToBytes(qrnProto.Value), qrn.Signature) {
		return ErrQrnInvalidSignature
	}
	return nil
}

func (qrn *Qrn) Bytes() []byte {
	pubKey, err := ce.PubKeyToProto(qrn.StandingMemberPubKey)
	if err != nil {
		panic(err)
	}

	pbv := tmproto.Qrn{
		Height:               qrn.Height,
		Timestamp:            qrn.Timestamp,
		StandingMemberPubKey: pubKey,
		StandingMemberIndex:  qrn.StandingMemberIndex,
		Value:                qrn.Value,
		Signature:            qrn.Signature,
	}

	bz, err := pbv.Marshal()
	if err != nil {
		panic(err)
	}
	return bz
}

func (qrn *Qrn) ToProto() (*tmproto.Qrn, error) {
	if qrn == nil {
		return nil, errors.New("nil qrn")
	}

	pubKey, err := ce.PubKeyToProto(qrn.StandingMemberPubKey)
	if err != nil {
		return nil, err
	}

	qrnProto := tmproto.Qrn{
		Height:               qrn.Height,
		Timestamp:            qrn.Timestamp,
		StandingMemberPubKey: pubKey,
		StandingMemberIndex:  qrn.StandingMemberIndex,
		Value:                qrn.Value,
		Signature:            qrn.Signature,
	}

	return &qrnProto, nil
}

func (qrn *Qrn) Copy() *Qrn {
	qrnCopy := *qrn
	return &qrnCopy
}

// TODO: Check
func QrnValueToBytes(value uint64) []byte {
	qrnValueBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(qrnValueBytes, value)

	return qrnValueBytes
}

func QrnListString(qrns []*Qrn) string {
	chunks := make([]string, len(qrns))
	for i, qrn := range qrns {
		chunks[i] = fmt.Sprintf("%s", qrn.StandingMemberPubKey.Address())
	}

	return strings.Join(chunks, ",")
}

func QrnFromProto(qrnProto *tmproto.Qrn) (*Qrn, error) {
	if qrnProto == nil {
		return nil, errors.New("nil qrn")
	}

	pubKey, err := ce.PubKeyFromProto(qrnProto.StandingMemberPubKey)
	if err != nil {
		return nil, err
	}

	qrn := new(Qrn)
	qrn.Height = qrnProto.Height
	qrn.Timestamp = qrnProto.Timestamp
	qrn.StandingMemberPubKey = pubKey
	qrn.StandingMemberIndex = qrnProto.StandingMemberIndex
	qrn.Value = qrnProto.Value
	qrn.Signature = qrnProto.Signature

	return qrn, nil
}
