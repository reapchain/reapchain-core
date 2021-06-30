package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"gitlab.reappay.net/sucs-lab/reapchain/crypto"
	tmbytes "gitlab.reappay.net/sucs-lab/reapchain/libs/bytes"
	tmjson "gitlab.reappay.net/sucs-lab/reapchain/libs/json"
	tmos "gitlab.reappay.net/sucs-lab/reapchain/libs/os"
	tmproto "gitlab.reappay.net/sucs-lab/reapchain/proto/reapchain/types"
	tmtime "gitlab.reappay.net/sucs-lab/reapchain/types/time"
)

const (
	// MaxChainIDLen is a maximum length of the chain ID.
	MaxChainIDLen = 50
)

//------------------------------------------------------------
// core types for a genesis definition
// NOTE: any changes to the genesis definition should
// be reflected in the documentation:
// docs/reapchain-core/using-reapchain.md

// GenesisValidator is an initial validator.
type GenesisValidator struct {
	Address Address       `json:"address"`
	PubKey  crypto.PubKey `json:"pub_key"`
	Power   int64         `json:"power"`
	Name    string        `json:"name"`
}

type GenesisStandingMember struct {
	Address Address       `json:"address"`
	PubKey  crypto.PubKey `json:"pub_key"`
	Name    string        `json:"name"`
}

// GenesisDoc defines the initial conditions for a reapchain blockchain, in particular its validator set.
type GenesisDoc struct {
	GenesisTime     time.Time                `json:"genesis_time"`
	ChainID         string                   `json:"chain_id"`
	InitialHeight   int64                    `json:"initial_height"`
	ConsensusParams *tmproto.ConsensusParams `json:"consensus_params,omitempty"`
	Validators      []GenesisValidator       `json:"validators,omitempty"`
	AppHash         tmbytes.HexBytes         `json:"app_hash"`
	AppState        json.RawMessage          `json:"app_state,omitempty"`

	StandingMembers    []GenesisStandingMember `json:"standing_members,omitempty"`
	Qns                []Qn                    `json:"qns,omitempty"`
	ConsensusRoundInfo ConsensusRound          `json:"consensus_round_info"`
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
		vals[i] = NewValidator(v.PubKey, v.Power)
	}
	vset := NewValidatorSet(vals)
	return vset.Hash()
}

func (genDoc *GenesisDoc) StandingMemberHash() []byte {
	sms := make([]*StandingMember, len(genDoc.StandingMembers))
	for i, v := range genDoc.StandingMembers {
		sms[i] = NewStandingMember(v.PubKey)
	}
	smSet := NewStandingMemberSet(sms)
	return smSet.Hash()
}

func (genDoc *GenesisDoc) QnHash() []byte {
	qns := make([]*Qn, len(genDoc.Qns))
	for i, qn := range genDoc.Qns {
		qns[i] = NewQn(qn.PubKey, qn.Value, qn.Height)
	}
	qnSet := NewQnSet(qns)
	return qnSet.Hash()
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

	// 상임위 검증
	for i, s := range genDoc.StandingMembers {
		if len(s.Address) > 0 && !bytes.Equal(s.PubKey.Address(), s.Address) {
			return fmt.Errorf("incorrect address for stending member %v in the genesis file, should be %v", s, s.PubKey.Address())
		}
		if len(s.Address) == 0 {
			genDoc.StandingMembers[i].Address = s.PubKey.Address()
		}
	}

	for i, qn := range genDoc.Qns {
		if len(qn.Address) > 0 && !bytes.Equal(qn.PubKey.Address(), qn.Address) {
			return fmt.Errorf("incorrect address for stending member %v in the genesis file, should be %v", qn, qn.PubKey.Address())
		}
		if len(qn.Address) == 0 {
			genDoc.Qns[i].Address = qn.PubKey.Address()
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
