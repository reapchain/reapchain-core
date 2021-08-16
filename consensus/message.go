package consensus

import (
	"github.com/gogo/protobuf/proto"

	tmjson "github.com/reapchain/reapchain-core/libs/json"
	tmcons "github.com/reapchain/reapchain-core/proto/reapchain/consensus"
)

// Message is a message that can be sent and received on the Reactor
type Message interface {
	ValidateBasic() error
}

func decodeMsg(bz []byte) (msg Message, err error) {
	pb := &tmcons.Message{}
	if err = proto.Unmarshal(bz, pb); err != nil {
		return msg, err
	}

	return MsgFromProto(pb)
}

func init() {
	tmjson.RegisterType(&NewRoundStepMessage{}, "reapchain/NewRoundStepMessage")
	tmjson.RegisterType(&NewValidBlockMessage{}, "reapchain/NewValidBlockMessage")
	tmjson.RegisterType(&ProposalMessage{}, "reapchain/Proposal")
	tmjson.RegisterType(&ProposalPOLMessage{}, "reapchain/ProposalPOL")
	tmjson.RegisterType(&BlockPartMessage{}, "reapchain/BlockPart")
	tmjson.RegisterType(&VoteMessage{}, "reapchain/Vote")
	tmjson.RegisterType(&HasVoteMessage{}, "reapchain/HasVote")
	tmjson.RegisterType(&VoteSetMaj23Message{}, "reapchain/VoteSetMaj23")
	tmjson.RegisterType(&VoteSetBitsMessage{}, "reapchain/VoteSetBits")
	tmjson.RegisterType(&QrnMessage{}, "reapchain/Qrn")
	tmjson.RegisterType(&HasQrnMessage{}, "reapchain/HasQrn")
	tmjson.RegisterType(&VrfMessage{}, "reapchain/Vrf")
	tmjson.RegisterType(&HasVrfMessage{}, "reapchain/HasVrf")
}
