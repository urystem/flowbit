package handler

import (
	"marketflow/internal/ports/inbound"
)

type handler struct {
	use inbound.UsecaseInter
}

func NewHandler(use inbound.UsecaseInter) HandleFuncs {
	return &handler{use: use}
}

type HandleFuncs interface {
	mode
}
