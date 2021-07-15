package consensus

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"time"

	"github.com/gogo/protobuf/proto"

	cfg "github.com/reapchain/reapchain-core/config"
	cstypes "github.com/reapchain/reapchain-core/consensus/types"
	"github.com/reapchain/reapchain-core/crypto"
	tmevents "github.com/reapchain/reapchain-core/libs/events"
	"github.com/reapchain/reapchain-core/libs/fail"
	tmjson "github.com/reapchain/reapchain-core/libs/json"
	"github.com/reapchain/reapchain-core/libs/log"
	tmmath "github.com/reapchain/reapchain-core/libs/math"
	tmos "github.com/reapchain/reapchain-core/libs/os"
	"github.com/reapchain/reapchain-core/libs/service"
	tmsync "github.com/reapchain/reapchain-core/libs/sync"
	"github.com/reapchain/reapchain-core/p2p"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
	sm "github.com/reapchain/reapchain-core/state"
	"github.com/reapchain/reapchain-core/types"
	tmtime "github.com/reapchain/reapchain-core/types/time"
)

// Consensus sentinel errors
var (
	ErrInvalidProposalSignature   = errors.New("error invalid proposal signature")
	ErrInvalidProposalPOLRound    = errors.New("error invalid proposal POL round")
	ErrAddingVote                 = errors.New("error adding vote")
	ErrAddingQrn                  = errors.New("error adding qrn")
	ErrSignatureFoundInPastBlocks = errors.New("found signature from the same key")

	errPubKeyIsNotSet = errors.New("pubkey is not set. Look for \"Can't get private validator pubkey\" errors")
)

var msgQueueSize = 1000

// msgs from the reactor which may update the state
type MsgInfo struct {
	Msg    Message `json:"msg"`
	PeerID p2p.ID  `json:"peer_key"`
}

// internally generated messages which may update the state
type timeoutInfo struct {
	Duration time.Duration         `json:"duration"`
	Height   int64                 `json:"height"`
	Round    int32                 `json:"round"`
	Step     cstypes.RoundStepType `json:"step"`
}

func (ti *timeoutInfo) String() string {
	return fmt.Sprintf("%v ; %d/%d %v", ti.Duration, ti.Height, ti.Round, ti.Step)
}

// interface to the mempool
type txNotifier interface {
	TxsAvailable() <-chan struct{}
}

// interface to the evidence pool
type evidencePool interface {
	// reports conflicting votes to the evidence pool to be processed into evidence
	ReportConflictingVotes(voteA, voteB *types.Vote)
}

// State handles execution of the consensus algorithm.
// It processes votes and proposals, and upon reaching agreement,
// commits blocks to the chain and executes them against the application.
// The internal state machine receives input from peers, the internal validator, and from a timer.
type State struct {
	service.BaseService

	// config details
	config        *cfg.ConsensusConfig
	privValidator types.PrivValidator // for signing votes

	// store blocks and commits
	blockStore sm.BlockStore

	// create and execute blocks
	blockExec *sm.BlockExecutor

	// notify us if txs are available
	txNotifier txNotifier

	// add evidence to the pool
	// when it's detected
	evpool evidencePool

	// internal state
	mtx tmsync.RWMutex
	cstypes.RoundState
	state sm.State // State until height-1.
	// privValidator pubkey, memoized for the duration of one block
	// to avoid extra requests to HSM
	privValidatorPubKey crypto.PubKey

	// state changes may be triggered by: msgs from peers,
	// msgs from ourself, or by timeouts
	peerMsgQueue     chan MsgInfo
	internalMsgQueue chan MsgInfo
	timeoutTicker    TimeoutTicker

	// information about about added votes and block parts are written on this channel
	// so statistics can be computed by reactor
	statsMsgQueue chan MsgInfo

	// we use eventBus to trigger msg broadcasts in the reactor,
	// and to notify external subscribers, eg. through a websocket
	eventBus *types.EventBus

	// a Write-Ahead Log ensures we can recover from any kind of crash
	// and helps us avoid signing conflicting votes
	wal          WAL
	replayMode   bool // so we don't log signing errors during replay
	doWALCatchup bool // determines if we even try to do the catchup

	// for tests where we want to limit the number of transitions the state makes
	nSteps int

	// some functions can be overwritten for testing
	decideProposal func(height int64, round int32)
	doPrevote      func(height int64, round int32)
	setProposal    func(proposal *types.Proposal) error

	// closed when we finish shutting down
	done chan struct{}

	// synchronous pubsub between consensus state and reactor.
	// state only emits EventNewRoundStep and EventVote
	evsw tmevents.EventSwitch

	// for reporting metrics
	metrics *Metrics
}

// StateOption sets an optional parameter on the State.
type StateOption func(*State)

// NewState returns a new State.
func NewState(
	config *cfg.ConsensusConfig,
	state sm.State,
	blockExec *sm.BlockExecutor,
	blockStore sm.BlockStore,
	txNotifier txNotifier,
	evpool evidencePool,
	options ...StateOption,
) *State {
	consensusState := &State{
		config:           config,
		blockExec:        blockExec,
		blockStore:       blockStore,
		txNotifier:       txNotifier,
		peerMsgQueue:     make(chan MsgInfo, msgQueueSize),
		internalMsgQueue: make(chan MsgInfo, msgQueueSize),
		timeoutTicker:    NewTimeoutTicker(),
		statsMsgQueue:    make(chan MsgInfo, msgQueueSize),
		done:             make(chan struct{}),
		doWALCatchup:     true,
		wal:              nilWAL{},
		evpool:           evpool,
		evsw:             tmevents.NewEventSwitch(),
		metrics:          NopMetrics(),
	}

	// set function defaults (may be overwritten before calling Start)
	consensusState.decideProposal = consensusState.defaultDecideProposal
	consensusState.doPrevote = consensusState.defaultDoPrevote
	consensusState.setProposal = consensusState.defaultSetProposal

	// We have no votes, so reconstruct LastCommit from SeenCommit.
	if state.LastBlockHeight > 0 {
		consensusState.reconstructLastCommit(state)
	}

	//fmt.Println("stompesi-으어", state.Qrns)
	consensusState.updateToState(state)

	// NOTE: we do not call scheduleRound0 yet, we do that upon Start()

	consensusState.BaseService = *service.NewBaseService(nil, "State", consensusState)
	for _, option := range options {
		option(consensusState)
	}

	return consensusState
}

// SetLogger implements Service.
func (consensusState *State) SetLogger(l log.Logger) {
	consensusState.BaseService.Logger = l
	consensusState.timeoutTicker.SetLogger(l)
}

// SetEventBus sets event bus.
func (consensusState *State) SetEventBus(b *types.EventBus) {
	consensusState.eventBus = b
	consensusState.blockExec.SetEventBus(b)
}

// StateMetrics sets the metrics.
func StateMetrics(metrics *Metrics) StateOption {
	return func(consensusState *State) { consensusState.metrics = metrics }
}

// String returns a string.
func (consensusState *State) String() string {
	// better not to access shared variables
	return "ConsensusState"
}

// GetState returns a copy of the chain state.
func (consensusState *State) GetState() sm.State {
	consensusState.mtx.RLock()
	defer consensusState.mtx.RUnlock()
	return consensusState.state.Copy()
}

// GetLastHeight returns the last height committed.
// If there were no blocks, returns 0.
func (consensusState *State) GetLastHeight() int64 {
	consensusState.mtx.RLock()
	defer consensusState.mtx.RUnlock()
	return consensusState.RoundState.Height - 1
}

// GetRoundState returns a shallow copy of the internal consensus state.
func (consensusState *State) GetRoundState() *cstypes.RoundState {
	consensusState.mtx.RLock()
	rs := consensusState.RoundState // copy
	consensusState.mtx.RUnlock()
	return &rs
}

// GetRoundStateJSON returns a json of RoundState.
func (consensusState *State) GetRoundStateJSON() ([]byte, error) {
	consensusState.mtx.RLock()
	defer consensusState.mtx.RUnlock()
	return tmjson.Marshal(consensusState.RoundState)
}

// GetRoundStateSimpleJSON returns a json of RoundStateSimple
func (consensusState *State) GetRoundStateSimpleJSON() ([]byte, error) {
	consensusState.mtx.RLock()
	defer consensusState.mtx.RUnlock()
	return tmjson.Marshal(consensusState.RoundState.RoundStateSimple())
}

// GetValidators returns a copy of the current validators.
func (consensusState *State) GetValidators() (int64, []*types.Validator) {
	consensusState.mtx.RLock()
	defer consensusState.mtx.RUnlock()
	return consensusState.state.LastBlockHeight, consensusState.state.Validators.Copy().Validators
}

// GetStandingMembers returns a copy of the current validators.
func (consensusState *State) GetStandingMembers() (int64, []*types.StandingMember) {
	consensusState.mtx.RLock()
	defer consensusState.mtx.RUnlock()
	return consensusState.state.LastBlockHeight, consensusState.state.StandingMemberSet.Copy().StandingMembers
}

func (consensusState *State) GetQrns() (int64, []*types.Qrn) {
	consensusState.mtx.RLock()
	defer consensusState.mtx.RUnlock()
	return consensusState.state.LastBlockHeight, consensusState.state.QrnSet.Copy().Qrns
}

// SetPrivValidator sets the private validator account for signing votes. It
// immediately requests pubkey and caches it.
func (consensusState *State) SetPrivValidator(priv types.PrivValidator) {
	consensusState.mtx.Lock()
	defer consensusState.mtx.Unlock()

	consensusState.privValidator = priv

	if err := consensusState.updatePrivValidatorPubKey(); err != nil {
		consensusState.Logger.Error("failed to get private validator pubkey", "err", err)
	}
}

