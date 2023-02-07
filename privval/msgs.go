package privval

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	privvalproto "github.com/reapchain/reapchain-core/proto/reapchain-core/privval"
)

// TODO: Add ChainIDRequest

func mustWrapMsg(pb proto.Message) privvalproto.Message {
	msg := privvalproto.Message{}

	switch pb := pb.(type) {
	case *privvalproto.Message:
		msg = *pb
	case *privvalproto.PubKeyRequest:
		msg.Sum = &privvalproto.Message_PubKeyRequest{PubKeyRequest: pb}
	case *privvalproto.TypeRequest:
		msg.Sum = &privvalproto.Message_TypeRequest{TypeRequest: pb}
	case *privvalproto.PubKeyResponse:
		msg.Sum = &privvalproto.Message_PubKeyResponse{PubKeyResponse: pb}
	case *privvalproto.SignVoteRequest:
		msg.Sum = &privvalproto.Message_SignVoteRequest{SignVoteRequest: pb}
	case *privvalproto.SignedVoteResponse:
		msg.Sum = &privvalproto.Message_SignedVoteResponse{SignedVoteResponse: pb}
	case *privvalproto.SignedProposalResponse:
		msg.Sum = &privvalproto.Message_SignedProposalResponse{SignedProposalResponse: pb}
	case *privvalproto.SignProposalRequest:
		msg.Sum = &privvalproto.Message_SignProposalRequest{SignProposalRequest: pb}
	case *privvalproto.PingRequest:
		msg.Sum = &privvalproto.Message_PingRequest{PingRequest: pb}
	case *privvalproto.PingResponse:
		msg.Sum = &privvalproto.Message_PingResponse{PingResponse: pb}

	case *privvalproto.SignQrnRequest:
		msg.Sum = &privvalproto.Message_SignQrnRequest{SignQrnRequest: pb}
	case *privvalproto.SignedQrnResponse:
		msg.Sum = &privvalproto.Message_SignedQrnResponse{SignedQrnResponse: pb}

	case *privvalproto.SignSettingSteeringMemberRequest:
		msg.Sum = &privvalproto.Message_SignSettingSteeringMemberRequest{SignSettingSteeringMemberRequest: pb}
	case *privvalproto.SignedSettingSteeringMemberResponse:
		msg.Sum = &privvalproto.Message_SignedSettingSteeringMemberResponse{SignedSettingSteeringMemberResponse: pb}

	case *privvalproto.SignVrfRequest:
		msg.Sum = &privvalproto.Message_SignVrfRequest{SignVrfRequest: pb}
	case *privvalproto.SignedVrfResponse:
		msg.Sum = &privvalproto.Message_SignedVrfResponse{SignedVrfResponse: pb}

	default:
		panic(fmt.Errorf("unknown message type %T", pb))
	}

	return msg
}
