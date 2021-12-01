package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/reapchain/reapchain-core/crypto"
	ce "github.com/reapchain/reapchain-core/crypto/encoding"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

const (
	nilSettingSteeringMemberStr string = "nil-SettingSteeringMember"
)

type SettingSteeringMember struct {
	Height                int64         `json:"height"`
	Timestamp             time.Time     `json:"timestamp"`
	CoordinatorPubKey     crypto.PubKey `json:"coordinator_pub_key"`
	SteeringMemberIndexes []int32       `json:"steering_member_indexes"`
	Signature             []byte        `json:"signature"`
}

func NewSettingSteeringMember(steeringMemberSize int) *SettingSteeringMember {
	settingSteeringMember := SettingSteeringMember{
		Timestamp:             time.Now(),
		SteeringMemberIndexes: make([]int32, steeringMemberSize),
	}

	return &settingSteeringMember
}

func (settingSteeringMember *SettingSteeringMember) ValidateBasic() error {
	if settingSteeringMember.Height < 0 {
		return errors.New("negative Height")
	}

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
		SteeringMemberIndexes: settingSteeringMember.SteeringMemberIndexes,
	}

	SettingSteeringMemberSignBytes, err := settingSteeringMemberProto.Marshal()
	if err != nil {
		panic(err)
	}
	return SettingSteeringMemberSignBytes
}

func (settingSteeringMember *SettingSteeringMember) GetSettingSteeringMemberBytes() []byte {
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
		SteeringMemberIndexes: settingSteeringMember.SteeringMemberIndexes,
		Signature:             settingSteeringMember.Signature,
	}

	settingSteeringMemberSignBytes, err := settingSteeringMemberProto.Marshal()
	if err != nil {
		panic(err)
	}
	return settingSteeringMemberSignBytes
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
		SteeringMemberIndexes: settingSteeringMember.SteeringMemberIndexes,
	}

	SettingSteeringMemberSignBytes, err := settingSteeringMemberProto.Marshal()
	if err != nil {
		panic(err)
	}
	return SettingSteeringMemberSignBytes
}

func (settingSteeringMember *SettingSteeringMember) VerifySign() bool {
	signBytes := settingSteeringMember.GetSettingSteeringMemberBytesForSign()
	if signBytes == nil {
		return false
	}

	return settingSteeringMember.CoordinatorPubKey.VerifySignature(signBytes, settingSteeringMember.Signature)
}

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
		SteeringMemberIndexes: settingSteeringMember.SteeringMemberIndexes,
		Signature:             settingSteeringMember.Signature,
	}

	return &settingSteeringMemberProto
}

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
	settingSteeringMember.SteeringMemberIndexes = settingSteeringMemberProto.SteeringMemberIndexes
	settingSteeringMember.Signature = settingSteeringMemberProto.Signature

	return settingSteeringMember
}