// SetTimeoutTicker sets the local timer. It may be useful to overwrite for
// testing.
func (consensusState *State) SetTimeoutTicker(timeoutTicker TimeoutTicker) {
	consensusState.mtx.Lock()
	consensusState.timeoutTicker = timeoutTicker
	consensusState.mtx.Unlock()
}

// LoadCommit loads the commit for a given height.
func (consensusState *State) LoadCommit(height int64) *types.Commit {
	consensusState.mtx.RLock()
	defer consensusState.mtx.RUnlock()

	if height == consensusState.blockStore.Height() {
		return consensusState.blockStore.LoadSeenCommit(height)
	}

	return consensusState.blockStore.LoadBlockCommit(height)
}

// OnStart loads the latest state via the WAL, and starts the timeout and
// receive routines.
func (consensusState *State) OnStart() error {
	// We may set the WAL in testing before calling Start, so only OpenWAL if its
	// still the nilWAL.
	if _, ok := consensusState.wal.(nilWAL); ok {
		if err := consensusState.loadWalFile(); err != nil {
			return err
		}
	}

	// We may have lost some votes if the process crashed reload from consensus
	// log to catchup.
	if consensusState.doWALCatchup {
		repairAttempted := false

	LOOP:
		for {
			err := consensusState.catchupReplay(consensusState.Height)
			switch {
			case err == nil:
				break LOOP

			case !IsDataCorruptionError(err):
				consensusState.Logger.Error("error on catchup replay; proceeding to start state anyway", "err", err)
				break LOOP

			case repairAttempted:
				return err
			}

			consensusState.Logger.Error("the WAL file is corrupted; attempting repair", "err", err)

			// 1) prep work
			if err := consensusState.wal.Stop(); err != nil {
				return err
			}

			repairAttempted = true

			// 2) backup original WAL file
			corruptedFile := fmt.Sprintf("%s.CORRUPTED", consensusState.config.WalFile())
			if err := tmos.CopyFile(consensusState.config.WalFile(), corruptedFile); err != nil {
				return err
			}

			consensusState.Logger.Debug("backed up WAL file", "src", consensusState.config.WalFile(), "dst", corruptedFile)

			// 3) try to repair (WAL file will be overwritten!)
			if err := repairWalFile(corruptedFile, consensusState.config.WalFile()); err != nil {
				consensusState.Logger.Error("the WAL repair failed", "err", err)
				return err
			}

			consensusState.Logger.Info("successful WAL repair")

			// reload WAL file
			if err := consensusState.loadWalFile(); err != nil {
				return err
			}
		}
	}

	if err := consensusState.evsw.Start(); err != nil {
		return err
	}

	// we need the timeoutRoutine for replay so
	// we don't block on the tick chan.
	// NOTE: we will get a build up of garbage go routines
	// firing on the tockChan until the receiveRoutine is started
	// to deal with them (by that point, at most one will be valid)
	if err := consensusState.timeoutTicker.Start(); err != nil {
		return err
	}

	// Double Signing Risk Reduction
	if err := consensusState.checkDoubleSigningRisk(consensusState.Height); err != nil {
		return err
	}

	// now start the receiveRoutine
	go consensusState.receiveRoutine(0)

	// schedule the first round!
	// use GetRoundState so we don't race the receiveRoutine for access
	consensusState.scheduleRound0(consensusState.GetRoundState())

	return nil
}

// timeoutRoutine: receive requests for timeouts on tickChan and fire timeouts on tockChan
// receiveRoutine: serializes processing of proposoals, block parts, votes; coordinates state transitions
func (consensusState *State) startRoutines(maxSteps int) {
	err := consensusState.timeoutTicker.Start()
	if err != nil {
		consensusState.Logger.Error("failed to start timeout ticker", "err", err)
		return
	}

	go consensusState.receiveRoutine(maxSteps)
}

// loadWalFile loads WAL data from file. It overwrites consensusState.wal.
func (consensusState *State) loadWalFile() error {
	wal, err := consensusState.OpenWAL(consensusState.config.WalFile())
	if err != nil {
		consensusState.Logger.Error("failed to load state WAL", "err", err)
		return err
	}

	consensusState.wal = wal
	return nil
}

// OnStop implements service.Service.
func (consensusState *State) OnStop() {
	if err := consensusState.evsw.Stop(); err != nil {
		consensusState.Logger.Error("failed trying to stop eventSwitch", "error", err)
	}

	if err := consensusState.timeoutTicker.Stop(); err != nil {
		consensusState.Logger.Error("failed trying to stop timeoutTicket", "error", err)
	}
	// WAL is stopped in receiveRoutine.
}

// Wait waits for the the main routine to return.
// NOTE: be sure to Stop() the event switch and drain
// any event channels or this may deadlock
func (consensusState *State) Wait() {
	<-consensusState.done
}

// OpenWAL opens a file to log all consensus messages and timeouts for
// deterministic accountability.
func (consensusState *State) OpenWAL(walFile string) (WAL, error) {
	wal, err := NewWAL(walFile)
	if err != nil {
		consensusState.Logger.Error("failed to open WAL", "file", walFile, "err", err)
		return nil, err
	}

	wal.SetLogger(consensusState.Logger.With("wal", walFile))

	if err := wal.Start(); err != nil {
		consensusState.Logger.Error("failed to start WAL", "err", err)
		return nil, err
	}

	return wal, nil
}

//------------------------------------------------------------
// internal functions for managing the state

func (consensusState *State) updateHeight(height int64) {
	consensusState.metrics.Height.Set(float64(height))
	consensusState.Height = height
}

func (consensusState *State) updateRoundStep(round int32, step cstypes.RoundStepType) {
	consensusState.Round = round
	consensusState.Step = step
}

// enterNewRound(height, 0) at consensusState.StartTime.
func (consensusState *State) scheduleRound0(rs *cstypes.RoundState) {
	// consensusState.Logger.Info("scheduleRound0", "now", tmtime.Now(), "startTime", consensusState.StartTime)
	sleepDuration := rs.StartTime.Sub(tmtime.Now())
	consensusState.scheduleTimeout(sleepDuration, rs.Height, 0, cstypes.RoundStepNewHeight)
}

// Attempt to schedule a timeout (by sending timeoutInfo on the tickChan)
func (consensusState *State) scheduleTimeout(duration time.Duration, height int64, round int32, step cstypes.RoundStepType) {
	consensusState.timeoutTicker.ScheduleTimeout(timeoutInfo{duration, height, round, step})
}

// send a msg into the receiveRoutine regarding our own proposal, block part, or vote
func (consensusState *State) sendInternalMessage(msgInfo MsgInfo) {
	select {
	case consensusState.internalMsgQueue <- msgInfo:
	default:
		// NOTE: using the go-routine means our votes can
		// be processed out of order.
		// TODO: use CList here for strict determinism and
		// attempt push to internalMsgQueue in receiveRoutine
		consensusState.Logger.Debug("internal msg queue is full; using a go-routine")
		go func() { consensusState.internalMsgQueue <- msgInfo }()
	}
}

// Reconstruct LastCommit from SeenCommit, which we saved along with the block,
// (which happens even before saving the state)
func (consensusState *State) reconstructLastCommit(state sm.State) {
	seenCommit := consensusState.blockStore.LoadSeenCommit(state.LastBlockHeight)
	if seenCommit == nil {
		panic(fmt.Sprintf(
			"failed to reconstruct last commit; seen commit for height %v not found",
			state.LastBlockHeight,
		))
	}

	lastPrecommits := types.CommitToVoteSet(state.ChainID, seenCommit, state.LastValidators)
	if !lastPrecommits.HasTwoThirdsMajority() {
		panic("failed to reconstruct last commit; does not have +2/3 maj")
	}

	consensusState.LastCommit = lastPrecommits
}

