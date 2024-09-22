package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/keccak254/pbft-implementation.git/consensus"
)

type Node struct {
		NodeID string
		NodeTable map[string]string
		View *View
		CurrentState *consensus.State
		CommittedMsgs []*consensus.RequestMsg
		MsgBuffer *MsgBuffer
		MsgEntrance chan interface{}
		MsgDelivery chan interface{}
		Alarm chan bool
}

type MsgBuffer struct {
		ReqMsgs	[]*consensus.RequestMsg
		PrePrepareMsgs	[]*consensus.PrePrepareMsg
		PrepareMsgs		[]*consensus.VoteMsg
		CommitMsgs	[]*consensus.VoteMsg
}

type View struct {
	Id int64
	Primary string
}

const ResolvingTimeDuration = time.Millisecond * 1000

func NewNode(nodeID string) *Node {
		const viewID = 10000000000

		node := &Node{
				NodeID: nodeID,
				NodeTable: map[string]string{
					"Apple": "localhost:1111",
					"MS": "localhost:1112",
					"Google": "localhost:1113",
					"IBM": "localhost:1114",
				},
				View: &View{
					ID: viewID,
					Primary: "Apple",
				},

				CurrentState: nil,
				CommittedMsgs: make([]*consensus.RequestMsg, 0),
				MsgBuffer: &MsgBuffer{
						ReqMsgs: make([]*consensus.RequestMsg, 0),
						PrePrepareMsgs: make([]*consensus.PrePrepareMsg, 0),
						PrepareMsgs: make([]*consensus.VoteMsg, 0),
						CommitMsgs: make([]*consensus.VoteMsg, 0),
				},

				MsgEntrance: make(chan interface{}),
				MsgDelivery:  make(chan interface{}),
				Alarm: make(chan bool),
		}

    go node.dispatchMsg()
    go node.alarmToDispatcher()
    go node.resolveMsg()

    return node
}

func (node *Node) Broadcast(msg interface{}, path string) map[string]error {
		errorMap := make(map[string]error)

		for nodeID, url := range node.NodeTable {
				if nodeID == node.NodeID {
					continue
				}

				jsonMsg, err := json.Marshal(msg)
				if err != nil{
						errorMap[nodeID] = err
						continue
				}

				send(url + path, jsonMsg)
		}

		if len(errorMap) == 0 {
			return nil
		} else {
				return errorMap
		}
}

func (node *Node) Reply(msg *consensus.ReplyMsg) error {
	for _, value := range node.CommittedMsgs {
			fmt.Printf("Committed value: %s, %d, %s, %d\n",
					value.ClientID, value.Timestamp, value.Operation, value.SequenceID)
	}
	fmt.Println()

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
			return err
	}

	send(node.NodeTable[node.View.Primary] + "/reply", jsonMsg)

	return nil
}

func (node *Node) GetReq(reqMsg *consensus.RequestMsg) error {
	LogMsg(reqMsg)

	err := node.createStateForNewConsensus()
	if err != nil {
			return err
	}

	prePrepareMsg, err := node.CurrentState.StartConsensus(reqMsg)
	if err != nil {
			return err
	}

	LogStage(fmt.Sprintf("Consensus Process (ViewID:%d)", node.CurrentState.ViewID), false)

	if prePrepareMsg != nil {
			node.Broadcast(prePrepareMsg, "/preprepare")
			LogStage("Pre-prepare", true)
	}
	return nil
}

func (node *Node) GetPrePrepare(prePrepareMsg *consensus.PrePrepareMsg) error {
	LogMsg(prePrepareMsg)

	err := node.createStateForNewConsensus()
	if err != nil {
			return err
	}

	prepareMsg, err := node.CurrentState.PrePrepare(prePrepareMsg)
	if err != nil {
			return err
	}

	if prepareMsg != nil {
			prepareMsg.NodeID = node.NodeID
			LogStage("Pre-prepare", true)
			node.Broadcast(prepareMsg, "/prepare")
			LogStage("Prepare", false)
	}

	return nil
}

func (node *Node) GetPrepare(prepareMsg *consensus.VoteMsg) error {
		LogMsg(prepareMsg)

		commitMsg, err := node.CurrentState.Prepare(prepareMsg)
		if err != nil {
				return err
		}

		if commitMsg != nil {
			commitMsg.NodeID = node.NodeID
			LogStage("Prepare", true)
			node.Broadcast(commitMsg, "/commit")
			LogStage("Commit,", false)
		}
		return nil
}

func (node *Node) GetCommit(commitMsg *consensus.VoteMsg)error{
	LogMsg(commitMsg)

	replyMsg, committedMsg, err := node.CurrentState.Commit(commitMsg)
	if err != nil {
		return err
	}

	if replyMsg != nil{
			if committedMsg == nil {
				return error.New("committed massage is nil, even though the replay message is not nil")
			}

			replyMsg.NodeID = node.NodeID
			node.CommittedMsgs = append(node.CommittedMsgs, committedMsg)

			LogStage("Commit", true)
			node.Reply(replyMsg)
			LogStage("Reply", true)
	}
	return nil
}





