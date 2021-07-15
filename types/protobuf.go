package types

import (
	abci "github.com/reapchain/reapchain-core/abci/types"
	"github.com/reapchain/reapchain-core/crypto"
	"github.com/reapchain/reapchain-core/crypto/ed25519"
	cryptoenc "github.com/reapchain/reapchain-core/crypto/encoding"
	"github.com/reapchain/reapchain-core/crypto/secp256k1"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

//-------------------------------------------------------
// Use strings to distinguish types in ABCI messages

const (
	ABCIPubKeyTypeEd25519   = ed25519.KeyType
	ABCIPubKeyTypeSecp256k1 = secp256k1.KeyType
)

// TODO: Make non-global by allowing for registration of more pubkey types

var ABCIPubKeyTypesToNames = map[string]string{
	ABCIPubKeyTypeEd25519:   ed25519.PubKeyName,
	ABCIPubKeyTypeSecp256k1: secp256k1.PubKeyName,
}

//-------------------------------------------------------

// TM2PB is used for converting Reapchain ABCI to protobuf ABCI.
// UNSTABLE
var TM2PB = tm2pb{}

type tm2pb struct{}

func (tm2pb) Header(header *Header) tmproto.Header {
	return tmproto.Header{
		Version: header.Version,
		ChainID: header.ChainID,
		Height:  header.Height,
		Time:    header.Time,

		LastBlockId: header.LastBlockID.ToProto(),

		LastCommitHash: header.LastCommitHash,
		DataHash:       header.DataHash,

		ValidatorsHash:     header.ValidatorsHash,
		NextValidatorsHash: header.NextValidatorsHash,
		ConsensusHash:      header.ConsensusHash,
		AppHash:            header.AppHash,
		LastResultsHash:    header.LastResultsHash,

		EvidenceHash:    header.EvidenceHash,
		ProposerAddress: header.ProposerAddress,

		StandingMembersHash: header.StandingMembersHash,
		ConsensusRound:      header.ConsensusRound.ToProto(),
		QrnsHash:            header.QrnsHash,
	}
}

func (tm2pb) Validator(val *Validator) abci.Validator {
	return abci.Validator{
		Address: val.PubKey.Address(),
		Power:   val.VotingPower,
	}
}

func (tm2pb) StandingMember(val *StandingMember) abci.StandingMember {
	return abci.StandingMember{
		Address: val.PubKey.Address(),
	}
}

func (tm2pb) Qrn(val *Qrn) abci.Qrn {
	return abci.Qrn{}
}

func (tm2pb) BlockID(blockID BlockID) tmproto.BlockID {
	return tmproto.BlockID{
		Hash:          blockID.Hash,
		PartSetHeader: TM2PB.PartSetHeader(blockID.PartSetHeader),
	}
}

func (tm2pb) PartSetHeader(header PartSetHeader) tmproto.PartSetHeader {
	return tmproto.PartSetHeader{
		Total: header.Total,
		Hash:  header.Hash,
	}
}

// XXX: panics on unknown pubkey type
func (tm2pb) ValidatorUpdate(val *Validator) abci.ValidatorUpdate {
	pk, err := cryptoenc.PubKeyToProto(val.PubKey)
	if err != nil {
		panic(err)
	}
	return abci.ValidatorUpdate{
		PubKey: pk,
		Power:  val.VotingPower,
	}
}

// XXX: panics on nil or unknown pubkey type
func (tm2pb) ValidatorUpdates(vals *ValidatorSet) []abci.ValidatorUpdate {
	validators := make([]abci.ValidatorUpdate, vals.Size())
	for i, val := range vals.Validators {
		validators[i] = TM2PB.ValidatorUpdate(val)
	}
	return validators
}

func (tm2pb) StandingMemberUpdate(val *StandingMember) abci.StandingMemberUpdate {
	pk, err := cryptoenc.PubKeyToProto(val.PubKey)
	if err != nil {
		panic(err)
	}
	return abci.StandingMemberUpdate{
		PubKey: pk,
	}
}

func (tm2pb) StandingMemberUpdates(sms *StandingMemberSet) []abci.StandingMemberUpdate {
	standingMembers := make([]abci.StandingMemberUpdate, sms.Size())
	for i, sm := range sms.StandingMembers {
		standingMembers[i] = TM2PB.StandingMemberUpdate(sm)
	}
	return standingMembers
}

func (tm2pb) QrnUpdate(qrn *Qrn) abci.QrnUpdate {
	pk, err := cryptoenc.PubKeyToProto(qrn.StandingMemberPubKey)
	if err != nil {
		panic(err)
	}
	return abci.QrnUpdate{
		StandingMemberPubKey: pk,
		Value:                qrn.Value,
		Height:               qrn.Height,
		Signature:            qrn.Signature,
	}
}

func (tm2pb) QrnUpdates(qrnSet *QrnSet) []abci.QrnUpdate {
	qrns := make([]abci.QrnUpdate, qrnSet.Size())

	for i, sm := range qrnSet.Qrns {
		qrns[i] = TM2PB.QrnUpdate(sm)
	}
	return qrns
}

func (tm2pb) ConsensusParams(params *tmproto.ConsensusParams) *abci.ConsensusParams {
	return &abci.ConsensusParams{
		Block: &abci.BlockParams{
			MaxBytes: params.Block.MaxBytes,
			MaxGas:   params.Block.MaxGas,
		},
		Evidence:  &params.Evidence,
		Validator: &params.Validator,
	}
}

// XXX: panics on nil or unknown pubkey type
func (tm2pb) NewValidatorUpdate(pubkey crypto.PubKey, power int64) abci.ValidatorUpdate {
	pubkeyABCI, err := cryptoenc.PubKeyToProto(pubkey)
	if err != nil {
		panic(err)
	}
	return abci.ValidatorUpdate{
		PubKey: pubkeyABCI,
		Power:  power,
	}
}

//----------------------------------------------------------------------------

// PB2TM is used for converting protobuf ABCI to Reapchain ABCI.
// UNSTABLE
var PB2TM = pb2tm{}

type pb2tm struct{}

func (pb2tm) ValidatorUpdates(vals []abci.ValidatorUpdate) ([]*Validator, error) {
	tmVals := make([]*Validator, len(vals))
	for i, v := range vals {
		pub, err := cryptoenc.PubKeyFromProto(v.PubKey)
		if err != nil {
			return nil, err
		}
		tmVals[i] = NewValidator(pub, v.Power)
	}
	return tmVals, nil
}

func (pb2tm) StandingMemberUpdates(sms []abci.StandingMemberUpdate) ([]*StandingMember, error) {
	smz := make([]*StandingMember, len(sms))
	for i, v := range sms {
		pub, err := cryptoenc.PubKeyFromProto(v.PubKey)
		if err != nil {
			return nil, err
		}
		smz[i] = NewStandingMember(pub)
	}
	return smz, nil
}

func (pb2tm) QrnUpdates(qrnUpdates []abci.QrnUpdate) ([]*Qrn, error) {
	qrns := make([]*Qrn, len(qrnUpdates))
	for i, qrn := range qrnUpdates {
		pubKey, err := cryptoenc.PubKeyFromProto(qrn.StandingMemberPubKey)
		if err != nil {
			return nil, err
		}
		qrns[i] = NewQrn(pubKey, qrn.Value, qrn.Height, qrn.Signature)
	}
	return qrns, nil
}