// Updates State and increments height to match that of state.
// The round becomes 0 and consensusState.Step becomes cstypes.RoundStepNewHeight.
func (consensusState *State) updateToState(state sm.State) {
	//fmt.Println("stompesi-jjjjjjkkkk", state.Qrns)

	if consensusState.CommitRound > -1 && 0 < consensusState.Height && consensusState.Height != state.LastBlockHeight {
		panic(fmt.Sprintf(
			"updateToState() expected state height of %v but found %v",
			consensusState.Height, state.LastBlockHeight,
		))
	}

	if !consensusState.state.IsEmpty() {
		if consensusState.state.LastBlockHeight > 0 && consensusState.state.LastBlockHeight+1 != consensusState.Height {
			// This might happen when someone else is mutating consensusState.state.
			// Someone forgot to pass in state.Copy() somewhere?!
			panic(fmt.Sprintf(
				"inconsistent consensusState.state.LastBlockHeight+1 %v vs consensusState.Height %v",
				consensusState.state.LastBlockHeight+1, consensusState.Height,
			))
		}
		if consensusState.state.LastBlockHeight > 0 && consensusState.Height == consensusState.state.InitialHeight {
			panic(fmt.Sprintf(
				"inconsistent consensusState.state.LastBlockHeight %v, expected 0 for initial height %v",
				consensusState.state.LastBlockHeight, consensusState.state.InitialHeight,
			))
		}

		// If state isn't further out than consensusState.state, just ignore.
		// This happens when SwitchToConsensus() is called in the reactor.
		// We don't want to reset e.g. the Votes, but we still want to
		// signal the new round step, because other services (eg. txNotifier)
		// depend on having an up-to-date peer state!
		if state.LastBlockHeight <= consensusState.state.LastBlockHeight {
			consensusState.Logger.Debug(
				"ignoring updateToState()",
				"new_height", state.LastBlockHeight+1,
				"old_height", consensusState.state.LastBlockHeight+1,
			)
			consensusState.newStep()
			return
		}
	}

	// Reset fields based on state.
	validators := state.Validators
	//fmt.Println("stompesi-last", state.StandingMembers)

	switch {
	case state.LastBlockHeight == 0: // Very first commit should be empty.
		consensusState.LastCommit = (*types.VoteSet)(nil)
	case consensusState.CommitRound > -1 && consensusState.Votes != nil: // Otherwise, use consensusState.Votes
		if !consensusState.Votes.Precommits(consensusState.CommitRound).HasTwoThirdsMajority() {
			panic(fmt.Sprintf(
				"wanted to form a commit, but precommits (H/R: %d/%d) didn't have 2/3+: %v",
				state.LastBlockHeight, consensusState.CommitRound, consensusState.Votes.Precommits(consensusState.CommitRound),
			))
		}

		consensusState.LastCommit = consensusState.Votes.Precommits(consensusState.CommitRound)

	case consensusState.LastCommit == nil:
		// NOTE: when Reapchain starts, it has no votes. reconstructLastCommit
		// must be called to reconstruct LastCommit from SeenCommit.
		panic(fmt.Sprintf(
			"last commit cannot be empty after initial block (H:%d)",
			state.LastBlockHeight+1,
		))
	}

	// Next desired block height
	height := state.LastBlockHeight + 1
	if height == 1 {
		height = state.InitialHeight
		state.ConsensusRound = types.NewConsensusRound(0, 0)
	}

	if state.ConsensusRound.ConsensusStartBlockHeight+state.ConsensusRound.Peorid == height {
		state.ConsensusRound.ConsensusStartBlockHeight = height
		state.QrnSet = types.NewQrnSet(height, state.StandingMemberSet, nil)
	}

	// RoundState fields
	consensusState.updateHeight(height)
	consensusState.updateRoundStep(0, cstypes.RoundStepNewHeight)

	if consensusState.CommitTime.IsZero() {
		// "Now" makes it easier to sync up dev nodes.
		// We add timeoutCommit to allow transactions
		// to be gathered for the first block.
		// And alternative solution that relies on clocks:
		// consensusState.StartTime = state.LastBlockTime.Add(timeoutCommit)
		consensusState.StartTime = consensusState.config.Commit(tmtime.Now())
	} else {
		consensusState.StartTime = consensusState.config.Commit(consensusState.CommitTime)
	}

	consensusState.Validators = validators
	consensusState.Proposal = nil
	consensusState.ProposalBlock = nil
	consensusState.ProposalBlockParts = nil
	consensusState.LockedRound = -1
	consensusState.LockedBlock = nil
	consensusState.LockedBlockParts = nil
	consensusState.ValidRound = -1
	consensusState.ValidBlock = nil
	consensusState.ValidBlockParts = nil
	consensusState.Votes = cstypes.NewHeightVoteSet(state.ChainID, height, validators)
	consensusState.CommitRound = -1
	consensusState.LastValidators = state.LastValidators
	consensusState.TriggeredTimeoutPrecommit = false
	consensusState.StandingMemberSet = state.StandingMemberSet
	consensusState.QrnSet = state.QrnSet
	consensusState.ConsensusRound = state.ConsensusRound

	fmt.Println("3stompesi-updateToState")
	consensusState.state = state

	// Finally, broadcast RoundState
	consensusState.newStep()
}

func (consensusState *State) newStep() {
	rs := consensusState.RoundStateEvent()
	if err := consensusState.wal.Write(rs); err != nil {
		consensusState.Logger.Error("failed writing to WAL", "err", err)
	}

	consensusState.nSteps++

	// newStep is called by updateToState in NewState before the eventBus is set!
	if consensusState.eventBus != nil {
		if err := consensusState.eventBus.PublishEventNewRoundStep(rs); err != nil {
			consensusState.Logger.Error("failed publishing new round step", "err", err)
		}

		consensusState.evsw.FireEvent(types.EventNewRoundStep, &consensusState.RoundState)
	}
}

//-----------------------------------------
// the main go routines

// receiveRoutine handles messages which may cause state transitions.
// it's argument (n) is the number of messages to process before exiting - use 0 to run forever
// It keeps the RoundState and is the only thing that updates it.
// Updates (state transitions) happen on timeouts, complete proposals, and 2/3 majorities.
// State must be locked before any internal state is updated.
func (consensusState *State) receiveRoutine(maxSteps int) {
	onExit := func(consensusState *State) {
		// NOTE: the internalMsgQueue may have signed messages from our
		// priv_val that haven't hit the WAL, but its ok because
		// priv_val tracks LastSig

		// close wal now that we're done writing to it
		if err := consensusState.wal.Stop(); err != nil {
			consensusState.Logger.Error("failed trying to stop WAL", "error", err)
		}

		consensusState.wal.Wait()
		close(consensusState.done)
	}

	defer func() {
		if r := recover(); r != nil {
			consensusState.Logger.Error("CONSENSUS FAILURE!!!", "err", r, "stack", string(debug.Stack()))
			// stop gracefully
			//
			// NOTE: We most probably shouldn't be running any further when there is
			// some unexpected panic. Some unknown error happened, and so we don't
			// know if that will result in the validator signing an invalid thing. It
			// might be worthwhile to explore a mechanism for manual resuming via
			// some console or secure RPC system, but for now, halting the chain upon
			// unexpected consensus bugs sounds like the better option.
			onExit(consensusState)
		}
	}()

	for {
		if maxSteps > 0 {
			if consensusState.nSteps >= maxSteps {
				consensusState.Logger.Debug("reached max steps; exiting receive routine")
				consensusState.nSteps = 0
				return
			}
		}

		roundState := consensusState.RoundState
		var mi MsgInfo

		select {
		case <-consensusState.txNotifier.TxsAvailable():
			consensusState.handleTxsAvailable()

		case mi = <-consensusState.peerMsgQueue:
			if err := consensusState.wal.Write(mi); err != nil {
				consensusState.Logger.Error("failed writing to WAL", "err", err)
			}

			// handles proposals, block parts, votes
			// may generate internal events (votes, complete proposals, 2/3 majorities)
			consensusState.handleMsg(mi)

		case mi = <-consensusState.internalMsgQueue:
			err := consensusState.wal.WriteSync(mi) // NOTE: fsync
			if err != nil {
				panic(fmt.Sprintf(
					"failed to write %v msg to consensus WAL due to %v; check your file system and restart the node",
					mi, err,
				))
			}

			if _, ok := mi.Msg.(*VoteMessage); ok {
				// we actually want to simulate failing during
				// the previous WriteSync, but this isn't easy to do.
				// Equivalent would be to fail here and manually remove
				// some bytes from the end of the wal.
				fail.Fail() // XXX
			}

			// handles proposals, block parts, votes
			consensusState.handleMsg(mi)

		case ti := <-consensusState.timeoutTicker.Chan(): // tockChan:
			if err := consensusState.wal.Write(ti); err != nil {
				consensusState.Logger.Error("failed writing to WAL", "err", err)
			}

			// if the timeout is relevant to the rs
			// go to the next step
			consensusState.handleTimeout(ti, roundState)

		case <-consensusState.Quit():
			onExit(consensusState)
			return
		}
	}
}

// state transitions on complete-proposal, 2/3-any, 2/3-one
func (consensusState *State) handleMsg(mi MsgInfo) {
	consensusState.mtx.Lock()
	defer consensusState.mtx.Unlock()

	var (
		added bool
		err   error
	)

	msg, peerID := mi.Msg, mi.PeerID

	switch msg := msg.(type) {
	case *ProposalMessage:
		// will not cause transition.
		// once proposal is set, we can receive block parts
		//fmt.Println("2stompesi-block-1")
		err = consensusState.setProposal(msg.Proposal)

	case *BlockPartMessage:
		// if the proposal is complete, we'll enterPrevote or tryFinalizeCommit
		//fmt.Println("2stompesi-block-2")
		added, err = consensusState.addProposalBlockPart(msg, peerID)
		if added {
			consensusState.statsMsgQueue <- mi
		}

		if err != nil && msg.Round != consensusState.Round {
			consensusState.Logger.Debug(
				"received block part from wrong round",
				"height", consensusState.Height,
				"cs_round", consensusState.Round,
				"block_round", msg.Round,
			)
			err = nil
		}

	case *QrnMessage:
		// attempt to add the vote and dupeout the validator if its a duplicate signature
		// if the vote gives us a 2/3-any or 2/3-one, we transition
		fmt.Println("handleMsg-QrnMessage")
		added, err = consensusState.addQrn(msg.Qrn)

	case *VoteMessage:
		// attempt to add the vote and dupeout the validator if its a duplicate signature
		// if the vote gives us a 2/3-any or 2/3-one, we transition
		//fmt.Println("2stompesi-block-3")
		added, err = consensusState.tryAddVote(msg.Vote, peerID)
		if added {
			consensusState.statsMsgQueue <- mi
		}

		// if err == ErrAddingVote {
		// TODO: punish peer
		// We probably don't want to stop the peer here. The vote does not
		// necessarily comes from a malicious peer but can be just broadcasted by
		// a typical peer.
		// https://github.com/reapchain/reapchain-core/issues/1281
		// }

		// NOTE: the vote is broadcast to peers by the reactor listening
		// for vote events

		// TODO: If rs.Height == vote.Height && rs.Round < vote.Round,
		// the peer is sending us CatchupCommit precommits.
		// We could make note of this and help filter in broadcastHasVoteMessage().

	default:
		consensusState.Logger.Error("unknown msg type", "type", fmt.Sprintf("%T", msg))
		return
	}

	if err != nil {
		consensusState.Logger.Error(
			"failed to process message",
			"height", consensusState.Height,
			"round", consensusState.Round,
			"peer", peerID,
			"err", err,
			"msg", msg,
		)
	}
}

