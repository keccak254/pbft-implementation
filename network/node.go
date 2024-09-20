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
