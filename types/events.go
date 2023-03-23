package types

import (
	"fmt"

	abci "github.com/reapchain/reapchain-core/abci/types"
	tmjson "github.com/reapchain/reapchain-core/libs/json"
	tmpubsub "github.com/reapchain/reapchain-core/libs/pubsub"
	tmquery "github.com/reapchain/reapchain-core/libs/pubsub/query"
)

// Reserved event types (alphabetically sorted).
const (
	// Block level events for mass consumption by users.
	// These events are triggered from the state package,
	// after a block has been committed.
	// These are also used by the tx indexer for async indexing.
	// All of this data can be fetched through the rpc.
	EventNewBlock            = "NewBlock"
	EventNewBlockHeader      = "NewBlockHeader"
	EventNewEvidence         = "NewEvidence"
	EventTx                  = "Tx"
	EventValidatorSetUpdates = "ValidatorSetUpdates"
	EventSteeringMemberCandidateSetUpdates = "SteeringMemberCandidateSetUpdates"
	EventStandingMemberSetUpdates          = "StandingMemberSetUpdates"


	// Internal consensus events.
	// These are used for testing the consensus state machine.
	// They can also be used to build real-time consensus visualizers.
	EventCompleteProposal = "CompleteProposal"
	EventLock             = "Lock"
	EventNewRound         = "NewRound"
	EventNewRoundStep     = "NewRoundStep"
	EventPolka            = "Polka"
	EventRelock           = "Relock"
	EventTimeoutPropose   = "TimeoutPropose"
	EventTimeoutWait      = "TimeoutWait"
	EventUnlock           = "Unlock"
	EventValidBlock       = "ValidBlock"
	EventVote             = "Vote"

	EventQrn             					= "Qrn"
	EventVrf             					= "Vrf"
	EventSettingSteeringMember    = "SettingSteeringMember"
)

// ENCODING / DECODING

// TMEventData implements events.EventData.
type TMEventData interface {
	// empty interface
}

func init() {
	tmjson.RegisterType(EventDataNewBlock{}, "reapchain-core/event/NewBlock")
	tmjson.RegisterType(EventDataNewBlockHeader{}, "reapchain-core/event/NewBlockHeader")
	tmjson.RegisterType(EventDataNewEvidence{}, "reapchain-core/event/NewEvidence")
	tmjson.RegisterType(EventDataTx{}, "reapchain-core/event/Tx")
	tmjson.RegisterType(EventDataRoundState{}, "reapchain-core/event/RoundState")
	tmjson.RegisterType(EventDataNewRound{}, "reapchain-core/event/NewRound")
	tmjson.RegisterType(EventDataCompleteProposal{}, "reapchain-core/event/CompleteProposal")
	tmjson.RegisterType(EventDataVote{}, "reapchain-core/event/Vote")
	tmjson.RegisterType(EventDataSteeringMemberCandidateSetUpdates{}, "reapchain/event/SteeringMemberCandidateSetUpdates")
	tmjson.RegisterType(EventDataStandingMemberSetUpdates{}, "reapchain/event/StandingMemberSetUpdates")
	tmjson.RegisterType(EventDataString(""), "reapchain-core/event/ProposalString")

	tmjson.RegisterType(EventDataQrn{}, "reapchain/event/Qrn")
	tmjson.RegisterType(EventDataVrf{}, "reapchain/event/Vrf")
	tmjson.RegisterType(EventDataSettingSteeringMember{}, "reapchain/event/SetEventDataSettingSteeringMember")

}

// Most event messages are basic types (a block, a transaction)
// but some (an input to a call tx or a receive) are more exotic

type EventDataNewBlock struct {
	Block *Block `json:"block"`

	ResultBeginBlock abci.ResponseBeginBlock `json:"result_begin_block"`
	ResultEndBlock   abci.ResponseEndBlock   `json:"result_end_block"`
	ConsensusInfo 	 abci.ConsensusInfo 	 `json:"consensus_info"`
}

type EventDataNewBlockHeader struct {
	Header Header `json:"header"`

	NumTxs           int64                   `json:"num_txs"` // Number of txs in a block
	ResultBeginBlock abci.ResponseBeginBlock `json:"result_begin_block"`
	ResultEndBlock   abci.ResponseEndBlock   `json:"result_end_block"`
}

type EventDataNewEvidence struct {
	Evidence Evidence `json:"evidence"`

	Height int64 `json:"height"`
}

// All txs fire EventDataTx
type EventDataTx struct {
	abci.TxResult
}