func (consensusState *State) handleTimeout(ti timeoutInfo, rs cstypes.RoundState) {
	consensusState.Logger.Debug("received tock", "timeout", ti.Duration, "height", ti.Height, "round", ti.Round, "step", ti.Step)

	// timeouts must be for current height, round, step
	if ti.Height != rs.Height || ti.Round < rs.Round || (ti.Round == rs.Round && ti.Step < rs.Step) {
		consensusState.Logger.Debug("ignoring tock because we are ahead", "height", rs.Height, "round", rs.Round, "step", rs.Step)
		return
	}

	// the timeout will now cause a state transition
	consensusState.mtx.Lock()
	defer consensusState.mtx.Unlock()

	switch ti.Step {
	case cstypes.RoundStepNewHeight:
		// NewRound event fired from enterNewRound.
		// XXX: should we fire timeout here (for timeout commit)?
		consensusState.enterNewRound(ti.Height, 0)

	case cstypes.RoundStepNewRound:
		consensusState.enterPropose(ti.Height, 0)

	case cstypes.RoundStepPropose:
		if err := consensusState.eventBus.PublishEventTimeoutPropose(consensusState.RoundStateEvent()); err != nil {
			consensusState.Logger.Error("failed publishing timeout propose", "err", err)
		}

		consensusState.enterPrevote(ti.Height, ti.Round)

	case cstypes.RoundStepPrevoteWait:
		if err := consensusState.eventBus.PublishEventTimeoutWait(consensusState.RoundStateEvent()); err != nil {
			consensusState.Logger.Error("failed publishing timeout wait", "err", err)
		}

		consensusState.enterPrecommit(ti.Height, ti.Round)

	case cstypes.RoundStepPrecommitWait:
		if err := consensusState.eventBus.PublishEventTimeoutWait(consensusState.RoundStateEvent()); err != nil {
			consensusState.Logger.Error("failed publishing timeout wait", "err", err)
		}

		consensusState.enterPrecommit(ti.Height, ti.Round)
		consensusState.enterNewRound(ti.Height, ti.Round+1)

	default:
		panic(fmt.Sprintf("invalid timeout step: %v", ti.Step))
	}

}

func (consensusState *State) handleTxsAvailable() {
	consensusState.mtx.Lock()
	defer consensusState.mtx.Unlock()

	// We only need to do this for round 0.
	if consensusState.Round != 0 {
		return
	}

	switch consensusState.Step {
	case cstypes.RoundStepNewHeight: // timeoutCommit phase
		if consensusState.needProofBlock(consensusState.Height) {
			// enterPropose will be called by enterNewRound
			return
		}

		// +1ms to ensure RoundStepNewRound timeout always happens after RoundStepNewHeight
		timeoutCommit := consensusState.StartTime.Sub(tmtime.Now()) + 1*time.Millisecond
		consensusState.scheduleTimeout(timeoutCommit, consensusState.Height, 0, cstypes.RoundStepNewRound)

	case cstypes.RoundStepNewRound: // after timeoutCommit
		consensusState.enterPropose(consensusState.Height, 0)
	}
}

//-----------------------------------------------------------------------------
// State functions
// Used internally by handleTimeout and handleMsg to make state transitions

// Enter: `timeoutNewHeight` by startTime (commitTime+timeoutCommit),
// 	or, if SkipTimeoutCommit==true, after receiving all precommits from (height,round-1)
// Enter: `timeoutPrecommits` after any +2/3 precommits from (height,round-1)
// Enter: +2/3 precommits for nil at (height,round-1)
// Enter: +2/3 prevotes any or +2/3 precommits for block or any from (height, round)
// NOTE: consensusState.StartTime was already set for height.
func (consensusState *State) enterNewRound(height int64, round int32) {
	logger := consensusState.Logger.With("height", height, "round", round)

	if consensusState.Height != height || round < consensusState.Round || (consensusState.Round == round && consensusState.Step != cstypes.RoundStepNewHeight) {
		logger.Debug(
			"entering new round with invalid args",
			"current", fmt.Sprintf("%v/%v/%v", consensusState.Height, consensusState.Round, consensusState.Step),
		)
		return
	}

	if now := tmtime.Now(); consensusState.StartTime.After(now) {
		logger.Debug("need to set a buffer and log message here for sanity", "start_time", consensusState.StartTime, "now", now)
	}

	logger.Debug("entering new round", "current", fmt.Sprintf("%v/%v/%v", consensusState.Height, consensusState.Round, consensusState.Step))

	// increment validators if necessary
	validators := consensusState.Validators
	if consensusState.Round < round {
		validators = validators.Copy()
		validators.IncrementProposerPriority(tmmath.SafeSubInt32(round, consensusState.Round))
	}

	// Setup new round
	// we don't fire newStep for this step,
	// but we fire an event, so update the round step first
	consensusState.updateRoundStep(round, cstypes.RoundStepNewRound)
	consensusState.Validators = validators
	if round == 0 {
		// We've already reset these upon new height,
		// and meanwhile we might have received a proposal
		// for round 0.
	} else {
		logger.Debug("resetting proposal info")
		consensusState.Proposal = nil
		consensusState.ProposalBlock = nil
		consensusState.ProposalBlockParts = nil
	}

	consensusState.Votes.SetRound(tmmath.SafeAddInt32(round, 1)) // also track next round (round+1) to allow round-skipping
	consensusState.TriggeredTimeoutPrecommit = false

	if err := consensusState.eventBus.PublishEventNewRound(consensusState.NewRoundEvent()); err != nil {
		consensusState.Logger.Error("failed publishing new round", "err", err)
	}

	consensusState.metrics.Rounds.Set(float64(round))

	// Wait for txs to be available in the mempool
	// before we enterPropose in round 0. If the last block changed the app hash,
	// we may need an empty "proof" block, and enterPropose immediately.
	waitForTxs := consensusState.config.WaitForTxs() && round == 0 && !consensusState.needProofBlock(height)
	if waitForTxs {
		if consensusState.config.CreateEmptyBlocksInterval > 0 {
			consensusState.scheduleTimeout(consensusState.config.CreateEmptyBlocksInterval, height, round,
				cstypes.RoundStepNewRound)
		}
	} else {
		consensusState.enterPropose(height, round)
	}
}

// needProofBlock returns true on the first height (so the genesis app hash is signed right away)
// and where the last block (height-1) caused the app hash to change
func (consensusState *State) needProofBlock(height int64) bool {
	if height == consensusState.state.InitialHeight {
		return true
	}

	lastBlockMeta := consensusState.blockStore.LoadBlockMeta(height - 1)
	if lastBlockMeta == nil {
		panic(fmt.Sprintf("needProofBlock: last block meta for height %d not found", height-1))
	}

	return !bytes.Equal(consensusState.state.AppHash, lastBlockMeta.Header.AppHash)
}

// Enter (CreateEmptyBlocks): from enterNewRound(height,round)
// Enter (CreateEmptyBlocks, CreateEmptyBlocksInterval > 0 ):
// 		after enterNewRound(height,round), after timeout of CreateEmptyBlocksInterval
// Enter (!CreateEmptyBlocks) : after enterNewRound(height,round), once txs are in the mempool
func (consensusState *State) enterPropose(height int64, round int32) {
	logger := consensusState.Logger.With("height", height, "round", round)

	if consensusState.Height != height || round < consensusState.Round || (consensusState.Round == round && cstypes.RoundStepPropose <= consensusState.Step) {
		logger.Debug(
			"entering propose step with invalid args",
			"current", fmt.Sprintf("%v/%v/%v", consensusState.Height, consensusState.Round, consensusState.Step),
		)
		return
	}

	logger.Debug("entering propose step", "current", fmt.Sprintf("%v/%v/%v", consensusState.Height, consensusState.Round, consensusState.Step))

	defer func() {
		// Done enterPropose:
		consensusState.updateRoundStep(round, cstypes.RoundStepPropose)
		consensusState.newStep()

		// If we have the whole proposal + POL, then goto Prevote now.
		// else, we'll enterPrevote when the rest of the proposal is received (in AddProposalBlockPart),
		// or else after timeoutPropose
		if consensusState.isProposalComplete() {
			consensusState.enterPrevote(height, consensusState.Round)
		}
	}()

	// If we don't get the proposal and all block parts quick enough, enterPrevote
	consensusState.scheduleTimeout(consensusState.config.Propose(round), height, round, cstypes.RoundStepPropose)

	// Nothing more to do if we're not a validator
	if consensusState.privValidator == nil {
		logger.Debug("node is not a validator")
		return
	}

	logger.Debug("node is a validator")

	if consensusState.privValidatorPubKey == nil {
		// If this node is a validator & proposer in the current round, it will
		// miss the opportunity to create a block.
		logger.Error("propose step; empty priv validator public key", "err", errPubKeyIsNotSet)
		return
	}

	address := consensusState.privValidatorPubKey.Address()

	// if not a validator, we're done
	if !consensusState.Validators.HasAddress(address) {
		logger.Debug("node is not a validator", "addr", address, "vals", consensusState.Validators)
		return
	}

	if consensusState.isCoordiantor(address) {
		// logger.Debug("propose step; our turn to propose", "proposer", address)
		fmt.Println("propose step; our turn to propose", "proposer", address)
		consensusState.decideProposal(height, round)
	} else {
		logger.Debug("propose step; not our turn to propose", "proposer", consensusState.StandingMemberSet.GetCoordinator().Address)
	}
}

