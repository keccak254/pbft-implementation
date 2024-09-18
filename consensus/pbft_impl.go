package consensus

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	consensus "github.com/keccak254/pbft-implementation.git/consensus/pbft_msg_types"
)

type State struct {
    ViewID         int64
    MsgLogs        *MsgLogs
    LastSequenceID int64
    CurrentStage   Stage
}

type MsgLogs struct {
    ReqMsg      *RequestMsg
    PrepareMsgs map[string]*VoteMsg
    CommitMsgs  map[string]*VoteMsg
}

type Stage int

const (
    Idle Stage = iota
    PrePrepared
    Prepared
    Committed
)

const f = 1

func CreateState(viewID int64, lastSequenceID int64) *State {
    return &State{
        ViewID: viewID,
        MsgLogs: &MsgLogs{
            ReqMsg:      nil,
            PrepareMsgs: make(map[string]*VoteMsg),
            CommitMsgs:  make(map[string]*VoteMsg),
        },
        LastSequenceID: lastSequenceID,
        CurrentStage:   Idle,
    }
}

func (state *State) StartConsensus(request *consensus.RequestMsg)(*consensus.PrePrepareMsg, error){
	sequenceID := time.Now().UnixNano()

	if state.LastSequenceID != -1 {
			for state.LastSequenceID >= sequenceID {
				sequenceID++
			}
	}

	request.SequenceID = sequenceID
	state.MsgLogs.ReqMsg = request

	digest, err := digest(request)
	if err != nil {
			return nil, err
	}

	state.CurrentStage = consensus.PrePrepared

	return &consensus.PrePrepareMsg{
		ViewID: state.ViewID,
		SequenceID: sequenceID,
		Digest: digest,
		RequestMsg: request,
	}, nil

}

func (state *State) PrePrepare(prePrepareMsg *PrePrepareMsg)(*VoteMsg, error) {

		state.MsgLogs.ReqMsg = prePrepareMsg.RequestMsg

		if !state.verifyMsg(prePrepareMsg.ViewID, prePrepareMsg.SequenceID, prePrepareMsg.Digest) {
			return nil, errors.New("pre-prepare message is corrupted")
		}

		state.CurrentStage = PrePrepared

		return &VoteMsg{
			ViewID:     state.ViewID,
			SequenceID: prePrepareMsg.SequenceID,
			Digest:     prePrepareMsg.Digest,
			MsgType:    PrepareMsg,
		}, nil
}