package handler

import "marketflow/internal/ports/inbound"

type handler struct {
	use inbound.UsecaseInter
}

type HandleFuncs interface {
	mode
	health
	latest
}

func NewHandler(use inbound.UsecaseInter) HandleFuncs {
	return &handler{use: use}
}