func (consensusState *State) isCoordiantor(address []byte) bool {
	// fmt.Println("isCoordinator", hex.EncodeToString(consensusState.state.QrnSet.Hash()))
	// fmt.Println("isCoordinator", hex.EncodeToString(address))
	// fmt.Println("isCoordinator", consensusState.StandingMemberSet.GetCoordinator().Address)

	return bytes.Equal(consensusState.StandingMemberSet.GetCoordinator().Address, address)
}

func (consensusState *State) defaultDecideProposal(height int64, round int32) {
	var block *types.Block
	var blockParts *types.PartSet

	// fmt.Println("jb")
	// time.Sleep(3000 * time.Millisecond)
	// fmt.Println("jb2")

	// Decide on block
	if consensusState.ValidBlock != nil {
		// If there is valid block, choose that.
		block, blockParts = consensusState.ValidBlock, consensusState.ValidBlockParts
	} else {
		// Create a new proposal block from state/txs from the mempool.
		fmt.Println("stompesi-createProposalBlock")

		block, blockParts = consensusState.createProposalBlock()
		if block == nil {
			return
		}
	}

	// Flush the WAL. Otherwise, we may not recompute the same proposal to sign,
	// and the privValidator will refuse to sign anything.
	if err := consensusState.wal.FlushAndSync(); err != nil {
		consensusState.Logger.Error("failed flushing WAL to disk")
	}

	// Make proposal
	propBlockID := types.BlockID{Hash: block.Hash(), PartSetHeader: blockParts.Header()}
	proposal := types.NewProposal(height, round, consensusState.ValidRound, propBlockID)
	p := proposal.ToProto()
	if err := consensusState.privValidator.SignProposal(consensusState.state.ChainID, p); err == nil {
		proposal.Signature = p.Signature

		// send proposal and block parts on internal msg queue
		//fmt.Println("stompesi-start-defaultDecideProposal")
		consensusState.sendInternalMessage(MsgInfo{&ProposalMessage{proposal}, ""})

		for i := 0; i < int(blockParts.Total()); i++ {
			part := blockParts.GetPart(i)
			consensusState.sendInternalMessage(MsgInfo{&BlockPartMessage{consensusState.Height, consensusState.Round, part}, ""})
		}

		consensusState.Logger.Debug("signed proposal", "height", height, "round", round, "proposal", proposal)
	} else if !consensusState.replayMode {
		consensusState.Logger.Error("propose step; failed signing proposal", "height", height, "round", round, "err", err)
	}
}

// Returns true if the proposal block is complete &&
// (if POLRound was proposed, we have +2/3 prevotes from there).
func (consensusState *State) isProposalComplete() bool {
	if consensusState.Proposal == nil || consensusState.ProposalBlock == nil {
		return false
	}
	// we have the proposal. if there's a POLRound,
	// make sure we have the prevotes from it too
	if consensusState.Proposal.POLRound < 0 {
		return true
	}
	// if this is false the proposer is lying or we haven't received the POL yet
	return consensusState.Votes.Prevotes(consensusState.Proposal.POLRound).HasTwoThirdsMajority()

}

// Create the next block to propose and return it. Returns nil block upon error.
//
// We really only need to return the parts, but the block is returned for
// convenience so we can log the proposal block.
//
// NOTE: keep it side-effect free for clarity.
// CONTRACT: consensusState.privValidator is not nil.
func (consensusState *State) createProposalBlock() (block *types.Block, blockParts *types.PartSet) {
	if consensusState.privValidator == nil {
		panic("entered createProposalBlock with privValidator being nil")
	}

	var commit *types.Commit
	switch {
	case consensusState.Height == consensusState.state.InitialHeight:
		// We're creating a proposal for the first block.
		// The commit is empty, but not nil.
		commit = types.NewCommit(0, 0, types.BlockID{}, nil)

	case consensusState.LastCommit.HasTwoThirdsMajority():
		// Make the commit from LastCommit
		commit = consensusState.LastCommit.MakeCommit()

	default: // This shouldn't happen.
		consensusState.Logger.Error("propose step; cannot propose anything without commit for the previous block")
		return
	}

	if consensusState.privValidatorPubKey == nil {
		// If this node is a validator & proposer in the current round, it will
		// miss the opportunity to create a block.
		consensusState.Logger.Error("propose step; empty priv validator public key", "err", errPubKeyIsNotSet)
		return
	}

	proposerAddr := consensusState.privValidatorPubKey.Address()
	fmt.Println("3stompesi-createProposalBlock", consensusState.state)
	return consensusState.blockExec.CreateProposalBlock(consensusState.Height, consensusState.state, commit, proposerAddr)
}

// Enter: `timeoutPropose` after entering Propose.
// Enter: proposal block and POL is ready.
// Prevote for LockedBlock if we're locked, or ProposalBlock if valid.
// Otherwise vote nil.
func (consensusState *State) enterPrevote(height int64, round int32) {
	logger := consensusState.Logger.With("height", height, "round", round)

	if consensusState.Height != height || round < consensusState.Round || (consensusState.Round == round && cstypes.RoundStepPrevote <= consensusState.Step) {
		logger.Debug(
			"entering prevote step with invalid args",
			"current", fmt.Sprintf("%v/%v/%v", consensusState.Height, consensusState.Round, consensusState.Step),
		)
		return
	}

	defer func() {
		// Done enterPrevote:
		consensusState.updateRoundStep(round, cstypes.RoundStepPrevote)
		consensusState.newStep()
	}()

	logger.Debug("entering prevote step", "current", fmt.Sprintf("%v/%v/%v", consensusState.Height, consensusState.Round, consensusState.Step))

	// Sign and broadcast vote as necessary
	consensusState.doPrevote(height, round)

	// Once `addVote` hits any +2/3 prevotes, we will go to PrevoteWait
	// (so we have more time to try and collect +2/3 prevotes for a single block)
}

func (consensusState *State) defaultDoPrevote(height int64, round int32) {
	logger := consensusState.Logger.With("height", height, "round", round)

	// If a block is locked, prevote that.
	if consensusState.LockedBlock != nil {
		logger.Debug("prevote step; already locked on a block; prevoting locked block")
		consensusState.signAddVote(tmproto.PrevoteType, consensusState.LockedBlock.Hash(), consensusState.LockedBlockParts.Header())
		return
	}

	// If ProposalBlock is nil, prevote nil.
	if consensusState.ProposalBlock == nil {
		logger.Debug("prevote step: ProposalBlock is nil")
		consensusState.signAddVote(tmproto.PrevoteType, nil, types.PartSetHeader{})
		return
	}

	// Validate proposal block
	err := consensusState.blockExec.ValidateBlock(consensusState.state, consensusState.ProposalBlock)
	if err != nil {
		// ProposalBlock is invalid, prevote nil.
		logger.Error("prevote step: ProposalBlock is invalid", "err", err)
		consensusState.signAddVote(tmproto.PrevoteType, nil, types.PartSetHeader{})
		return
	}

	// Prevote consensusState.ProposalBlock
	// NOTE: the proposal signature is validated when it is received,
	// and the proposal block parts are validated as they are received (against the merkle hash in the proposal)
	logger.Debug("prevote step: ProposalBlock is valid")
	consensusState.signAddVote(tmproto.PrevoteType, consensusState.ProposalBlock.Hash(), consensusState.ProposalBlockParts.Header())
}

// Enter: any +2/3 prevotes at next round.
func (consensusState *State) enterPrevoteWait(height int64, round int32) {
	logger := consensusState.Logger.With("height", height, "round", round)

	if consensusState.Height != height || round < consensusState.Round || (consensusState.Round == round && cstypes.RoundStepPrevoteWait <= consensusState.Step) {
		logger.Debug(
			"entering prevote wait step with invalid args",
			"current", fmt.Sprintf("%v/%v/%v", consensusState.Height, consensusState.Round, consensusState.Step),
		)
		return
	}

	if !consensusState.Votes.Prevotes(round).HasTwoThirdsAny() {
		panic(fmt.Sprintf(
			"entering prevote wait step (%v/%v), but prevotes does not have any +2/3 votes",
			height, round,
		))
	}

	logger.Debug("entering prevote wait step", "current", fmt.Sprintf("%v/%v/%v", consensusState.Height, consensusState.Round, consensusState.Step))

	defer func() {
		// Done enterPrevoteWait:
		consensusState.updateRoundStep(round, cstypes.RoundStepPrevoteWait)
		consensusState.newStep()
	}()

	// Wait for some more prevotes; enterPrecommit
	consensusState.scheduleTimeout(consensusState.config.Prevote(round), height, round, cstypes.RoundStepPrevoteWait)
}

