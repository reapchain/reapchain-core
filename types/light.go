package types

import (
	"bytes"
	"errors"
	"fmt"

	tmproto "github.com/reapchain/reapchain-core/proto/podc/types"
)

// LightBlock is a SignedHeader and a ValidatorSet.
// It is the basis of the light client
type LightBlock struct {
	*SignedHeader `json:"signed_header"`
	ValidatorSet  *ValidatorSet `json:"validator_set"`
	StandingMemberSet          *StandingMemberSet          `json:"standing_member_set"`
	SteeringMemberCandidateSet *SteeringMemberCandidateSet `json:"steering_member_candidate_set"`
	QrnSet                     *QrnSet                     `json:"qrn_set"`
	VrfSet                     *VrfSet                     `json:"vrf_set"`

	NextQrnSet *QrnSet `json:"next_qrn_set"`
	NextVrfSet *VrfSet `json:"next_vrf_set"`

	SettingSteeringMember *SettingSteeringMember `json:"setting_steering_member"`

}

// ValidateBasic checks that the data is correct and consistent
//
// This does no verification of the signatures
func (lb LightBlock) ValidateBasic(chainID string) error {
	if lb.SignedHeader == nil {
		return errors.New("missing signed header")
	}
	if lb.ValidatorSet == nil {
		return errors.New("missing validator set")
	}

	if lb.StandingMemberSet == nil {
		return errors.New("missing standing member set")
	}

	if lb.QrnSet == nil {
		return errors.New("missing qrn set")
	}

	if err := lb.SignedHeader.ValidateBasic(chainID); err != nil {
		return fmt.Errorf("invalid signed header: %w", err)
	}
	if err := lb.ValidatorSet.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid validator set: %w", err)
	}

	if err := lb.StandingMemberSet.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid standing member set: %w", err)
	}

	if err := lb.QrnSet.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid qrn set: %w", err)
	}

	// make sure the validator set is consistent with the header
	if valSetHash := lb.ValidatorSet.Hash(); !bytes.Equal(lb.SignedHeader.ValidatorsHash, valSetHash) {
		return fmt.Errorf("expected validator hash of header to match validator set hash (%X != %X)",
			lb.SignedHeader.ValidatorsHash, valSetHash,
		)
	}

	if steeringMemberCandidateSetHash := lb.SteeringMemberCandidateSet.Hash(); !bytes.Equal(lb.SignedHeader.SteeringMemberCandidatesHash, steeringMemberCandidateSetHash) {
		return fmt.Errorf("expected sttering member candidate hash of header to match sttering member candidate set hash (%X != %X)",
			lb.SignedHeader.SteeringMemberCandidatesHash, steeringMemberCandidateSetHash,
		)
	}

	if standingMemberSetHash := lb.StandingMemberSet.Hash(); !bytes.Equal(lb.SignedHeader.StandingMembersHash, standingMemberSetHash) {
		return fmt.Errorf("expected standing member hash of header to match standing member set hash (%X != %X)",
			lb.SignedHeader.StandingMembersHash, standingMemberSetHash,
		)
	}

	if qrnSetHash := lb.QrnSet.Hash(); !bytes.Equal(lb.SignedHeader.QrnsHash, qrnSetHash) {
		return fmt.Errorf("expected qrn hash of header to match qrn set hash (%X != %X)",
			lb.SignedHeader.QrnsHash, qrnSetHash,
		)
	}

	if vrfSetHash := lb.VrfSet.Hash(); !bytes.Equal(lb.SignedHeader.VrfsHash, vrfSetHash) {
		return fmt.Errorf("expected vrf hash of header to match vrf set hash (%X != %X)",
			lb.SignedHeader.VrfsHash, vrfSetHash,
		)
	}
	

	return nil
}

// String returns a string representation of the LightBlock
func (lb LightBlock) String() string {
	return lb.StringIndented("")
}

// StringIndented returns an indented string representation of the LightBlock
//
// SignedHeader
// ValidatorSet
func (lb LightBlock) StringIndented(indent string) string {
	return fmt.Sprintf(`LightBlock{
%s  %v
%s  %v
%s}`,
		indent, lb.SignedHeader.StringIndented(indent+"  "),
		indent, lb.ValidatorSet.StringIndented(indent+"  "),
		indent)
}

// ToProto converts the LightBlock to protobuf
func (lb *LightBlock) ToProto() (*tmproto.LightBlock, error) {
	if lb == nil {
		return nil, nil
	}

	lbp := new(tmproto.LightBlock)
	var err error
	if lb.SignedHeader != nil {
		lbp.SignedHeader = lb.SignedHeader.ToProto()
	}
	if lb.ValidatorSet != nil {
		lbp.ValidatorSet, err = lb.ValidatorSet.ToProto()
		if err != nil {
			return nil, err
		}
	}

	if lb.StandingMemberSet != nil {
		lbp.StandingMemberSet, err = lb.StandingMemberSet.ToProto()
		if err != nil {
			return nil, err
		}
	}

	if lb.SteeringMemberCandidateSet != nil {
		lbp.SteeringMemberCandidateSet, err = lb.SteeringMemberCandidateSet.ToProto()
		if err != nil {
			return nil, err
		}
	}

	if lb.QrnSet != nil {
		lbp.QrnSet, err = lb.QrnSet.ToProto()
		if err != nil {
			return nil, err
		}
	}

	if lb.VrfSet != nil {
		lbp.VrfSet, err = lb.VrfSet.ToProto()
		if err != nil {
			return nil, err
		}
	}

	if lb.NextQrnSet != nil {
		lbp.NextQrnSet, err = lb.NextQrnSet.ToProto()
		if err != nil {
			return nil, err
		}
	}

	if lb.NextVrfSet != nil {
		lbp.NextVrfSet, err = lb.NextVrfSet.ToProto()
		if err != nil {
			return nil, err
		}
	}

	if lb.SettingSteeringMember != nil {
		lbp.SettingSteeringMember = lb.SettingSteeringMember.ToProto()
	}


	return lbp, nil
}

