package consensus

type PBFT interface {
		StartConsensus(request *RequestMsg)(* PrePrepareMsg, error)
		PrePrepare(prePrepareMsg *PrePrepareMsg) (*VoteMsg, error)
		Prepare(prepareMsg *VoteMsg) (*VoteMsg, error)
		Commit(CommitMsg *VoteMsg) (*ReplyMsg, RequestMsg, error)
}