// Enter: `timeoutPrevote` after any +2/3 prevotes.
// Enter: `timeoutPrecommit` after any +2/3 precommits.
// Enter: +2/3 precomits for block or nil.
// Lock & precommit the ProposalBlock if we have enough prevotes for it (a POL in this round)
// else, unlock an existing lock and precommit nil if +2/3 of prevotes were nil,
// else, precommit nil otherwise.
func (consensusState *State) enterPrecommit(height int64, round int32) {
	logger := consensusState.Logger.With("height", height, "round", round)

	if consensusState.Height != height || round < consensusState.Round || (consensusState.Round == round && cstypes.RoundStepPrecommit <= consensusState.Step) {
		logger.Debug(
			"entering precommit step with invalid args",
			"current", fmt.Sprintf("%v/%v/%v", consensusState.Height, consensusState.Round, consensusState.Step),
		)
		return
	}

	logger.Debug("entering precommit step", "current", fmt.Sprintf("%v/%v/%v", consensusState.Height, consensusState.Round, consensusState.Step))

	defer func() {
		// Done enterPrecommit:
		consensusState.updateRoundStep(round, cstypes.RoundStepPrecommit)
		consensusState.newStep()
	}()

	// check for a polka
	blockID, ok := consensusState.Votes.Prevotes(round).TwoThirdsMajority()

	// If we don't have a polka, we must precommit nil.
	if !ok {
		if consensusState.LockedBlock != nil {
			logger.Debug("precommit step; no +2/3 prevotes during enterPrecommit while we are locked; precommitting nil")
		} else {
			logger.Debug("precommit step; no +2/3 prevotes during enterPrecommit; precommitting nil")
		}

		consensusState.signAddVote(tmproto.PrecommitType, nil, types.PartSetHeader{})
		return
	}

	// At this point +2/3 prevoted for a particular block or nil.
	if err := consensusState.eventBus.PublishEventPolka(consensusState.RoundStateEvent()); err != nil {
		logger.Error("failed publishing polka", "err", err)
	}

	// the latest POLRound should be this round.
	polRound, _ := consensusState.Votes.POLInfo()
	if polRound < round {
		panic(fmt.Sprintf("this POLRound should be %v but got %v", round, polRound))
	}

	// +2/3 prevoted nil. Unlock and precommit nil.
	if len(blockID.Hash) == 0 {
		if consensusState.LockedBlock == nil {
			logger.Debug("precommit step; +2/3 prevoted for nil")
		} else {
			logger.Debug("precommit step; +2/3 prevoted for nil; unlocking")
			consensusState.LockedRound = -1
			consensusState.LockedBlock = nil
			consensusState.LockedBlockParts = nil

			if err := consensusState.eventBus.PublishEventUnlock(consensusState.RoundStateEvent()); err != nil {
				logger.Error("failed publishing event unlock", "err", err)
			}
		}

		consensusState.signAddVote(tmproto.PrecommitType, nil, types.PartSetHeader{})
		return
	}

	// At this point, +2/3 prevoted for a particular block.

	// If we're already locked on that block, precommit it, and update the LockedRound
	if consensusState.LockedBlock.HashesTo(blockID.Hash) {
		logger.Debug("precommit step; +2/3 prevoted locked block; relocking")
		consensusState.LockedRound = round

		if err := consensusState.eventBus.PublishEventRelock(consensusState.RoundStateEvent()); err != nil {
			logger.Error("failed publishing event relock", "err", err)
		}

		consensusState.signAddVote(tmproto.PrecommitType, blockID.Hash, blockID.PartSetHeader)
		return
	}

	// If +2/3 prevoted for proposal block, stage and precommit it
	if consensusState.ProposalBlock.HashesTo(blockID.Hash) {
		logger.Debug("precommit step; +2/3 prevoted proposal block; locking", "hash", blockID.Hash)

		// Validate the block.
		if err := consensusState.blockExec.ValidateBlock(consensusState.state, consensusState.ProposalBlock); err != nil {
			panic(fmt.Sprintf("precommit step; +2/3 prevoted for an invalid block: %v", err))
		}

		consensusState.LockedRound = round
		consensusState.LockedBlock = consensusState.ProposalBlock
		consensusState.LockedBlockParts = consensusState.ProposalBlockParts

		if err := consensusState.eventBus.PublishEventLock(consensusState.RoundStateEvent()); err != nil {
			logger.Error("failed publishing event lock", "err", err)
		}

		consensusState.signAddVote(tmproto.PrecommitType, blockID.Hash, blockID.PartSetHeader)
		return
	}

	// There was a polka in this round for a block we don't have.
	// Fetch that block, unlock, and precommit nil.
	// The +2/3 prevotes for this round is the POL for our unlock.
	logger.Debug("precommit step; +2/3 prevotes for a block we do not have; voting nil", "block_id", blockID)

	consensusState.LockedRound = -1
	consensusState.LockedBlock = nil
	consensusState.LockedBlockParts = nil

	if !consensusState.ProposalBlockParts.HasHeader(blockID.PartSetHeader) {
		consensusState.ProposalBlock = nil
		consensusState.ProposalBlockParts = types.NewPartSetFromHeader(blockID.PartSetHeader)
	}

	if err := consensusState.eventBus.PublishEventUnlock(consensusState.RoundStateEvent()); err != nil {
		logger.Error("failed publishing event unlock", "err", err)
	}

	consensusState.signAddVote(tmproto.PrecommitType, nil, types.PartSetHeader{})
}

// Enter: any +2/3 precommits for next round.
func (consensusState *State) enterPrecommitWait(height int64, round int32) {
	logger := consensusState.Logger.With("height", height, "round", round)

	if consensusState.Height != height || round < consensusState.Round || (consensusState.Round == round && consensusState.TriggeredTimeoutPrecommit) {
		logger.Debug(
			"entering precommit wait step with invalid args",
			"triggered_timeout", consensusState.TriggeredTimeoutPrecommit,
			"current", fmt.Sprintf("%v/%v", consensusState.Height, consensusState.Round),
		)
		return
	}

	if !consensusState.Votes.Precommits(round).HasTwoThirdsAny() {
		panic(fmt.Sprintf(
			"entering precommit wait step (%v/%v), but precommits does not have any +2/3 votes",
			height, round,
		))
	}

	logger.Debug("entering precommit wait step", "current", fmt.Sprintf("%v/%v/%v", consensusState.Height, consensusState.Round, consensusState.Step))

	defer func() {
		// Done enterPrecommitWait:
		consensusState.TriggeredTimeoutPrecommit = true
		consensusState.newStep()
	}()

	// wait for some more precommits; enterNewRound
	consensusState.scheduleTimeout(consensusState.config.Precommit(round), height, round, cstypes.RoundStepPrecommitWait)
}

// Enter: +2/3 precommits for block
func (consensusState *State) enterCommit(height int64, commitRound int32) {
	logger := consensusState.Logger.With("height", height, "commit_round", commitRound)

	if consensusState.Height != height || cstypes.RoundStepCommit <= consensusState.Step {
		logger.Debug(
			"entering commit step with invalid args",
			"current", fmt.Sprintf("%v/%v/%v", consensusState.Height, consensusState.Round, consensusState.Step),
		)
		return
	}

	logger.Debug("entering commit step", "current", fmt.Sprintf("%v/%v/%v", consensusState.Height, consensusState.Round, consensusState.Step))

	defer func() {
		// Done enterCommit:
		// keep consensusState.Round the same, commitRound points to the right Precommits set.
		consensusState.updateRoundStep(consensusState.Round, cstypes.RoundStepCommit)
		consensusState.CommitRound = commitRound
		consensusState.CommitTime = tmtime.Now()
		consensusState.newStep()

		// Maybe finalize immediately.
		consensusState.tryFinalizeCommit(height)
	}()

	blockID, ok := consensusState.Votes.Precommits(commitRound).TwoThirdsMajority()
	if !ok {
		panic("RunActionCommit() expects +2/3 precommits")
	}

	// The Locked* fields no longer matter.
	// Move them over to ProposalBlock if they match the commit hash,
	// otherwise they'll be cleared in updateToState.
	if consensusState.LockedBlock.HashesTo(blockID.Hash) {
		logger.Debug("commit is for a locked block; set ProposalBlock=LockedBlock", "block_hash", blockID.Hash)
		consensusState.ProposalBlock = consensusState.LockedBlock
		consensusState.ProposalBlockParts = consensusState.LockedBlockParts
	}

	// If we don't have the block being committed, set up to get it.
	if !consensusState.ProposalBlock.HashesTo(blockID.Hash) {
		if !consensusState.ProposalBlockParts.HasHeader(blockID.PartSetHeader) {
			logger.Info(
				"commit is for a block we do not know about; set ProposalBlock=nil",
				"proposal", consensusState.ProposalBlock.Hash(),
				"commit", blockID.Hash,
			)

			// We're getting the wrong block.
			// Set up ProposalBlockParts and keep waiting.
			consensusState.ProposalBlock = nil
			consensusState.ProposalBlockParts = types.NewPartSetFromHeader(blockID.PartSetHeader)

			if err := consensusState.eventBus.PublishEventValidBlock(consensusState.RoundStateEvent()); err != nil {
				logger.Error("failed publishing valid block", "err", err)
			}

			consensusState.evsw.FireEvent(types.EventValidBlock, &consensusState.RoundState)
		}
	}
}

