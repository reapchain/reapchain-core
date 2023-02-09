package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/reapchain/reapchain-core/crypto"
	ce "github.com/reapchain/reapchain-core/crypto/encoding"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain-core/types"
)

const (
	nilSettingSteeringMembertr string = "nil-SettingSteeringMember"
)

// It manages steering member information which are generated via coordinator
type SettingSteeringMember struct {
	Height                int64         `json:"height"`
	Timestamp             time.Time     `json:"timestamp"`
	CoordinatorPubKey     crypto.PubKey `json:"coordinator_pub_key"`
	SteeringMemberAddresses [][]byte       `json:"steering_member_addresses"`
	Signature             []byte        `json:"signature"`
}

func NewSettingSteeringMember(steeringMemberSize int) *SettingSteeringMember {
	settingSteeringMember := SettingSteeringMember{
		Timestamp:             time.Now(),
		SteeringMemberAddresses: make([][]byte, 0, steeringMemberSize),
	}

	return &settingSteeringMember
}

func (settingSteeringMember *SettingSteeringMember) ValidateBasic() error {
	// check block height
	if settingSteeringMember.Height < 0 {
		return errors.New("negative Height")
	}

	// check coordinator address
	coordinatorAddress := settingSteeringMember.CoordinatorPubKey.Address()
	if len(coordinatorAddress) != crypto.AddressSize {
		return fmt.Errorf("expected StandingMemberAddress size to be %d bytes, got %d bytes",
			crypto.AddressSize,
			len(coordinatorAddress),
		)
	}

	return nil
}

func (settingSteeringMember *SettingSteeringMember) Copy() *SettingSteeringMember {
	if settingSteeringMember == nil {
		return nil
	}
	settingSteeringMemberCopy := *settingSteeringMember
	return &settingSteeringMemberCopy
}

// It is called via coordinator for signing the steering members who are selected as validator next round.
// It returns the bytes of the steering members
func (settingSteeringMember *SettingSteeringMember) GetSettingSteeringMemberBytesForSign() []byte {
	if settingSteeringMember == nil {
		return nil
	}

	pubKeyProto, err := ce.PubKeyToProto(settingSteeringMember.CoordinatorPubKey)
	if err != nil {
		panic(err)
	}


	settingSteeringMemberProto := tmproto.SettingSteeringMember{
		Height:                settingSteeringMember.Height,
		Timestamp:             settingSteeringMember.Timestamp,
		CoordinatorPubKey:     pubKeyProto,
		SteeringMemberAddresses: settingSteeringMember.SteeringMemberAddresses,
	}

	SettingSteeringMemberignBytes, err := settingSteeringMemberProto.Marshal()
	if err != nil {
		panic(err)
	}
	return SettingSteeringMemberignBytes
}

func (settingSteeringMember *SettingSteeringMember) GetBytesForSign() []byte {
	if settingSteeringMember == nil {
		return nil
	}

	pubKeyProto, err := ce.PubKeyToProto(settingSteeringMember.CoordinatorPubKey)
	if err != nil {
		panic(err)
	}

	settingSteeringMemberProto := tmproto.SettingSteeringMember{
		Height:                settingSteeringMember.Height,
		Timestamp:             settingSteeringMember.Timestamp,
		CoordinatorPubKey:     pubKeyProto,
		SteeringMemberAddresses: settingSteeringMember.SteeringMemberAddresses,
	}

	SettingSteeringMemberignBytes, err := settingSteeringMemberProto.Marshal()
	if err != nil {
		panic(err)
	}
	return SettingSteeringMemberignBytes
}

// Validate the coordinator sign.
// Each node gets own cornidator information and checks the settingSteeringMember is valid with the signature which is included the type
func (settingSteeringMember *SettingSteeringMember) VerifySign() bool {
	signBytes := settingSteeringMember.GetSettingSteeringMemberBytesForSign()
	if signBytes == nil {
		return false
	}

	return settingSteeringMember.CoordinatorPubKey.VerifySignature(signBytes, settingSteeringMember.Signature)
}

// Convert the type to proto puffer type to send the type to other peer or SDK
func (settingSteeringMember *SettingSteeringMember) ToProto() *tmproto.SettingSteeringMember {
	if settingSteeringMember == nil {
		return nil
	}

	pubKey, err := ce.PubKeyToProto(settingSteeringMember.CoordinatorPubKey)
	if err != nil {
		return nil
	}

	settingSteeringMemberProto := tmproto.SettingSteeringMember{
		Height:                settingSteeringMember.Height,
		Timestamp:             settingSteeringMember.Timestamp,
		CoordinatorPubKey:     pubKey,
		SteeringMemberAddresses: settingSteeringMember.SteeringMemberAddresses,
		Signature:             settingSteeringMember.Signature,
	}

	return &settingSteeringMemberProto
}

// Convert the setting steering member's proto puffer type to this type to apply the reapchain-core
func SettingSteeringMemberFromProto(settingSteeringMemberProto *tmproto.SettingSteeringMember) *SettingSteeringMember {
	if settingSteeringMemberProto == nil {
		return nil
	}

	pubKey, err := ce.PubKeyFromProto(settingSteeringMemberProto.CoordinatorPubKey)
	if err != nil {
		return nil
	}

	settingSteeringMember := new(SettingSteeringMember)
	settingSteeringMember.Height = settingSteeringMemberProto.Height
	settingSteeringMember.Timestamp = settingSteeringMemberProto.Timestamp
	settingSteeringMember.CoordinatorPubKey = pubKey
	settingSteeringMember.SteeringMemberAddresses = settingSteeringMemberProto.SteeringMemberAddresses
	settingSteeringMember.Signature = settingSteeringMemberProto.Signature

	return settingSteeringMember
}
