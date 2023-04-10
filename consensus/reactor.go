package consensus

import (
	"fmt"
	"reflect"
	"time"

	bc "github.com/reapchain/reapchain-core/blockchain"
	cstypes "github.com/reapchain/reapchain-core/consensus/types"
	"github.com/reapchain/reapchain-core/libs/bits"
	tmevents "github.com/reapchain/reapchain-core/libs/events"
	"github.com/reapchain/reapchain-core/libs/log"
	tmsync "github.com/reapchain/reapchain-core/libs/sync"
	"github.com/reapchain/reapchain-core/p2p"
	bcproto "github.com/reapchain/reapchain-core/proto/reapchain-core/blockchain"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain-core/types"
	sm "github.com/reapchain/reapchain-core/state"
	"github.com/reapchain/reapchain-core/types"
)

const (
	StateChannel                 = byte(0x20)
	DataChannel                  = byte(0x21)
	VoteChannel                  = byte(0x22)
	VoteSetBitsChannel           = byte(0x23)
	QrnChannel                   = byte(0x24)
	VrfChannel                   = byte(0x25)
	SettingSteeringMemberChannel = byte(0x26)
	CatchUpChannel = byte(0x27)

	maxMsgSize = 1048576 // 1MB; NOTE/TODO: keep in sync with types.PartSet sizes.

	blocksToContributeToBecomeGoodPeer = 10000
	votesToContributeToBecomeGoodPeer  = 10000
)

//-----------------------------------------------------------------------------

// Reactor defines a reactor for the consensus service.
type Reactor struct {
	p2p.BaseReactor // BaseService + p2p.Switch

	chainID string
	conS *State

	mtx      tmsync.RWMutex
	waitSync bool
	eventBus *types.EventBus
	rs       *cstypes.RoundState

	Metrics *Metrics

	stateStore sm.Store
	CatchupQrnMessages []*QrnMessage
	CatchupVrfMessages []*VrfMessage
	CatchupSettingSteeringMemberMessage *SettingSteeringMemberMessage
}

type ReactorOption func(*Reactor)

// NewReactor returns a new Reactor with the given
// consensusState.
func NewReactor(chainID string, consensusState *State, stateStore sm.Store, waitSync bool, options ...ReactorOption) *Reactor {
	conR := &Reactor{
		conS:     consensusState,
		chainID: chainID,
		stateStore: stateStore,
		waitSync: waitSync,
		rs:       consensusState.GetRoundState(),
		Metrics:  NopMetrics(),
		CatchupQrnMessages: make([]*QrnMessage, 0),
		CatchupVrfMessages: make([]*VrfMessage, 0),
		CatchupSettingSteeringMemberMessage: nil,
	}
	conR.BaseReactor = *p2p.NewBaseReactor("Consensus", conR)

	for _, option := range options {
		option(conR)
	}

	return conR
}

