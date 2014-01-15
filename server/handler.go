package server

import (
	"log"
)

type EchoMessageHandler struct {
	// nil
}

func NewEchoMessageHandler() *EchoMessageHandler {
	return &EchoMessageHandler{}
}

func (this *EchoMessageHandler) Handle(in []byte, s Session) []byte {
	log.Println(s.Name(), string(in))
	return in
}
