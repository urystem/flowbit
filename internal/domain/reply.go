package domain

type basicInputReply struct {
	ReplyToID uint64
	basicInputCommentReply
}

type ReplyForm struct {
	basicInputReply
	File *InPutObject // for s3
}

type InsertReply struct { // sql
	basicInputReply
	HasImage bool
}