// OnStart implements BaseService by subscribing to events, which later will be
// broadcasted to other peers and starting state if we're not in fast sync.
func (conR *Reactor) OnStart() error {
	conR.Logger.Info("Reactor ", "waitSync", conR.WaitSync())

	// start routine that computes peer statistics for evaluating peer quality
	go conR.peerStatsRoutine()

	conR.subscribeToBroadcastEvents()
	go conR.updateRoundStateRoutine()

	if !conR.WaitSync() {
		err := conR.conS.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

// OnStop implements BaseService by unsubscribing from events and stopping
// state.
func (conR *Reactor) OnStop() {
	conR.unsubscribeFromBroadcastEvents()
	if err := conR.conS.Stop(); err != nil {
		conR.Logger.Error("Error stopping consensus state", "err", err)
	}
	if !conR.WaitSync() {
		conR.conS.Wait()
	}
}

// SwitchToConsensus switches from fast_sync mode to consensus mode.
// It resets the state, turns off fast_sync, and starts the consensus state-machine
func (conR *Reactor) SwitchToConsensus(state sm.State, skipWAL bool) {
	conR.Logger.Info("SwitchToConsensus")

	// We have no votes, so reconstruct LastCommit from SeenCommit.
	if state.LastBlockHeight > 0 {
		conR.conS.reconstructLastCommit(state)
	}
	
	// if we have CatchupQrnMessages, add the qrns
	for _, currentQrnMessage := range conR.CatchupQrnMessages {
		state.NextQrnSet.AddQrn(state.ChainID, currentQrnMessage.Qrn)
	}

	// if we have CatchupVrfMessages, add the vrfs
	for _, currentVrfMessage := range conR.CatchupVrfMessages {
		state.NextVrfSet.AddVrf(currentVrfMessage.Vrf)
	}
	
	// if we have CatchupSettingSteeringMemberMessage, and it is releated current consensus round, we replace the steering member
	if (conR.CatchupSettingSteeringMemberMessage != nil) {
		if (state.ConsensusRound.ConsensusStartBlockHeight + int64(state.ConsensusRound.Period) == conR.CatchupSettingSteeringMemberMessage.SettingSteeringMember.Height)  {
			state.SettingSteeringMember = conR.CatchupSettingSteeringMemberMessage.SettingSteeringMember
		}
	}

	// NOTE: The line below causes broadcastNewRoundStepRoutine() to broadcast a
	// NewRoundStepMessage.
	conR.conS.updateToState(state)

	conR.mtx.Lock()
	conR.waitSync = false
	conR.mtx.Unlock()
	conR.Metrics.FastSyncing.Set(0)
	conR.Metrics.StateSyncing.Set(0)

	if skipWAL {
		conR.conS.doWALCatchup = false
	}
	err := conR.conS.Start()
	if err != nil {
		panic(fmt.Sprintf(`Failed to start consensus state: %v

conS:
%+v

conR:
%+v`, err, conR.conS, conR))
	}
}

// GetChannels implements Reactor
func (conR *Reactor) GetChannels() []*p2p.ChannelDescriptor {
	// TODO optimize
	return []*p2p.ChannelDescriptor{
		{
			ID:                  StateChannel,
			Priority:            6,
			SendQueueCapacity:   100,
			RecvMessageCapacity: maxMsgSize,
		},
		{
			ID: DataChannel, // maybe split between gossiping current block and catchup stuff
			// once we gossip the whole block there's nothing left to send until next height or round
			Priority:            10,
			SendQueueCapacity:   100,
			RecvBufferCapacity:  50 * 4096,
			RecvMessageCapacity: maxMsgSize,
		},
		{
			ID:                  VoteChannel,
			Priority:            7,
			SendQueueCapacity:   100,
			RecvBufferCapacity:  100 * 100,
			RecvMessageCapacity: maxMsgSize,
		},
		{
			ID:                  VoteSetBitsChannel,
			Priority:            1,
			SendQueueCapacity:   2,
			RecvBufferCapacity:  1024,
			RecvMessageCapacity: maxMsgSize,
		},
		{
			ID:                  QrnChannel,
			Priority:            10,
			SendQueueCapacity:   100,
			RecvBufferCapacity:  50 * 4096,
			RecvMessageCapacity: maxMsgSize,
		},
		{
			ID:                  VrfChannel,
			Priority:            10,
			SendQueueCapacity:   100,
			RecvBufferCapacity:  50 * 4096,
			RecvMessageCapacity: maxMsgSize,
		},
		{
			ID:                  SettingSteeringMemberChannel,
			Priority:            10,
			SendQueueCapacity:   100,
			RecvBufferCapacity:  50 * 4096,
			RecvMessageCapacity: maxMsgSize,
		},
				{
			ID:                  CatchUpChannel,
			Priority:            10,
			SendQueueCapacity:   100,
			RecvBufferCapacity:  50 * 4096,
			RecvMessageCapacity: maxMsgSize,
		},
	}
}

// InitPeer implements Reactor by creating a state for the peer.
func (conR *Reactor) InitPeer(peer p2p.Peer) p2p.Peer {
	peerState := NewPeerState(peer).SetLogger(conR.Logger)
	peer.Set(types.PeerStateKey, peerState)
	return peer
}

// AddPeer implements Reactor by spawning multiple gossiping goroutines for the
// peer.
func (conR *Reactor) AddPeer(peer p2p.Peer) {
	if !conR.IsRunning() {
		return
	}

	peerState, ok := peer.Get(types.PeerStateKey).(*PeerState)
	if !ok {
		panic(fmt.Sprintf("peer %v has no state", peer))
	}

	// Begin routines for this peer.
	go conR.gossipDataRoutine(peer, peerState)
	go conR.gossipVotesRoutine(peer, peerState)
	go conR.queryMaj23Routine(peer, peerState)

	go conR.gossipQrnsRoutine(peer, peerState)
	go conR.gossipVrfsRoutine(peer, peerState)
	go conR.gossipSettingSteeringMemberRoutine(peer, peerState)

	// Send our state to peer.
	// If we're fast_syncing, broadcast a RoundStepMessage later upon SwitchToConsensus().
	if !conR.WaitSync() {
		conR.sendNewRoundStepMessage(peer)
	}
}

// RemovePeer is a noop.
func (conR *Reactor) RemovePeer(peer p2p.Peer, reason interface{}) {
	if !conR.IsRunning() {
		return
	}
	// TODO
	// ps, ok := peer.Get(PeerStateKey).(*PeerState)
	// if !ok {
	// 	panic(fmt.Sprintf("Peer %v has no state", peer))
	// }
	// ps.Disconnect()
}

// Receive implements Reactor
// NOTE: We process these messages even when we're fast_syncing.
// Messages affect either a peer state or the consensus state.
// Peer state updates can happen in parallel, but processing of
// proposals, block parts, and votes are ordered by the receiveRoutine
// NOTE: blocks on consensus state for proposals, block parts, and votes
func (conR *Reactor) Receive(chID byte, src p2p.Peer, msgBytes []byte) {
	if !conR.IsRunning() {
		conR.Logger.Debug("Receive", "src", src, "chId", chID, "bytes", msgBytes)
		return
	}

	switch chID {
	case CatchUpChannel:
		msg, _ := bc.DecodeMsg(msgBytes)
		switch msg := msg.(type) {
		case *bcproto.StateResponse:
			state, err := sm.SyncStateFromProto(msg.State)

			if err != nil {
				conR.Logger.Error("State content is invalid", "err", err)
				return
			}
			
			for _, catchupState := range conR.conS.CatchupStates {
				if (catchupState.LastBlockHeight == state.LastBlockHeight){
					return
				}
			}
			conR.conS.CatchupStates = append(conR.conS.CatchupStates, state)
		}
		return		
	}

	msg, err := decodeMsg(msgBytes)
	if err != nil {
		conR.Logger.Error("Error decoding message", "src", src, "chId", chID, "err", err)
		conR.Switch.StopPeerForError(src, err)
		return
	}

	if err = msg.ValidateBasic(); err != nil {
		conR.Logger.Error("Peer sent us invalid msg2", "peer", src, "msg", msg, "err", err)
		conR.Switch.StopPeerForError(src, err)
		return
	}

	conR.Logger.Debug("Receive", "src", src, "chId", chID, "msg", msg)

	// Get peer states
	ps, ok := src.Get(types.PeerStateKey).(*PeerState)
	if !ok {
		panic(fmt.Sprintf("Peer %v has no state", src))
	}

	switch chID {
	case StateChannel:
		switch msg := msg.(type) {
		case *NewRoundStepMessage:
			conR.conS.mtx.Lock()
			initialHeight := conR.conS.state.InitialHeight
			conR.conS.mtx.Unlock()
			if err = msg.ValidateHeight(initialHeight); err != nil {
				conR.Logger.Error("Peer sent us invalid msg1", "peer", src, "msg", msg, "err", err)
				conR.Switch.StopPeerForError(src, err)
				return
			}
			ps.ApplyNewRoundStepMessage(msg, conR.conS.state.ConsensusRound.ConsensusStartBlockHeight+int64(conR.conS.state.ConsensusRound.Period))
		case *NewValidBlockMessage:
			ps.ApplyNewValidBlockMessage(msg)
		case *HasVoteMessage:
			ps.ApplyHasVoteMessage(msg)
		case *HasSettingSteeringMemberMessage:
			ps.ApplyHasSettingSteeringMemberMessage(msg)
		case *VoteSetMaj23Message:
			cs := conR.conS
			cs.mtx.Lock()
			height, votes := cs.Height, cs.Votes
			cs.mtx.Unlock()
			if height != msg.Height {
				return
			}
			// Peer claims to have a maj23 for some BlockID at H,R,S,
			err := votes.SetPeerMaj23(msg.Round, msg.Type, ps.peer.ID(), msg.BlockID)
			if err != nil {
				conR.Switch.StopPeerForError(src, err)
				return
			}
			// Respond with a VoteSetBitsMessage showing which votes we have.
			// (and consequently shows which we don't have)
			var ourVotes *bits.BitArray
			switch msg.Type {
			case tmproto.PrevoteType:
				ourVotes = votes.Prevotes(msg.Round).BitArrayByBlockID(msg.BlockID)
			case tmproto.PrecommitType:
				ourVotes = votes.Precommits(msg.Round).BitArrayByBlockID(msg.BlockID)
			default:
				panic("Bad VoteSetBitsMessage field Type. Forgot to add a check in ValidateBasic?")
			}
			src.TrySend(VoteSetBitsChannel, MustEncode(&VoteSetBitsMessage{
				Height:  msg.Height,
				Round:   msg.Round,
				Type:    msg.Type,
				BlockID: msg.BlockID,
				Votes:   ourVotes,
			}))
		default:
			conR.Logger.Error(fmt.Sprintf("Unknown message type %v", reflect.TypeOf(msg)))
		}

	case DataChannel:
		if conR.WaitSync() {
			conR.Logger.Info("Ignoring message received during sync", "msg", msg)
			return
		}
		switch msg := msg.(type) {
		case *ProposalMessage:
			ps.SetHasProposal(msg.Proposal)
			conR.conS.peerMsgQueue <- msgInfo{msg, src.ID()}
		case *ProposalPOLMessage:
			ps.ApplyProposalPOLMessage(msg)
		case *BlockPartMessage:
			ps.SetHasProposalBlockPart(msg.Height, msg.Round, int(msg.Part.Index))
			conR.Metrics.BlockParts.With("peer_id", string(src.ID())).Add(1)
			conR.conS.peerMsgQueue <- msgInfo{msg, src.ID()}
		default:
			conR.Logger.Error(fmt.Sprintf("Unknown message type %v", reflect.TypeOf(msg)))
		}

	case QrnChannel:
		// if fast sync the message is added in catchup state
		if conR.WaitSync() {
			switch msg := msg.(type) {
				case *QrnMessage:
					conR.mtx.Lock()
					conR.tryAddCatchupQrnMessage(conR.chainID, msg)
					conR.mtx.Unlock()
			}
			
			return
		}
		switch msg := msg.(type) {
		case *HasQrnMessage:
			ps.ApplyHasQrnMessage(msg)
			
		case *QrnMessage:
			cs := conR.conS
			cs.mtx.RLock()
			height, standingMemberSize := cs.Height, cs.RoundState.StandingMemberSet.Size()
			cs.mtx.RUnlock()
			ps.EnsureQrnBitArrays(height, standingMemberSize)
			ps.SetHasQrn(msg.Qrn)

			cs.peerMsgQueue <- msgInfo{msg, src.ID()}

		default:
			// don't punish (leave room for soft upgrades)
			conR.Logger.Error(fmt.Sprintf("Unknown message type %v", reflect.TypeOf(msg)))
		}

	case VrfChannel:
		// if fast sync the message is added in catchup state
		if conR.WaitSync() {
			switch msg := msg.(type) {
				case *VrfMessage:
					conR.mtx.Lock()
					conR.tryAddCatchupVrfMessage(msg)
					conR.mtx.Unlock()
			}
			return
		}
		switch msg := msg.(type) {
		case *HasVrfMessage:
			ps.ApplyHasVrfMessage(msg)
		case *VrfMessage:
			cs := conR.conS
			cs.mtx.RLock()
			height, steeringMemberCandidate := cs.Height, cs.RoundState.SteeringMemberCandidateSet.Size()
			cs.mtx.RUnlock()
			ps.EnsureVrfBitArrays(height, steeringMemberCandidate)
			ps.SetHasVrf(msg.Vrf)

			cs.peerMsgQueue <- msgInfo{msg, src.ID()}

		default:
			// don't punish (leave room for soft upgrades)
			conR.Logger.Error(fmt.Sprintf("Unknown message type %v", reflect.TypeOf(msg)))
		}

	case SettingSteeringMemberChannel:
		// if fast sync the message is added in catchup state
		if conR.WaitSync() {
			switch msg := msg.(type) {
				case *SettingSteeringMemberMessage:
					conR.mtx.Lock()
					conR.tryAddCatchupSettingSteeringMemberMessage(conR.chainID, msg)
					conR.mtx.Unlock()
			}
			return
		}
		switch msg := msg.(type) {
		case *SettingSteeringMemberMessage:
			cs := conR.conS
			cs.mtx.RLock()
			cs.mtx.RUnlock()
			ps.SetHasSettingSteeringMember(msg.SettingSteeringMember.Height)

			cs.peerMsgQueue <- msgInfo{msg, src.ID()}

		default:
			// don't punish (leave room for soft upgrades)
			conR.Logger.Error(fmt.Sprintf("Unknown message type %v", reflect.TypeOf(msg)))
		}

	case VoteChannel:
		if conR.WaitSync() {
			conR.Logger.Info("Ignoring message received during sync", "msg", msg)
			return
		}
		
		switch msg := msg.(type) {
		case *VoteMessage:
			cs := conR.conS
			cs.mtx.RLock()
			height, valSize, lastCommitSize := cs.Height, cs.Validators.Size(), cs.LastCommit.Size()
			cs.mtx.RUnlock()
			ps.EnsureVoteBitArrays(height, valSize)
			ps.EnsureVoteBitArrays(height-1, lastCommitSize)
			ps.SetHasVote(msg.Vote)

			cs.peerMsgQueue <- msgInfo{msg, src.ID()}

		default:
			// don't punish (leave room for soft upgrades)
			conR.Logger.Error(fmt.Sprintf("Unknown message type %v", reflect.TypeOf(msg)))
		}

	case VoteSetBitsChannel:
		if conR.WaitSync() {
			conR.Logger.Info("Ignoring message received during sync", "msg", msg)
			return
		}
		switch msg := msg.(type) {
		case *VoteSetBitsMessage:
			cs := conR.conS
			cs.mtx.Lock()
			height, votes := cs.Height, cs.Votes
			cs.mtx.Unlock()

			if height == msg.Height {
				var ourVotes *bits.BitArray
				switch msg.Type {
				case tmproto.PrevoteType:
					ourVotes = votes.Prevotes(msg.Round).BitArrayByBlockID(msg.BlockID)
				case tmproto.PrecommitType:
					ourVotes = votes.Precommits(msg.Round).BitArrayByBlockID(msg.BlockID)
				default:
					panic("Bad VoteSetBitsMessage field Type. Forgot to add a check in ValidateBasic?")
				}
				ps.ApplyVoteSetBitsMessage(msg, ourVotes)
			} else {
				ps.ApplyVoteSetBitsMessage(msg, nil)
			}
		default:
			// don't punish (leave room for soft upgrades)
			conR.Logger.Error(fmt.Sprintf("Unknown message type %v", reflect.TypeOf(msg)))
		}

	default:
		conR.Logger.Error(fmt.Sprintf("Unknown chId %X", chID))
	}
}

// SetEventBus sets event bus.
func (conR *Reactor) SetEventBus(b *types.EventBus) {
	conR.eventBus = b
	conR.conS.SetEventBus(b)
}

// WaitSync returns whether the consensus reactor is waiting for state/fast sync.
func (conR *Reactor) WaitSync() bool {
	conR.mtx.RLock()
	defer conR.mtx.RUnlock()
	return conR.waitSync
}

//--------------------------------------

// subscribeToBroadcastEvents subscribes for new round steps and votes
// using internal pubsub defined on state to broadcast
// them to peers upon receiving.
func (conR *Reactor) subscribeToBroadcastEvents() {
	const subscriber = "consensus-reactor"
	if err := conR.conS.evsw.AddListenerForEvent(subscriber, types.EventNewRoundStep,
		func(data tmevents.EventData) {
			conR.broadcastNewRoundStepMessage(data.(*cstypes.RoundState))
		}); err != nil {
		conR.Logger.Error("Error adding listener for events", "err", err)
	}

	if err := conR.conS.evsw.AddListenerForEvent(subscriber, types.EventValidBlock,
		func(data tmevents.EventData) {
			conR.broadcastNewValidBlockMessage(data.(*cstypes.RoundState))
		}); err != nil {
		conR.Logger.Error("Error adding listener for events", "err", err)
	}

	if err := conR.conS.evsw.AddListenerForEvent(subscriber, types.EventVote,
		func(data tmevents.EventData) {
			conR.broadcastHasVoteMessage(data.(*types.Vote))
		}); err != nil {
		conR.Logger.Error("Error adding listener for events", "err", err)
	}

	if err := conR.conS.evsw.AddListenerForEvent(subscriber, types.EventQrn,
		func(data tmevents.EventData) {
			conR.broadcastQrnMessage(data.(*types.Qrn))
		}); err != nil {
		conR.Logger.Error("Error adding listener for events", "err", err)
	}

	if err := conR.conS.evsw.AddListenerForEvent(subscriber, types.EventVrf,
		func(data tmevents.EventData) {
			conR.broadcastVrfMessage(data.(*types.Vrf))
		}); err != nil {
		conR.Logger.Error("Error adding listener for events", "err", err)
	}

	if err := conR.conS.evsw.AddListenerForEvent(subscriber, types.EventSettingSteeringMember,
		func(data tmevents.EventData) {
			conR.broadcastSettingSteeringMemberMessage(data.(*types.SettingSteeringMember))
		}); err != nil {
		conR.Logger.Error("Error adding listener for events", "err", err)
	}
}

func (conR *Reactor) unsubscribeFromBroadcastEvents() {
	const subscriber = "consensus-reactor"
	conR.conS.evsw.RemoveListener(subscriber)
}

func (conR *Reactor) broadcastNewRoundStepMessage(rs *cstypes.RoundState) {
	nrsMsg := makeRoundStepMessage(rs)
	conR.Switch.Broadcast(StateChannel, MustEncode(nrsMsg))
}

func (conR *Reactor) broadcastNewValidBlockMessage(rs *cstypes.RoundState) {
	csMsg := &NewValidBlockMessage{
		Height:             rs.Height,
		Round:              rs.Round,
		BlockPartSetHeader: rs.ProposalBlockParts.Header(),
		BlockParts:         rs.ProposalBlockParts.BitArray(),
		IsCommit:           rs.Step == cstypes.RoundStepCommit,
	}
	conR.Switch.Broadcast(StateChannel, MustEncode(csMsg))
}

// Broadcasts HasVoteMessage to peers that care.
func (conR *Reactor) broadcastHasVoteMessage(vote *types.Vote) {
	msg := &HasVoteMessage{
		Height: vote.Height,
		Round:  vote.Round,
		Type:   vote.Type,
		Index:  vote.ValidatorIndex,
	}
	conR.Switch.Broadcast(StateChannel, MustEncode(msg))
	/*
		// TODO: Make this broadcast more selective.
		for _, peer := range conR.Switch.Peers().List() {
			ps, ok := peer.Get(PeerStateKey).(*PeerState)
			if !ok {
				panic(fmt.Sprintf("Peer %v has no state", peer))
			}
			prs := ps.GetRoundState()
			if prs.Height == vote.Height {
				// TODO: Also filter on round?
				peer.TrySend(StateChannel, struct{ ConsensusMessage }{msg})
			} else {
				// Height doesn't match
				// TODO: check a field, maybe CatchupCommitRound?
				// TODO: But that requires changing the struct field comment.
			}
		}
	*/
}

func (conR *Reactor) broadcastQrnMessage(qrn *types.Qrn) {
	msg := &QrnMessage{
		Qrn: qrn.Copy(),
	}
	conR.Switch.Broadcast(QrnChannel, MustEncode(msg))
}

func (conR *Reactor) broadcastVrfMessage(vrf *types.Vrf) {
	msg := &VrfMessage{
		Vrf: vrf.Copy(),
	}
	conR.Switch.Broadcast(VrfChannel, MustEncode(msg))
}

func (conR *Reactor) broadcastSettingSteeringMemberMessage(settingSteeringMember *types.SettingSteeringMember) {
	msg := &SettingSteeringMemberMessage{
		SettingSteeringMember: settingSteeringMember.Copy(),
	}
	conR.Switch.Broadcast(SettingSteeringMemberChannel, MustEncode(msg))
}

func makeRoundStepMessage(rs *cstypes.RoundState) (nrsMsg *NewRoundStepMessage) {
	nrsMsg = &NewRoundStepMessage{
		Height:                rs.Height,
		Round:                 rs.Round,
		Step:                  rs.Step,
		SecondsSinceStartTime: int64(time.Since(rs.StartTime).Seconds()),
		LastCommitRound:       rs.LastCommit.GetRound(),
	}
	return
}

func (conR *Reactor) sendNewRoundStepMessage(peer p2p.Peer) {
	rs := conR.getRoundState()
	nrsMsg := makeRoundStepMessage(rs)
	peer.Send(StateChannel, MustEncode(nrsMsg))
}

func (conR *Reactor) updateRoundStateRoutine() {
	t := time.NewTicker(100 * time.Microsecond)
	defer t.Stop()
	for range t.C {
		if !conR.IsRunning() {
			return
		}
		rs := conR.conS.GetRoundState()
		conR.mtx.Lock()
		conR.rs = rs
		conR.mtx.Unlock()
	}
}

func (conR *Reactor) getRoundState() *cstypes.RoundState {
	conR.mtx.RLock()
	defer conR.mtx.RUnlock()
	return conR.rs
}

func (conR *Reactor) gossipDataRoutine(peer p2p.Peer, ps *PeerState) {
	logger := conR.Logger.With("peer", peer)

	OUTER_LOOP:
		for {
			// Manage disconnects from self or peer.
			if !peer.IsRunning() || !conR.IsRunning() {
				logger.Info("Stopping gossipDataRoutine for peer")
				return
			}
			rs := conR.getRoundState()
			prs := ps.GetRoundState()

			// Send proposal Block parts?
			if rs.ProposalBlockParts.HasHeader(prs.ProposalBlockPartSetHeader) {
				if index, ok := rs.ProposalBlockParts.BitArray().Sub(prs.ProposalBlockParts.Copy()).PickRandom(); ok {
					part := rs.ProposalBlockParts.GetPart(index)
					msg := &BlockPartMessage{
						Height: rs.Height, // This tells peer that this part applies to us.
						Round:  rs.Round,  // This tells peer that this part applies to us.
						Part:   part,
					}
					logger.Debug("Sending block part", "height", prs.Height, "round", prs.Round)
					if peer.Send(DataChannel, MustEncode(msg)) {
						ps.SetHasProposalBlockPart(prs.Height, prs.Round, index)
					}
					continue OUTER_LOOP
				}
			}

			// If the peer is on a previous height that we have, help catch up.
			blockStoreBase := conR.conS.blockStore.Base()
			if blockStoreBase > 0 && 0 < prs.Height && prs.Height < rs.Height && prs.Height >= blockStoreBase {
				heightLogger := logger.With("height", prs.Height)

				// if we never received the commit message from the peer, the block parts wont be initialized
				if prs.ProposalBlockParts == nil {
					blockMeta := conR.conS.blockStore.LoadBlockMeta(prs.Height)
					if blockMeta == nil {
						heightLogger.Error("Failed to load block meta",
							"blockstoreBase", blockStoreBase, "blockstoreHeight", conR.conS.blockStore.Height())
						time.Sleep(conR.conS.config.PeerGossipSleepDuration)
					} else {
						ps.InitProposalBlockParts(blockMeta.BlockID.PartSetHeader)
					}
					// continue the loop since prs is a copy and not effected by this initialization
					continue OUTER_LOOP
				}

				// if the peer is on syncying
				if ( prs.Height < rs.Height - 1) {
					conR.respondStateToPeer(prs.Height, peer)
				}

				conR.gossipDataForCatchup(heightLogger, rs, prs, ps, peer)
				continue OUTER_LOOP
			}

			// If height and round don't match, sleep.
			if (rs.Height != prs.Height) || (rs.Round != prs.Round) {
				// logger.Info("Peer Height|Round mismatch, sleeping",
				// "peerHeight", prs.Height, "peerRound", prs.Round, "peer", peer)
				time.Sleep(conR.conS.config.PeerGossipSleepDuration)
				continue OUTER_LOOP
			}

			// By here, height and round match.
			// Proposal block parts were already matched and sent if any were wanted.
			// (These can match on hash so the round doesn't matter)
			// Now consider sending other things, like the Proposal itself.

			// Send Proposal && ProposalPOL BitArray?
			if rs.Proposal != nil && !prs.Proposal {
				// Proposal: share the proposal metadata with peer.
				{
					msg := &ProposalMessage{Proposal: rs.Proposal}
					logger.Debug("Sending proposal", "height", prs.Height, "round", prs.Round)
					if peer.Send(DataChannel, MustEncode(msg)) {
						// NOTE[ZM]: A peer might have received different proposal msg so this Proposal msg will be rejected!
						ps.SetHasProposal(rs.Proposal)
					}
				}
				// ProposalPOL: lets peer know which POL votes we have so far.
				// Peer must receive ProposalMessage first.
				// rs.Proposal was validated, so rs.Proposal.POLRound <= rs.Round,
				// so we definitely have rs.Votes.Prevotes(rs.Proposal.POLRound).
				if 0 <= rs.Proposal.POLRound {
					msg := &ProposalPOLMessage{
						Height:           rs.Height,
						ProposalPOLRound: rs.Proposal.POLRound,
						ProposalPOL:      rs.Votes.Prevotes(rs.Proposal.POLRound).BitArray(),
					}
					logger.Debug("Sending POL", "height", prs.Height, "round", prs.Round)
					peer.Send(DataChannel, MustEncode(msg))
				}
				continue OUTER_LOOP
			}

			// Nothing to do. Sleep.
			time.Sleep(conR.conS.config.PeerGossipSleepDuration)
			continue OUTER_LOOP
		}
}

func (conR *Reactor) gossipDataForCatchup(logger log.Logger, rs *cstypes.RoundState,
	prs *cstypes.PeerRoundState, ps *PeerState, peer p2p.Peer) {

	if index, ok := prs.ProposalBlockParts.Not().PickRandom(); ok {
		// Ensure that the peer's PartSetHeader is correct
		blockMeta := conR.conS.blockStore.LoadBlockMeta(prs.Height)
		if blockMeta == nil {
			logger.Error("Failed to load block meta", "ourHeight", rs.Height,
				"blockstoreBase", conR.conS.blockStore.Base(), "blockstoreHeight", conR.conS.blockStore.Height())
			time.Sleep(conR.conS.config.PeerGossipSleepDuration)
			return
		} else if !blockMeta.BlockID.PartSetHeader.Equals(prs.ProposalBlockPartSetHeader) {
			logger.Info("Peer ProposalBlockPartSetHeader mismatch, sleeping",
				"blockPartSetHeader", blockMeta.BlockID.PartSetHeader, "peerBlockPartSetHeader", prs.ProposalBlockPartSetHeader)
			time.Sleep(conR.conS.config.PeerGossipSleepDuration)
			return
		}
		// Load the part
		part := conR.conS.blockStore.LoadBlockPart(prs.Height, index)
		if part == nil {
			logger.Error("Could not load part", "index", index,
				"blockPartSetHeader", blockMeta.BlockID.PartSetHeader, "peerBlockPartSetHeader", prs.ProposalBlockPartSetHeader)
			time.Sleep(conR.conS.config.PeerGossipSleepDuration)
			return
		}
		// Send the part
		msg := &BlockPartMessage{
			Height: prs.Height, // Not our height, so it doesn't matter.
			Round:  prs.Round,  // Not our height, so it doesn't matter.
			Part:   part,
		}
		logger.Debug("Sending block part for catchup", "round", prs.Round, "index", index)
		if peer.Send(DataChannel, MustEncode(msg)) {
			ps.SetHasProposalBlockPart(prs.Height, prs.Round, index)
		} else {
			logger.Debug("Sending block part for catchup failed")
		}
		return
	}
	//  logger.Info("No parts to send in catch-up, sleeping")
	time.Sleep(conR.conS.config.PeerGossipSleepDuration)
}

func (conR *Reactor) gossipVotesRoutine(peer p2p.Peer, ps *PeerState) {
	logger := conR.Logger.With("peer", peer)

	// Simple hack to throttle logs upon sleep.
	var sleeping = 0

OUTER_LOOP:
	for {
		// Manage disconnects from self or peer.
		if !peer.IsRunning() || !conR.IsRunning() {
			logger.Info("Stopping gossipVotesRoutine for peer")
			return
		}
		rs := conR.getRoundState()
		prs := ps.GetRoundState()

		switch sleeping {
		case 1: // First sleep
			sleeping = 2
		case 2: // No more sleep
			sleeping = 0
		}

		// logger.Debug("gossipVotesRoutine", "rsHeight", rs.Height, "rsRound", rs.Round,
		// "prsHeight", prs.Height, "prsRound", prs.Round, "prsStep", prs.Step)

		// If height matches, then send LastCommit, Prevotes, Precommits.
		if rs.Height == prs.Height {
			heightLogger := logger.With("height", prs.Height)
			if conR.gossipVotesForHeight(heightLogger, rs, prs, ps) {
				continue OUTER_LOOP
			}
		}

		// Special catchup logic.
		// If peer is lagging by height 1, send LastCommit.
		if prs.Height != 0 && rs.Height == prs.Height+1 {
			if ps.PickSendVote(rs.LastCommit) {
				logger.Debug("Picked rs.LastCommit to send", "height", prs.Height)
				continue OUTER_LOOP
			}
		}

		// Catchup logic
		// If peer is lagging by more than 1, send Commit.
		blockStoreBase := conR.conS.blockStore.Base()
		if blockStoreBase > 0 && prs.Height != 0 && rs.Height >= prs.Height+2 && prs.Height >= blockStoreBase {
			// Load the block commit for prs.Height,
			// which contains precommit signatures for prs.Height.
			if commit := conR.conS.blockStore.LoadBlockCommit(prs.Height); commit != nil {
				if ps.PickSendVote(commit) {
					logger.Debug("Picked Catchup commit to send", "height", prs.Height)
					continue OUTER_LOOP
				}
			}
		}

		if sleeping == 0 {
			// We sent nothing. Sleep...
			sleeping = 1
			logger.Debug("No votes to send, sleeping", "rs.Height", rs.Height, "prs.Height", prs.Height,
				"localPV", rs.Votes.Prevotes(rs.Round).BitArray(), "peerPV", prs.Prevotes,
				"localPC", rs.Votes.Precommits(rs.Round).BitArray(), "peerPC", prs.Precommits)
		} else if sleeping == 2 {
			// Continued sleep...
			sleeping = 1
		}

		time.Sleep(conR.conS.config.PeerGossipSleepDuration)
		continue OUTER_LOOP
	}
}

func (conR *Reactor) gossipVotesForHeight(
	logger log.Logger,
	rs *cstypes.RoundState,
	prs *cstypes.PeerRoundState,
	ps *PeerState,
) bool {

	// If there are lastCommits to send...
	if prs.Step == cstypes.RoundStepNewHeight {
		if ps.PickSendVote(rs.LastCommit) {
			logger.Debug("Picked rs.LastCommit to send")
			return true
		}
	}
	// If there are POL prevotes to send...
	if prs.Step <= cstypes.RoundStepPropose && prs.Round != -1 && prs.Round <= rs.Round && prs.ProposalPOLRound != -1 {
		if polPrevotes := rs.Votes.Prevotes(prs.ProposalPOLRound); polPrevotes != nil {
			if ps.PickSendVote(polPrevotes) {
				logger.Debug("Picked rs.Prevotes(prs.ProposalPOLRound) to send",
					"round", prs.ProposalPOLRound)
				return true
			}
		}
	}
	// If there are prevotes to send...
	if prs.Step <= cstypes.RoundStepPrevoteWait && prs.Round != -1 && prs.Round <= rs.Round {
		if ps.PickSendVote(rs.Votes.Prevotes(prs.Round)) {
			logger.Debug("Picked rs.Prevotes(prs.Round) to send", "round", prs.Round)
			return true
		}
	}
	// If there are precommits to send...
	if prs.Step <= cstypes.RoundStepPrecommitWait && prs.Round != -1 && prs.Round <= rs.Round {
		if ps.PickSendVote(rs.Votes.Precommits(prs.Round)) {
			logger.Debug("Picked rs.Precommits(prs.Round) to send", "round", prs.Round)
			return true
		}
	}
	// If there are prevotes to send...Needed because of validBlock mechanism
	if prs.Round != -1 && prs.Round <= rs.Round {
		if ps.PickSendVote(rs.Votes.Prevotes(prs.Round)) {
			logger.Debug("Picked rs.Prevotes(prs.Round) to send", "round", prs.Round)
			return true
		}
	}
	// If there are POLPrevotes to send...
	if prs.ProposalPOLRound != -1 {
		if polPrevotes := rs.Votes.Prevotes(prs.ProposalPOLRound); polPrevotes != nil {
			if ps.PickSendVote(polPrevotes) {
				logger.Debug("Picked rs.Prevotes(prs.ProposalPOLRound) to send",
					"round", prs.ProposalPOLRound)
				return true
			}
		}
	}

	return false
}

// NOTE: `queryMaj23Routine` has a simple crude design since it only comes
// into play for liveness when there's a signature DDoS attack happening.
func (conR *Reactor) queryMaj23Routine(peer p2p.Peer, ps *PeerState) {
	logger := conR.Logger.With("peer", peer)

OUTER_LOOP:
	for {
		// Manage disconnects from self or peer.
		if !peer.IsRunning() || !conR.IsRunning() {
			logger.Info("Stopping queryMaj23Routine for peer")
			return
		}

		// Maybe send Height/Round/Prevotes
		{
			rs := conR.getRoundState()
			prs := ps.GetRoundState()
			if rs.Height == prs.Height {
				if maj23, ok := rs.Votes.Prevotes(prs.Round).TwoThirdsMajority(); ok {
					peer.TrySend(StateChannel, MustEncode(&VoteSetMaj23Message{
						Height:  prs.Height,
						Round:   prs.Round,
						Type:    tmproto.PrevoteType,
						BlockID: maj23,
					}))
					time.Sleep(conR.conS.config.PeerQueryMaj23SleepDuration)
				}
			}
		}

		// Maybe send Height/Round/Precommits
		{
			rs := conR.getRoundState()
			prs := ps.GetRoundState()
			if rs.Height == prs.Height {
				if maj23, ok := rs.Votes.Precommits(prs.Round).TwoThirdsMajority(); ok {
					peer.TrySend(StateChannel, MustEncode(&VoteSetMaj23Message{
						Height:  prs.Height,
						Round:   prs.Round,
						Type:    tmproto.PrecommitType,
						BlockID: maj23,
					}))
					time.Sleep(conR.conS.config.PeerQueryMaj23SleepDuration)
				}
			}
		}

		// Maybe send Height/Round/ProposalPOL
		{
			rs := conR.getRoundState()
			prs := ps.GetRoundState()
			if rs.Height == prs.Height && prs.ProposalPOLRound >= 0 {
				if maj23, ok := rs.Votes.Prevotes(prs.ProposalPOLRound).TwoThirdsMajority(); ok {
					peer.TrySend(StateChannel, MustEncode(&VoteSetMaj23Message{
						Height:  prs.Height,
						Round:   prs.ProposalPOLRound,
						Type:    tmproto.PrevoteType,
						BlockID: maj23,
					}))
					time.Sleep(conR.conS.config.PeerQueryMaj23SleepDuration)
				}
			}
		}

		// Little point sending LastCommitRound/LastCommit,
		// These are fleeting and non-blocking.

		// Maybe send Height/CatchupCommitRound/CatchupCommit.
		{
			prs := ps.GetRoundState()
			if prs.CatchupCommitRound != -1 && prs.Height > 0 && prs.Height <= conR.conS.blockStore.Height() &&
				prs.Height >= conR.conS.blockStore.Base() {
				if commit := conR.conS.LoadCommit(prs.Height); commit != nil {
					peer.TrySend(StateChannel, MustEncode(&VoteSetMaj23Message{
						Height:  prs.Height,
						Round:   commit.Round,
						Type:    tmproto.PrecommitType,
						BlockID: commit.BlockID,
					}))
					time.Sleep(conR.conS.config.PeerQueryMaj23SleepDuration)
				}
			}
		}

		time.Sleep(conR.conS.config.PeerQueryMaj23SleepDuration)

		continue OUTER_LOOP
	}
}

func (conR *Reactor) peerStatsRoutine() {
	for {
		if !conR.IsRunning() {
			conR.Logger.Info("Stopping peerStatsRoutine")
			return
		}

		select {
		case msg := <-conR.conS.statsMsgQueue:
			// Get peer
			peer := conR.Switch.Peers().Get(msg.PeerID)
			if peer == nil {
				conR.Logger.Debug("Attempt to update stats for non-existent peer",
					"peer", msg.PeerID)
				continue
			}
			// Get peer state
			ps, ok := peer.Get(types.PeerStateKey).(*PeerState)
			if !ok {
				panic(fmt.Sprintf("Peer %v has no state", peer))
			}
			switch msg.Msg.(type) {
			case *VoteMessage:
				if numVotes := ps.RecordVote(); numVotes%votesToContributeToBecomeGoodPeer == 0 {
					conR.Switch.MarkPeerAsGood(peer)
				}
			case *BlockPartMessage:
				if numParts := ps.RecordBlockPart(); numParts%blocksToContributeToBecomeGoodPeer == 0 {
					conR.Switch.MarkPeerAsGood(peer)
				}
			}
		case <-conR.conS.Quit():
			return

		case <-conR.Quit():
			return
		}
	}
}

// String returns a string representation of the Reactor.
// NOTE: For now, it is just a hard-coded string to avoid accessing unprotected shared variables.
// TODO: improve!
func (conR *Reactor) String() string {
	// better not to access shared variables
	return "ConsensusReactor" // conR.StringIndented("")
}

// StringIndented returns an indented string representation of the Reactor
func (conR *Reactor) StringIndented(indent string) string {
	s := "ConsensusReactor{\n"
	s += indent + "  " + conR.conS.StringIndented(indent+"  ") + "\n"
	for _, peer := range conR.Switch.Peers().List() {
		ps, ok := peer.Get(types.PeerStateKey).(*PeerState)
		if !ok {
			panic(fmt.Sprintf("Peer %v has no state", peer))
		}
		s += indent + "  " + ps.StringIndented(indent+"  ") + "\n"
	}
	s += indent + "}"
	return s
}

// ReactorMetrics sets the metrics
func ReactorMetrics(metrics *Metrics) ReactorOption {
	return func(conR *Reactor) { conR.Metrics = metrics }
}

//-----------------------------------------------------------------------------
// Messages