// If we have the block AND +2/3 commits for it, finalize.
func (consensusState *State) tryFinalizeCommit(height int64) {
	logger := consensusState.Logger.With("height", height)

	if consensusState.Height != height {
		panic(fmt.Sprintf("tryFinalizeCommit() consensusState.Height: %v vs height: %v", consensusState.Height, height))
	}

	blockID, ok := consensusState.Votes.Precommits(consensusState.CommitRound).TwoThirdsMajority()
	if !ok || len(blockID.Hash) == 0 {
		logger.Error("failed attempt to finalize commit; there was no +2/3 majority or +2/3 was for nil")
		return
	}

	if !consensusState.ProposalBlock.HashesTo(blockID.Hash) {
		// TODO: this happens every time if we're not a validator (ugly logs)
		// TODO: ^^ wait, why does it matter that we're a validator?
		logger.Debug(
			"failed attempt to finalize commit; we do not have the commit block",
			"proposal_block", consensusState.ProposalBlock.Hash(),
			"commit_block", blockID.Hash,
		)
		return
	}

	consensusState.finalizeCommit(height)
}

// Increment height and goto cstypes.RoundStepNewHeight
func (consensusState *State) finalizeCommit(height int64) {
	logger := consensusState.Logger.With("height", height)

	if consensusState.Height != height || consensusState.Step != cstypes.RoundStepCommit {
		logger.Debug(
			"entering finalize commit step",
			"current", fmt.Sprintf("%v/%v/%v", consensusState.Height, consensusState.Round, consensusState.Step),
		)
		return
	}

	blockID, ok := consensusState.Votes.Precommits(consensusState.CommitRound).TwoThirdsMajority()
	block, blockParts := consensusState.ProposalBlock, consensusState.ProposalBlockParts

	if !ok {
		panic("cannot finalize commit; commit does not have 2/3 majority")
	}
	if !blockParts.HasHeader(blockID.PartSetHeader) {
		panic("expected ProposalBlockParts header to be commit header")
	}
	if !block.HashesTo(blockID.Hash) {
		panic("cannot finalize commit; proposal block does not hash to commit hash")
	}

	if err := consensusState.blockExec.ValidateBlock(consensusState.state, block); err != nil {
		panic(fmt.Errorf("+2/3 committed an invalid block: %w", err))
	}

	logger.Info(
		"finalizing commit of block",
		"hash", block.Hash(),
		"root", block.AppHash,
		"num_txs", len(block.Txs),
	)
	logger.Debug(fmt.Sprintf("%v", block))

	fail.Fail() // XXX

	// Save to blockStore.
	if consensusState.blockStore.Height() < block.Height {
		// NOTE: the seenCommit is local justification to commit this block,
		// but may differ from the LastCommit included in the next block
		precommits := consensusState.Votes.Precommits(consensusState.CommitRound)
		seenCommit := precommits.MakeCommit()
		consensusState.blockStore.SaveBlock(block, blockParts, seenCommit)
	} else {
		// Happens during replay if we already saved the block but didn't commit
		logger.Debug("calling finalizeCommit on already stored block", "height", block.Height)
	}

	fail.Fail() // XXX

	// Write EndHeightMessage{} for this height, implying that the blockstore
	// has saved the block.
	//
	// If we crash before writing this EndHeightMessage{}, we will recover by
	// running ApplyBlock during the ABCI handshake when we restart.  If we
	// didn't save the block to the blockstore before writing
	// EndHeightMessage{}, we'd have to change WAL replay -- currently it
	// complains about replaying for heights where an #ENDHEIGHT entry already
	// exists.
	//
	// Either way, the State should not be resumed until we
	// successfully call ApplyBlock (ie. later here, or in Handshake after
	// restart).
	endMsg := EndHeightMessage{height}
	if err := consensusState.wal.WriteSync(endMsg); err != nil { // NOTE: fsync
		panic(fmt.Sprintf(
			"failed to write %v msg to consensus WAL due to %v; check your file system and restart the node",
			endMsg, err,
		))
	}

	fail.Fail() // XXX

	// Create a copy of the state for staging and an event cache for txs.
	//fmt.Println("stompesi-jjjjjjkkkk2", consensusState.state.Qrns)
	stateCopy := consensusState.state.Copy()

	// Execute and commit the block, update and save the state, and update the mempool.
	// NOTE The block.AppHash wont reflect these txs until the next block.
	var (
		err          error
		retainHeight int64
	)

	//fmt.Println("stompesi-start-종빈", len(stateCopy.StandingMembers.StandingMembers))

	stateCopy, retainHeight, err = consensusState.blockExec.ApplyBlock(
		stateCopy,
		types.BlockID{
			Hash:          block.Hash(),
			PartSetHeader: blockParts.Header(),
		},
		block,
	)
	if err != nil {
		logger.Error("failed to apply block", "err", err)
		return
	}

	fail.Fail() // XXX

	// Prune old heights, if requested by ABCI app.
	if retainHeight > 0 {
		pruned, err := consensusState.pruneBlocks(retainHeight)
		if err != nil {
			logger.Error("failed to prune blocks", "retain_height", retainHeight, "err", err)
		} else {
			logger.Debug("pruned blocks", "pruned", pruned, "retain_height", retainHeight)
		}
	}

	// must be called before we update state
	consensusState.recordMetrics(height, block)

	// NewHeightStep!
	consensusState.updateToState(stateCopy)

	fail.Fail() // XXX

	// Private validator might have changed it's key pair => refetch pubkey.
	if err := consensusState.updatePrivValidatorPubKey(); err != nil {
		logger.Error("failed to get private validator pubkey", "err", err)
	}

	// consensusState.StartTime is already set.
	// Schedule Round0 to start soon.
	consensusState.scheduleRound0(&consensusState.RoundState)

	// By here,
	// * consensusState.Height has been increment to height+1
	// * consensusState.Step is now cstypes.RoundStepNewHeight
	// * consensusState.StartTime is set to when we will start round0.

	// genDoc.Qrns = []types.Qrn{{
	// 	Address: pubKey.Address(),
	// 	PubKey:  pubKey,
	// 	Value:   tmrand.Uint64(),
	// }}

	//TODO: stompesi

	if consensusState.state.ConsensusRound.ConsensusStartBlockHeight-1 == height {

		pubKey, err := consensusState.privValidator.GetPubKey()
		if err != nil {
			logger.Error("can't get pubkey", "err", err)
		}

		if index, _ := consensusState.state.StandingMemberSet.GetByAddress(pubKey.Address()); index != -1 {
			consensusState.signAddQrn()
		}
	}

	// consensusState.state.QrnSet.Qrns[0].Value = tmrand.Uint64()

	// Value:   tmrand.Uint64(),

}

func (consensusState *State) pruneBlocks(retainHeight int64) (uint64, error) {
	base := consensusState.blockStore.Base()
	if retainHeight <= base {
		return 0, nil
	}
	pruned, err := consensusState.blockStore.PruneBlocks(retainHeight)
	if err != nil {
		return 0, fmt.Errorf("failed to prune block store: %w", err)
	}
	err = consensusState.blockExec.Store().PruneStates(base, retainHeight)
	if err != nil {
		return 0, fmt.Errorf("failed to prune state database: %w", err)
	}
	return pruned, nil
}

func (consensusState *State) recordMetrics(height int64, block *types.Block) {
	//fmt.Println("stompesi-start-recordMetrics")

	consensusState.metrics.Validators.Set(float64(consensusState.Validators.Size()))
	consensusState.metrics.ValidatorsPower.Set(float64(consensusState.Validators.TotalVotingPower()))

	consensusState.metrics.StandingMembers.Set(float64(consensusState.StandingMemberSet.Size()))

	//fmt.Println("stompesi-start-recordMetrics-1")
	var (
		missingValidators      int
		missingValidatorsPower int64
	)
	// height=0 -> MissingValidators and MissingValidatorsPower are both 0.
	// Remember that the first LastCommit is intentionally empty, so it's not
	// fair to increment missing validators number.
	//fmt.Println("stompesi-start-recordMetrics-2")
	if height > consensusState.state.InitialHeight {
		//fmt.Println("stompesi-start-recordMetrics-3")
		// Sanity check that commit size matches validator set size - only applies
		// after first block.
		var (
			commitSize = block.LastCommit.Size()
			valSetLen  = len(consensusState.LastValidators.Validators)
			address    types.Address
		)
		if commitSize != valSetLen {
			panic(fmt.Sprintf("commit size (%d) doesn't match valset length (%d) at height %d\n\n%v\n\n%v",
				commitSize, valSetLen, block.Height, block.LastCommit.Signatures, consensusState.LastValidators.Validators))
		}

		if consensusState.privValidator != nil {
			if consensusState.privValidatorPubKey == nil {
				// Metrics won't be updated, but it's not critical.
				consensusState.Logger.Error(fmt.Sprintf("recordMetrics: %v", errPubKeyIsNotSet))
			} else {
				address = consensusState.privValidatorPubKey.Address()
			}
		}

		//fmt.Println("stompesi-start-recordMetrics-4")
		for i, val := range consensusState.LastValidators.Validators {
			commitSig := block.LastCommit.Signatures[i]
			if commitSig.Absent() {
				missingValidators++
				missingValidatorsPower += val.VotingPower
			}

			if bytes.Equal(val.Address, address) {
				label := []string{
					"validator_address", val.Address.String(),
				}
				consensusState.metrics.ValidatorPower.With(label...).Set(float64(val.VotingPower))
				if commitSig.ForBlock() {
					consensusState.metrics.ValidatorLastSignedHeight.With(label...).Set(float64(height))
				} else {
					consensusState.metrics.ValidatorMissedBlocks.With(label...).Add(float64(1))
				}
			}

		}
	}
	consensusState.metrics.MissingValidators.Set(float64(missingValidators))
	consensusState.metrics.MissingValidatorsPower.Set(float64(missingValidatorsPower))

	// NOTE: byzantine validators power and count is only for consensus evidence i.e. duplicate vote
	var (
		byzantineValidatorsPower = int64(0)
		byzantineValidatorsCount = int64(0)
	)
	for _, ev := range block.Evidence.Evidence {
		if dve, ok := ev.(*types.DuplicateVoteEvidence); ok {
			if _, val := consensusState.Validators.GetByAddress(dve.VoteA.ValidatorAddress); val != nil {
				byzantineValidatorsCount++
				byzantineValidatorsPower += val.VotingPower
			}
		}
	}
	consensusState.metrics.ByzantineValidators.Set(float64(byzantineValidatorsCount))
	consensusState.metrics.ByzantineValidatorsPower.Set(float64(byzantineValidatorsPower))

	if height > 1 {
		lastBlockMeta := consensusState.blockStore.LoadBlockMeta(height - 1)
		if lastBlockMeta != nil {
			consensusState.metrics.BlockIntervalSeconds.Observe(
				block.Time.Sub(lastBlockMeta.Header.Time).Seconds(),
			)
		}
	}

	//fmt.Println("stompesi-start-recordMetrics-5")
	consensusState.metrics.NumTxs.Set(float64(len(block.Data.Txs)))
	consensusState.metrics.TotalTxs.Add(float64(len(block.Data.Txs)))
	consensusState.metrics.BlockSizeBytes.Set(float64(block.Size()))
	consensusState.metrics.CommittedHeight.Set(float64(block.Height))
	//fmt.Println("stompesi-start-recordMetrics-6")
}

