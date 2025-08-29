package domain

import (
	"io"
	"time"
)

type InPutObject struct {
	io.Reader
	ObjName string
	ConType string
	Size    int64
}

type OutputObject struct {
	io.ReadSeekCloser
	Modified time.Time
	ConType  string
}
