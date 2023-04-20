package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/reapchain/reapchain-core/crypto"
	tmbytes "github.com/reapchain/reapchain-core/libs/bytes"
	tmjson "github.com/reapchain/reapchain-core/libs/json"
	tmos "github.com/reapchain/reapchain-core/libs/os"
	tmproto "github.com/reapchain/reapchain-core/proto/podc/types"
	tmtime "github.com/reapchain/reapchain-core/types/time"
)

const (
	// MaxChainIDLen is a maximum length of the chain ID.
	MaxChainIDLen = 50
)

//------------------------------------------------------------
// core types for a genesis definition
// NOTE: any changes to the genesis definition should
// be reflected in the documentation:
// docs/reapchain-core-core/using-reapchain-core.md

// GenesisValidator is an initial validator.
type GenesisValidator struct {
	Address Address       `json:"address"`
	PubKey  crypto.PubKey `json:"pub_key"`
	Type  	string 				`json:"type"`
	Power   int64         `json:"power"`
	Name    string        `json:"name"`
}

// GenesisStandingMember is an initial standingMember.
type GenesisMember struct {
	Address Address       `json:"address"`
	PubKey  crypto.PubKey `json:"pub_key"`
	Name    string        `json:"name"`
	Power   int64         `json:"power"`
}

// GenesisDoc defines the initial conditions for a reapchain-core blockchain, in particular its validator set.
type GenesisDoc struct {
	GenesisTime     time.Time                `json:"genesis_time"`
	ChainID         string                   `json:"chain_id"`
	InitialHeight   int64                    `json:"initial_height"`
	ConsensusParams *tmproto.ConsensusParams `json:"consensus_params,omitempty"`
	Validators      []GenesisValidator       `json:"validators,omitempty"`
	AppHash         tmbytes.HexBytes         `json:"app_hash"`
	AppState        json.RawMessage          `json:"app_state,omitempty"`
	StandingMembers          []GenesisMember          `json:"standing_members,omitempty"`
	ConsensusRound           ConsensusRound           `json:"consensus_round"`
	Qrns                     []Qrn                    `json:"qrns,omitempty"`
	SteeringMemberCandidates []GenesisMember          `json:"steering_member_candidates"`
	Vrfs                     []Vrf                    `json:"vrfs"`
}

// SaveAs is a utility method for saving GenensisDoc as a JSON file.
func (genDoc *GenesisDoc) SaveAs(file string) error {
	genDocBytes, err := tmjson.MarshalIndent(genDoc, "", "  ")
	if err != nil {
		return err
	}
	return tmos.WriteFile(file, genDocBytes, 0644)
}

// ValidatorHash returns the hash of the validator set contained in the GenesisDoc
func (genDoc *GenesisDoc) ValidatorHash() []byte {
	vals := make([]*Validator, len(genDoc.Validators))
	for i, v := range genDoc.Validators {
		vals[i] = NewValidator(v.PubKey, v.Power, v.Type)
	}
	vset := NewValidatorSet(vals)
	return vset.Hash()
}

// ValidateAndComplete checks that all necessary fields are present
// and fills in defaults for optional fields left empty
func (genDoc *GenesisDoc) ValidateAndComplete() error {
	if genDoc.ChainID == "" {
		return errors.New("genesis doc must include non-empty chain_id")
	}
	if len(genDoc.ChainID) > MaxChainIDLen {
		return fmt.Errorf("chain_id in genesis doc is too long (max: %d)", MaxChainIDLen)
	}
	if genDoc.InitialHeight < 0 {
		return fmt.Errorf("initial_height cannot be negative (got %v)", genDoc.InitialHeight)
	}
	if genDoc.InitialHeight == 0 {
		genDoc.InitialHeight = 1
	}

	if genDoc.ConsensusParams == nil {
		genDoc.ConsensusParams = DefaultConsensusParams()
	} else if err := ValidateConsensusParams(*genDoc.ConsensusParams); err != nil {
		return err
	}

	for i, v := range genDoc.Validators {
		if v.Power == 0 {
			return fmt.Errorf("the genesis file cannot contain validators with no voting power: %v", v)
		}
		if len(v.Address) > 0 && !bytes.Equal(v.PubKey.Address(), v.Address) {
			return fmt.Errorf("incorrect address for validator %v in the genesis file, should be %v", v, v.PubKey.Address())
		}
		if len(v.Address) == 0 {
			genDoc.Validators[i].Address = v.PubKey.Address()
		}
	}

	for i, standingMember := range genDoc.StandingMembers {

		if standingMember.Power == 0 {
			return fmt.Errorf("the genesis file cannot contain standingMembers with no voting power: %v", standingMember)
		}
		if len(standingMember.Address) > 0 && !bytes.Equal(standingMember.PubKey.Address(), standingMember.Address) {
			return fmt.Errorf("incorrect address for validator %v in the genesis file, should be %v", standingMember, standingMember.PubKey.Address())
		}
		if len(standingMember.Address) == 0 {
			genDoc.StandingMembers[i].Address = standingMember.PubKey.Address()
		}
	}

	for _, qrn := range genDoc.Qrns {
		if qrn.Timestamp.Sub(genDoc.GenesisTime) > 0 {
			return fmt.Errorf("Invalid qrn timestamp: qrnTimestamp = %v / genesisTimestamp = %v", qrn.Timestamp, genDoc.GenesisTime)
		}

		if err := qrn.ValidateBasic(); err != nil {
			return fmt.Errorf("Qrn error: %v", err)
		}

		if qrn.VerifySign(genDoc.ChainID) == false {
			return fmt.Errorf("Incorrect sign of qrn")
		}
	}


	if genDoc.GenesisTime.IsZero() {
		genDoc.GenesisTime = tmtime.Now()
	}

	return nil
}

//------------------------------------------------------------
// Make genesis state from file

// GenesisDocFromJSON unmarshalls JSON data into a GenesisDoc.
func GenesisDocFromJSON(jsonBlob []byte) (*GenesisDoc, error) {
	genDoc := GenesisDoc{}
	err := tmjson.Unmarshal(jsonBlob, &genDoc)
	if err != nil {
		return nil, err
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return nil, err
	}

	return &genDoc, err
}

// GenesisDocFromFile reads JSON data from a file and unmarshalls it into a GenesisDoc.
func GenesisDocFromFile(genDocFile string) (*GenesisDoc, error) {
	jsonBlob, err := ioutil.ReadFile(genDocFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't read GenesisDoc file: %w", err)
	}
	genDoc, err := GenesisDocFromJSON(jsonBlob)
	if err != nil {
		return nil, fmt.Errorf("error reading GenesisDoc at %s: %w", genDocFile, err)
	}
	return genDoc, nil
}
