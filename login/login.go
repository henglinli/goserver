package login

import (
	"log"
	//	"../server"
	"code.google.com/p/goprotobuf/proto"
)

const (
	kMaxOnlineClients = 2048
)

//
type MessageHandler interface {
	Handle([]byte) []byte
}

// impl server.MessageHandler
type Login struct {
	request  Request
	response Response
}

//
func (l *Login) Handle(input []byte) []byte {
	var err error
	var message []byte
	//
	log.Println("unmarshal" , len(input))
	err = proto.Unmarshal(input, &l.request)
	log.Println("unmarshal done")
	if err != nil {
		log.Println(err.Error())
		l.bad()
	} else {
		//
		command := l.request.GetCommand()
		switch command {
		case Request_kPing:
			l.pong()
		case Request_kRegister:
			l.register()
		case Request_kLogin:
			l.login()
		case Request_kEnd:
			fallthrough
		default:
			l.bad()
		}
	}
	//
	message, err = proto.Marshal(&l.response)
	if err != nil {
		log.Println(err.Error)
		return []byte{0}
	}
	return message
}

// bad
func (l *Login) bad() {
	*l.response.Status = Response_kOthers
	l.response.Error = proto.String("BadMessage")
}

// pong
func (l *Login) pong() {
	*l.response.Status = Response_kOk
}

// register
func (l *Login) register() {
	*l.response.Status = Response_kOk
	l.response.User = l.request.GetRegister()
}

// login
func (l *Login) login() {
	*l.response.Status = Response_kOk
}

// login manager
type LoginManager struct {
	logins chan Login
}

//
func NewLoginManager() *LoginManager {
	manager := &LoginManager{
		logins: make(chan Login, kMaxOnlineClients),
	}
	//
	return manager
}

//
func (manager *LoginManager) Handle(input []byte) (output []byte) {
	select {
	// get login
	case l := <-manager.logins:
		// handle
		output = l.Handle(input)
		// recycle login
		select {
		case manager.logins <- l:
			// nil
		default:
			// nil
		}
	default:
		// make login
		l := Login{
			request:  Request{},
			response: Response{},
		}
		// handle
		output = l.Handle(input)
		// recycle login
		select {
		case manager.logins <- l:
			// nil
		default:
			// nil
		}
	}
	//
	return output
}