//-----------------------------------------------------------------------------

func (consensusState *State) defaultSetProposal(proposal *types.Proposal) error {
	// Already have one
	// TODO: possibly catch double proposals
	if consensusState.Proposal != nil {
		return nil
	}

	// Does not apply
	if proposal.Height != consensusState.Height || proposal.Round != consensusState.Round {
		return nil
	}

	// Verify POLRound, which must be -1 or in range [0, proposal.Round).
	if proposal.POLRound < -1 ||
		(proposal.POLRound >= 0 && proposal.POLRound >= proposal.Round) {
		return ErrInvalidProposalPOLRound
	}

	p := proposal.ToProto()
	// Verify signature
	if !consensusState.StandingMemberSet.GetCoordinator().PubKey.VerifySignature(
		types.ProposalSignBytes(consensusState.state.ChainID, p), proposal.Signature,
	) {
		return ErrInvalidProposalSignature
	}

	proposal.Signature = p.Signature
	consensusState.Proposal = proposal
	// We don't update consensusState.ProposalBlockParts if it is already set.
	// This happens if we're already in cstypes.RoundStepCommit or if there is a valid block in the current round.
	// TODO: We can check if Proposal is for a different block as this is a sign of misbehavior!
	if consensusState.ProposalBlockParts == nil {
		consensusState.ProposalBlockParts = types.NewPartSetFromHeader(proposal.BlockID.PartSetHeader)
	}

	consensusState.Logger.Info("received proposal", "proposal", proposal)
	return nil
}

// NOTE: block is not necessarily valid.
// Asynchronously triggers either enterPrevote (before we timeout of propose) or tryFinalizeCommit,
// once we have the full block.
func (consensusState *State) addProposalBlockPart(msg *BlockPartMessage, peerID p2p.ID) (added bool, err error) {
	height, round, part := msg.Height, msg.Round, msg.Part

	// Blocks might be reused, so round mismatch is OK
	if consensusState.Height != height {
		consensusState.Logger.Debug("received block part from wrong height", "height", height, "round", round)
		return false, nil
	}

	// We're not expecting a block part.
	if consensusState.ProposalBlockParts == nil {
		// NOTE: this can happen when we've gone to a higher round and
		// then receive parts from the previous round - not necessarily a bad peer.
		consensusState.Logger.Debug(
			"received a block part when we are not expecting any",
			"height", height,
			"round", round,
			"index", part.Index,
			"peer", peerID,
		)
		return false, nil
	}

	//fmt.Println("2stompesi-block-tt", part)
	added, err = consensusState.ProposalBlockParts.AddPart(part)
	if err != nil {
		return added, err
	}
	if consensusState.ProposalBlockParts.ByteSize() > consensusState.state.ConsensusParams.Block.MaxBytes {
		return added, fmt.Errorf("total size of proposal block parts exceeds maximum block bytes (%d > %d)",
			consensusState.ProposalBlockParts.ByteSize(), consensusState.state.ConsensusParams.Block.MaxBytes,
		)
	}

	if added && consensusState.ProposalBlockParts.IsComplete() {
		bz, err := ioutil.ReadAll(consensusState.ProposalBlockParts.GetReader())
		if err != nil {
			return added, err
		}

		var pbb = new(tmproto.Block)
		//fmt.Println("2stompesi-block-tt", pbb.Header.Height)
		err = proto.Unmarshal(bz, pbb)
		if err != nil {
			return added, err
		}
		//fmt.Println("2stompesi-block-1", pbb)

		block, err := types.BlockFromProto(pbb)
		//fmt.Println("2stompesi-block-2", block)

		if err != nil {
			return added, err
		}

		consensusState.ProposalBlock = block

		// NOTE: it's possible to receive complete proposal blocks for future rounds without having the proposal
		consensusState.Logger.Info("received complete proposal block", "height", consensusState.ProposalBlock.Height, "hash", consensusState.ProposalBlock.Hash())

		if err := consensusState.eventBus.PublishEventCompleteProposal(consensusState.CompleteProposalEvent()); err != nil {
			consensusState.Logger.Error("failed publishing event complete proposal", "err", err)
		}

		// Update Valid* if we can.
		prevotes := consensusState.Votes.Prevotes(consensusState.Round)
		blockID, hasTwoThirds := prevotes.TwoThirdsMajority()
		if hasTwoThirds && !blockID.IsZero() && (consensusState.ValidRound < consensusState.Round) {
			if consensusState.ProposalBlock.HashesTo(blockID.Hash) {
				consensusState.Logger.Debug(
					"updating valid block to new proposal block",
					"valid_round", consensusState.Round,
					"valid_block_hash", consensusState.ProposalBlock.Hash(),
				)

				consensusState.ValidRound = consensusState.Round
				consensusState.ValidBlock = consensusState.ProposalBlock
				consensusState.ValidBlockParts = consensusState.ProposalBlockParts
			}
			// TODO: In case there is +2/3 majority in Prevotes set for some
			// block and consensusState.ProposalBlock contains different block, either
			// proposer is faulty or voting power of faulty processes is more
			// than 1/3. We should trigger in the future accountability
			// procedure at this point.
		}

		if consensusState.Step <= cstypes.RoundStepPropose && consensusState.isProposalComplete() {
			// Move onto the next step
			consensusState.enterPrevote(height, consensusState.Round)
			if hasTwoThirds { // this is optimisation as this will be triggered when prevote is added
				consensusState.enterPrecommit(height, consensusState.Round)
			}
		} else if consensusState.Step == cstypes.RoundStepCommit {
			// If we're waiting on the proposal block...
			consensusState.tryFinalizeCommit(height)
		}

		return added, nil
	}

	return added, nil
}

// updatePrivValidatorPubKey get's the private validator public key and
// memoizes it. This func returns an error if the private validator is not
// responding or responds with an error.
func (consensusState *State) updatePrivValidatorPubKey() error {
	if consensusState.privValidator == nil {
		return nil
	}

	pubKey, err := consensusState.privValidator.GetPubKey()
	if err != nil {
		return err
	}
	consensusState.privValidatorPubKey = pubKey
	return nil
}

// look back to check existence of the node's consensus votes before joining consensus
func (consensusState *State) checkDoubleSigningRisk(height int64) error {
	if consensusState.privValidator != nil && consensusState.privValidatorPubKey != nil && consensusState.config.DoubleSignCheckHeight > 0 && height > 0 {
		valAddr := consensusState.privValidatorPubKey.Address()
		doubleSignCheckHeight := consensusState.config.DoubleSignCheckHeight
		if doubleSignCheckHeight > height {
			doubleSignCheckHeight = height
		}

		for i := int64(1); i < doubleSignCheckHeight; i++ {
			lastCommit := consensusState.blockStore.LoadSeenCommit(height - i)
			if lastCommit != nil {
				for sigIdx, s := range lastCommit.Signatures {
					if s.BlockIDFlag == types.BlockIDFlagCommit && bytes.Equal(s.ValidatorAddress, valAddr) {
						consensusState.Logger.Info("found signature from the same key", "sig", s, "idx", sigIdx, "height", height-i)
						return ErrSignatureFoundInPastBlocks
					}
				}
			}
		}
	}

	return nil
}

//---------------------------------------------------------

func CompareHRS(h1 int64, r1 int32, s1 cstypes.RoundStepType, h2 int64, r2 int32, s2 cstypes.RoundStepType) int {
	if h1 < h2 {
		return -1
	} else if h1 > h2 {
		return 1
	}
	if r1 < r2 {
		return -1
	} else if r1 > r2 {
		return 1
	}
	if s1 < s2 {
		return -1
	} else if s1 > s2 {
		return 1
	}
	return 0
}

// repairWalFile decodes messages from src (until the decoder errors) and
// writes them to dst.
func repairWalFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	var (
		dec = NewWALDecoder(in)
		enc = NewWALEncoder(out)
	)

	// best-case repair (until first error is encountered)
	for {
		msg, err := dec.Decode()
		if err != nil {
			break
		}

		err = enc.Encode(msg)
		if err != nil {
			return fmt.Errorf("failed to encode msg: %w", err)
		}
	}

	return nil
}
