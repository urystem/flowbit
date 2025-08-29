package domain

import (
	"time"

	"github.com/google/uuid"
)

// for catalog
type PostNonContent struct {
	ID       uint64
	Title    string
	HasImage bool
}

// create-post(like template) basic
type basicInputPost struct {
	Uuid    uuid.UUID
	Name    string
	Subject string
	Content string
}

// create-post db
type InsertPost struct {
	basicInputPost
	HasImage bool
}

// create-post (general)
type Form struct {
	basicInputPost
	File *InPutObject
}

// output
type Post struct {
	ID uint64
	PostX
}

type PostX struct {
	UserName string
	Title    string
	Content  string
	HasImage bool
	DataTime time.Time
}

// archived post
type ArchivePost struct {
	Post     Post
	Comments []Comment
}

type ActivePost struct {
	Post         Post
	CommentTries []*CommentTree
}