// NOTE: This goes into the replay WAL
type EventDataRoundState struct {
	Height int64  `json:"height"`
	Round  int32  `json:"round"`
	Step   string `json:"step"`
}

type ValidatorInfo struct {
	Address Address `json:"address"`
	Index   int32   `json:"index"`
}

type EventDataNewRound struct {
	Height int64  `json:"height"`
	Round  int32  `json:"round"`
	Step   string `json:"step"`

	Proposer ValidatorInfo `json:"proposer"`
}

type EventDataCompleteProposal struct {
	Height int64  `json:"height"`
	Round  int32  `json:"round"`
	Step   string `json:"step"`

	BlockID BlockID `json:"block_id"`
}

type EventDataVote struct {
	Vote *Vote
}

// Add event
type EventDataQrn struct {
	Qrn *Qrn
}

type EventDataVrf struct {
	Vrf *Vrf
}

type EventDataSettingSteeringMember struct {
	SettingSteeringMember *SettingSteeringMember
}


type EventDataString string

type EventDataSteeringMemberCandidateSetUpdates struct {
	SteeringMemberCandidateUpdates []*SteeringMemberCandidate `json:"steering_member_candidate_updates"`
}

type EventDataStandingMemberSetUpdates struct {
	StandingMemberUpdates []*StandingMember `json:"standing_member_updates"`
}

// PUBSUB

const (
	// EventTypeKey is a reserved composite key for event name.
	EventTypeKey = "tm.event"
	// TxHashKey is a reserved key, used to specify transaction's hash.
	// see EventBus#PublishEventTx
	TxHashKey = "tx.hash"
	// TxHeightKey is a reserved key, used to specify transaction block's height.
	// see EventBus#PublishEventTx
	TxHeightKey = "tx.height"

	// BlockHeightKey is a reserved key used for indexing BeginBlock and Endblock
	// events.
	BlockHeightKey = "block.height"
)

var (
	EventQueryCompleteProposal    = QueryForEvent(EventCompleteProposal)
	EventQueryLock                = QueryForEvent(EventLock)
	EventQueryNewBlock            = QueryForEvent(EventNewBlock)
	EventQueryNewBlockHeader      = QueryForEvent(EventNewBlockHeader)
	EventQueryNewEvidence         = QueryForEvent(EventNewEvidence)
	EventQueryNewRound            = QueryForEvent(EventNewRound)
	EventQueryNewRoundStep        = QueryForEvent(EventNewRoundStep)
	EventQueryPolka               = QueryForEvent(EventPolka)
	EventQueryRelock              = QueryForEvent(EventRelock)
	EventQueryTimeoutPropose      = QueryForEvent(EventTimeoutPropose)
	EventQueryTimeoutWait         = QueryForEvent(EventTimeoutWait)
	EventQueryTx                  = QueryForEvent(EventTx)
	EventQueryUnlock              = QueryForEvent(EventUnlock)
	EventQuerySteeringMemberCandidateSetUpdates = QueryForEvent(EventSteeringMemberCandidateSetUpdates)
	EventQueryStandingMemberSetUpdates          = QueryForEvent(EventStandingMemberSetUpdates)
	EventQueryValidBlock          = QueryForEvent(EventValidBlock)
	EventQueryVote                = QueryForEvent(EventVote)

	EventQueryQrn                              = QueryForEvent(EventQrn)
	EventQueryVrf                              = QueryForEvent(EventVrf)
	EventQuerySettingSteeringMember						 = QueryForEvent(EventSettingSteeringMember)

)

func EventQueryTxFor(tx Tx) tmpubsub.Query {
	return tmquery.MustParse(fmt.Sprintf("%s='%s' AND %s='%X'", EventTypeKey, EventTx, TxHashKey, tx.Hash()))
}

func QueryForEvent(eventType string) tmpubsub.Query {
	return tmquery.MustParse(fmt.Sprintf("%s='%s'", EventTypeKey, eventType))
}

// BlockEventPublisher publishes all block related events
type BlockEventPublisher interface {
	PublishEventNewBlock(block EventDataNewBlock) error
	PublishEventNewBlockHeader(header EventDataNewBlockHeader) error
	PublishEventNewEvidence(evidence EventDataNewEvidence) error
	PublishEventTx(EventDataTx) error
	PublishEventSteeringMemberCandidateSetUpdates(EventDataSteeringMemberCandidateSetUpdates) error
	PublishEventStandingMemberSetUpdates(EventDataStandingMemberSetUpdates) error
}

type TxEventPublisher interface {
	PublishEventTx(EventDataTx) error
}
