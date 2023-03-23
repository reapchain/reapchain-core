package types

import (
	abci "github.com/reapchain/reapchain-core/abci/types"
	"github.com/reapchain/reapchain-core/crypto"
	"github.com/reapchain/reapchain-core/crypto/ed25519"
	cryptoenc "github.com/reapchain/reapchain-core/crypto/encoding"
	"github.com/reapchain/reapchain-core/crypto/secp256k1"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain-core/types"
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

// TM2PB is used for converting ReapchainCore ABCI to protobuf ABCI.
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

		StandingMembersHash:          header.StandingMembersHash,
		SteeringMemberCandidatesHash: header.SteeringMemberCandidatesHash,
		ConsensusRound:               header.ConsensusRound.ToProto(),
	}
}

func (tm2pb) Validator(val *Validator) abci.Validator {
	return abci.Validator{
		Address: val.PubKey.Address(),
		Power:   val.VotingPower,
	}
}

func (tm2pb) StandingMember(standingMember *StandingMember) abci.StandingMember {
	return abci.StandingMember{
		Address: standingMember.PubKey.Address(),
	}
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
		Type: val.Type,
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


func (tm2pb) StandingMemberUpdate(standingMember *StandingMember) abci.StandingMemberUpdate {
	pubKeyProto, err := cryptoenc.PubKeyToProto(standingMember.PubKey)
	if err != nil {
		panic(err)
	}
	return abci.StandingMemberUpdate{
		PubKey: pubKeyProto,
	}
}

func (tm2pb) SteeringMemberCandidateUpdate(steeringMemberCandidate *SteeringMemberCandidate) abci.SteeringMemberCandidateUpdate {
	pubKeyProto, err := cryptoenc.PubKeyToProto(steeringMemberCandidate.PubKey)
	if err != nil {
		panic(err)
	}
	return abci.SteeringMemberCandidateUpdate{
		PubKey: pubKeyProto,
	}
}

func (tm2pb) StandingMemberSetUpdate(standingMemberSet *StandingMemberSet) []abci.StandingMemberUpdate {
	standingMembers := make([]abci.StandingMemberUpdate, standingMemberSet.Size())
	for i, sm := range standingMemberSet.StandingMembers {
		standingMembers[i] = TM2PB.StandingMemberUpdate(sm)
	}
	return standingMembers
}

func (tm2pb) SteeringMemberCandidateSetUpdate(steeringMemberCandidateSet *SteeringMemberCandidateSet) []abci.SteeringMemberCandidateUpdate {
	steeringMemberCandidates := make([]abci.SteeringMemberCandidateUpdate, steeringMemberCandidateSet.Size())
	for i, sm := range steeringMemberCandidateSet.SteeringMemberCandidates {
		steeringMemberCandidates[i] = TM2PB.SteeringMemberCandidateUpdate(sm)
	}
	return steeringMemberCandidates
}

func (tm2pb) QrnUpdate(qrn *Qrn) abci.QrnUpdate {
	pubKeyProto, err := cryptoenc.PubKeyToProto(qrn.StandingMemberPubKey)
	if err != nil {
		panic(err)
	}
	return abci.QrnUpdate{
		Height:               qrn.Height,
		Timestamp:            qrn.Timestamp,
		StandingMemberPubKey: pubKeyProto,
		Value:                qrn.Value,
		Signature:            qrn.Signature,
	}
}

func (tm2pb) QrnSetUpdate(qrnSet *QrnSet) []abci.QrnUpdate {
	qrnUpdates := make([]abci.QrnUpdate, qrnSet.Size())
	for i, qrn := range qrnSet.Qrns {
		qrnUpdates[i] = TM2PB.QrnUpdate(qrn)
	}
	return qrnUpdates
}

func (tm2pb) VrfUpdate(vrf *Vrf) abci.VrfUpdate {
	pubKeyProto, err := cryptoenc.PubKeyToProto(vrf.SteeringMemberCandidatePubKey)
	if err != nil {
		panic(err)
	}
	return abci.VrfUpdate{
		Height:                        vrf.Height,
		Timestamp:                     vrf.Timestamp,
		SteeringMemberCandidatePubKey: pubKeyProto,
		Value:                         vrf.Value,
		Proof:                         vrf.Proof,
	}
}

func (tm2pb) VrfSetUpdate(vrfSet *VrfSet) []abci.VrfUpdate {
	vrfUpdates := make([]abci.VrfUpdate, vrfSet.Size())
	for i, vrf := range vrfSet.Vrfs {
		vrfUpdates[i] = TM2PB.VrfUpdate(vrf)
	}
	return vrfUpdates
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

func (tm2pb) ConsensusRound(consensudRoundProto *tmproto.ConsensusRound) *abci.ConsensusRound {
	return &abci.ConsensusRound{
		ConsensusStartBlockHeight: consensudRoundProto.ConsensusStartBlockHeight,
		QrnPeriod:                 consensudRoundProto.QrnPeriod,
		VrfPeriod:                 consensudRoundProto.VrfPeriod,
		ValidatorPeriod:           consensudRoundProto.ValidatorPeriod,
		Period:                    consensudRoundProto.Period,
	}
}

// XXX: panics on nil or unknown pubkey type
func (tm2pb) NewValidatorUpdate(pubkey crypto.PubKey, power int64, validatorType string) abci.ValidatorUpdate {
	pubkeyABCI, err := cryptoenc.PubKeyToProto(pubkey)
	if err != nil {
		panic(err)
	}
	return abci.ValidatorUpdate{
		PubKey: pubkeyABCI,
		Power:  power,
		Type: validatorType,
	}
}

//----------------------------------------------------------------------------

// PB2TM is used for converting protobuf ABCI to ReapchainCore ABCI.
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
		tmVals[i] = NewValidator(pub, v.Power, v.Type)
	}
	return tmVals, nil
}

func (pb2tm) StandingMemberUpdates(vals []abci.ValidatorUpdate) ([]*StandingMember, error) {
	var smz []*StandingMember

	for _, v := range vals {
		if v.GetType() == "standing" {
			pubKey, err := cryptoenc.PubKeyFromProto(v.PubKey)
			if err != nil {
				return nil, err
			}
			smz = append(smz, NewStandingMember(pubKey, v.Power))
		}
	}
	return smz, nil
}

func (pb2tm) SteeringMemberCandidateUpdates(vals []abci.ValidatorUpdate) ([]*SteeringMemberCandidate, error) {
	var smz []*SteeringMemberCandidate

	for _, v := range vals {
		if v.GetType() == "steering" {
			pubKey, err := cryptoenc.PubKeyFromProto(v.PubKey)
			if err != nil {
				return nil, err
			}
			smz = append(smz, NewSteeringMemberCandidate(pubKey, v.Power))
		}
	}
	return smz, nil
}

func (pb2tm) QrnUpdates(qrnUpdates []abci.QrnUpdate) ([]*Qrn, error) {
	qrns := make([]*Qrn, len(qrnUpdates))
	for i, qrnUpdate := range qrnUpdates {
		qrns[i] = QrnFromAbci(&qrnUpdate)
	}
	return qrns, nil
}

func (pb2tm) VrfUpdates(vrfUpdates []abci.VrfUpdate) ([]*Vrf, error) {
	vrfs := make([]*Vrf, len(vrfUpdates))
	for i, vrfUpdate := range vrfUpdates {
		vrfs[i] = VrfFromAbci(&vrfUpdate)
	}
	return vrfs, nil
}

