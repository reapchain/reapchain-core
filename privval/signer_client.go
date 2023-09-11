package privval

import (
	"fmt"
	"time"

	"github.com/reapchain/reapchain-core/crypto"
	cryptoenc "github.com/reapchain/reapchain-core/crypto/encoding"
	privvalproto "github.com/reapchain/reapchain-core/proto/podc/privval"
	tmproto "github.com/reapchain/reapchain-core/proto/podc/types"
	"github.com/reapchain/reapchain-core/types"
)

// SignerClient implements PrivValidator.
// Handles remote validator connections that provide signing services
type SignerClient struct {
	endpoint *SignerListenerEndpoint
	chainID  string
}

var _ types.PrivValidator = (*SignerClient)(nil)

// NewSignerClient returns an instance of SignerClient.
// it will start the endpoint (if not already started)
func NewSignerClient(endpoint *SignerListenerEndpoint, chainID string) (*SignerClient, error) {
	if !endpoint.IsRunning() {
		if err := endpoint.Start(); err != nil {
			return nil, fmt.Errorf("failed to start listener endpoint: %w", err)
		}
	}

	return &SignerClient{endpoint: endpoint, chainID: chainID}, nil
}

// Close closes the underlying connection
func (sc *SignerClient) Close() error {
	return sc.endpoint.Close()
}

// IsConnected indicates with the signer is connected to a remote signing service
func (sc *SignerClient) IsConnected() bool {
	return sc.endpoint.IsConnected()
}

// WaitForConnection waits maxWait for a connection or returns a timeout error
func (sc *SignerClient) WaitForConnection(maxWait time.Duration) error {
	return sc.endpoint.WaitForConnection(maxWait)
}

//--------------------------------------------------------
// Implement PrivValidator

// Ping sends a ping request to the remote signer
func (sc *SignerClient) Ping() error {
	response, err := sc.endpoint.SendRequest(mustWrapMsg(&privvalproto.PingRequest{}))
	if err != nil {
		sc.endpoint.Logger.Error("SignerClient::Ping", "err", err)
		return nil
	}

	pb := response.GetPingResponse()
	if pb == nil {
		return err
	}

	return nil
}

// GetPubKey retrieves a public key from a remote signer
// returns an error if client is not able to provide the key
func (sc *SignerClient) GetPubKey() (crypto.PubKey, error) {
	response, err := sc.endpoint.SendRequest(mustWrapMsg(&privvalproto.PubKeyRequest{ChainId: sc.chainID}))
	if err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	resp := response.GetPubKeyResponse()
	if resp == nil {
		return nil, ErrUnexpectedResponse
	}
	if resp.Error != nil {
		return nil, &RemoteSignerError{Code: int(resp.Error.Code), Description: resp.Error.Description}
	}

	pk, err := cryptoenc.PubKeyFromProto(resp.PubKey)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

// GetPubKey retrieves a validator type from a remote signer
// returns an error if client is not able to provide the type
func (sc *SignerClient) GetType() (string, error) {
	response, err := sc.endpoint.SendRequest(mustWrapMsg(&privvalproto.TypeRequest{ChainId: sc.chainID}))
	if err != nil {
		return "", fmt.Errorf("send: %w", err)
	}

	resp := response.GetTypeResponse()
	if resp == nil {
		return "", ErrUnexpectedResponse
	}
	if resp.Error != nil {
		return "", &RemoteSignerError{Code: int(resp.Error.Code), Description: resp.Error.Description}
	}

	return resp.Type, nil
}


// SignVote requests a remote signer to sign a vote
func (sc *SignerClient) SignVote(chainID string, vote *tmproto.Vote) error {
	response, err := sc.endpoint.SendRequest(mustWrapMsg(&privvalproto.SignVoteRequest{Vote: vote, ChainId: chainID}))
	if err != nil {
		return err
	}

	resp := response.GetSignedVoteResponse()
	if resp == nil {
		return ErrUnexpectedResponse
	}
	if resp.Error != nil {
		return &RemoteSignerError{Code: int(resp.Error.Code), Description: resp.Error.Description}
	}

	*vote = resp.Vote

	return nil
}

func (sc *SignerClient) SignQrn(qrn *types.Qrn) error {
	response, err := sc.endpoint.SendRequest(mustWrapMsg(&privvalproto.SignQrnRequest{Qrn: qrn.ToProto()}))
	if err != nil {
		return err
	}

	resp := response.GetSignedQrnResponse()
	if resp == nil {
		return ErrUnexpectedResponse
	}
	if resp.Error != nil {
		return &RemoteSignerError{Code: int(resp.Error.Code), Description: resp.Error.Description}
	}

	qrn = types.QrnFromProto(&resp.Qrn)

	return nil
}

func (sc *SignerClient) SignSettingSteeringMember(settingSteeringMember *types.SettingSteeringMember) error {
	response, err := sc.endpoint.SendRequest(mustWrapMsg(&privvalproto.SignSettingSteeringMemberRequest{SettingSteeringMember: settingSteeringMember.ToProto()}))
	if err != nil {
		return err
	}

	resp := response.GetSignedSettingSteeringMemberResponse()
	if resp == nil {
		return ErrUnexpectedResponse
	}
	
	if resp.Error != nil {
		return &RemoteSignerError{Code: int(resp.Error.Code), Description: resp.Error.Description}
	}

	return nil
}

//TODO: mssong
func (sc *SignerClient) ProveVrf(vrf *types.Vrf) error {
	response, err := sc.endpoint.SendRequest(mustWrapMsg(&privvalproto.SignVrfRequest{Vrf: vrf.ToProto()}))
	if err != nil {
		return err
	}

	resp := response.GetSignedVrfResponse()
	if resp == nil {
		return ErrUnexpectedResponse
	}
	if resp.Error != nil {
		return &RemoteSignerError{Code: int(resp.Error.Code), Description: resp.Error.Description}
	}

	vrf = types.VrfFromProto(&resp.Vrf)

	return nil
}

// SignProposal requests a remote signer to sign a proposal
func (sc *SignerClient) SignProposal(chainID string, proposal *tmproto.Proposal) error {
	response, err := sc.endpoint.SendRequest(mustWrapMsg(
		&privvalproto.SignProposalRequest{Proposal: proposal, ChainId: chainID},
	))
	if err != nil {
		return err
	}

	resp := response.GetSignedProposalResponse()
	if resp == nil {
		return ErrUnexpectedResponse
	}
	if resp.Error != nil {
		return &RemoteSignerError{Code: int(resp.Error.Code), Description: resp.Error.Description}
	}

	*proposal = resp.Proposal

	return nil
}