// LightBlockFromProto converts from protobuf back into the Lightblock.
// An error is returned if either the validator set or signed header are invalid
func LightBlockFromProto(pb *tmproto.LightBlock) (*LightBlock, error) {
	if pb == nil {
		return nil, errors.New("nil light block")
	}

	lb := new(LightBlock)

	if pb.SignedHeader != nil {
		sh, err := SignedHeaderFromProto(pb.SignedHeader)
		if err != nil {
			return nil, err
		}
		lb.SignedHeader = sh
	}

	if pb.ValidatorSet != nil {
		vals, err := ValidatorSetFromProto(pb.ValidatorSet)
		if err != nil {
			return nil, err
		}
		lb.ValidatorSet = vals
	}

	if pb.StandingMemberSet != nil {
		standingMemberSet, err := StandingMemberSetFromProto(pb.StandingMemberSet)
		if err != nil {
			return nil, err
		}
		lb.StandingMemberSet = standingMemberSet
	}

	if pb.SteeringMemberCandidateSet != nil {
		steeringMemberCandidateSet, err := SteeringMemberCandidateSetFromProto(pb.SteeringMemberCandidateSet)
		if err != nil {
			return nil, err
		}
		lb.SteeringMemberCandidateSet = steeringMemberCandidateSet
	}

	if pb.QrnSet != nil {
		qrnSet, err := QrnSetFromProto(pb.QrnSet)
		if err != nil {
			return nil, err
		}
		lb.QrnSet = qrnSet
	}

	if pb.NextQrnSet != nil {
		nextQrnSet, err := QrnSetFromProto(pb.NextQrnSet)
		if err != nil {
			return nil, err
		}
		lb.NextQrnSet = nextQrnSet
	}

	if pb.VrfSet != nil {
		vrfSet, err := VrfSetFromProto(pb.VrfSet)
		if err != nil {
			return nil, err
		}
		lb.VrfSet = vrfSet
	}

	if pb.NextVrfSet != nil {
		nextVrfSet, err := VrfSetFromProto(pb.NextVrfSet)
		if err != nil {
			return nil, err
		}
		lb.NextVrfSet = nextVrfSet
	}

	if pb.SettingSteeringMember != nil {
		lb.SettingSteeringMember = SettingSteeringMemberFromProto(pb.SettingSteeringMember)
	}

	return lb, nil
}

//-----------------------------------------------------------------------------

// SignedHeader is a header along with the commits that prove it.
type SignedHeader struct {
	*Header `json:"header"`

	Commit *Commit `json:"commit"`
}

// ValidateBasic does basic consistency checks and makes sure the header
// and commit are consistent.
//
// NOTE: This does not actually check the cryptographic signatures.  Make sure
// to use a Verifier to validate the signatures actually provide a
// significantly strong proof for this header's validity.
func (sh SignedHeader) ValidateBasic(chainID string) error {
	if sh.Header == nil {
		return errors.New("missing header")
	}
	if sh.Commit == nil {
		return errors.New("missing commit")
	}

	if err := sh.Header.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid header: %w", err)
	}
	if err := sh.Commit.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid commit: %w", err)
	}

	if sh.ChainID != chainID {
		return fmt.Errorf("header belongs to another chain %q, not %q", sh.ChainID, chainID)
	}

	// Make sure the header is consistent with the commit.
	if sh.Commit.Height != sh.Height {
		return fmt.Errorf("header and commit height mismatch: %d vs %d", sh.Height, sh.Commit.Height)
	}
	if hhash, chash := sh.Header.Hash(), sh.Commit.BlockID.Hash; !bytes.Equal(hhash, chash) {
		return fmt.Errorf("commit signs block %X, header is block %X", chash, hhash)
	}

	return nil
}

// String returns a string representation of SignedHeader.
func (sh SignedHeader) String() string {
	return sh.StringIndented("")
}

// StringIndented returns an indented string representation of SignedHeader.
//
// Header
// Commit
func (sh SignedHeader) StringIndented(indent string) string {
	return fmt.Sprintf(`SignedHeader{
%s  %v
%s  %v
%s}`,
		indent, sh.Header.StringIndented(indent+"  "),
		indent, sh.Commit.StringIndented(indent+"  "),
		indent)
}

// ToProto converts SignedHeader to protobuf
func (sh *SignedHeader) ToProto() *tmproto.SignedHeader {
	if sh == nil {
		return nil
	}

	psh := new(tmproto.SignedHeader)
	if sh.Header != nil {
		psh.Header = sh.Header.ToProto()
	}
	if sh.Commit != nil {
		psh.Commit = sh.Commit.ToProto()
	}

	return psh
}

// FromProto sets a protobuf SignedHeader to the given pointer.
// It returns an error if the header or the commit is invalid.
func SignedHeaderFromProto(shp *tmproto.SignedHeader) (*SignedHeader, error) {
	if shp == nil {
		return nil, errors.New("nil SignedHeader")
	}

	sh := new(SignedHeader)

	if shp.Header != nil {
		h, err := HeaderFromProto(shp.Header)
		if err != nil {
			return nil, err
		}
		sh.Header = &h
	}

	if shp.Commit != nil {
		c, err := CommitFromProto(shp.Commit)
		if err != nil {
			return nil, err
		}
		sh.Commit = c
	}

	return sh, nil
}
