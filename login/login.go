package login

import (
	"os"
	"log"
	"errors"
	"../message"
	//	"../server"
	"code.google.com/p/goprotobuf/proto"
)

const (
	kMaxOnlineClients = 2048
	kBufferSize = 2048
)

//
type MessageHandler interface {
	Handle([]byte) []byte
}

// impl server.MessageHandler
type Login struct {
	logger *log.Logger
	request  Request
	response Response
}

// decode message
func (this *Login) decode(input []byte) error {
	err := proto.Unmarshal(input, &this.request)
	if err != nil {
		this.logger.Println(err.Error())
		return errors.New("Illegal protocol")
	}
	return nil
}

// check message type 
func (this *Login) check(desired uint32) error {
	expected := this.request.GetType()
	if expected != desired {
		if expected > message.KNil && expected < message.KEnd {
			return errors.New("Bad message type")
		} else {
			return errors.New("Illegal message type")
		}
	}
	return nil
}

// encode message
func (this *Login) encode() []byte {
	message, err := proto.Marshal(&this.response)
	if err != nil {
		this.logger.Println(err.Error)
		return []byte("Internal Error: proto.Marshal")
	}
	return message
}

// handle
func (this *Login) Handle(input []byte) []byte {
	// var
	var err error
	// decode message
	err = this.decode(input) 
	if err != nil {
		return []byte(err.Error())	
	}
	// message type check
	err = this.check(message.KLoginRequest)
	if err != nil {		
		*this.response.Status = Response_kError
		this.response.Error = proto.String(err.Error())
	} else {
	// command
		command := this.request.GetCommand()
		switch command {
		case Request_kPing:
			this.pong()
		case Request_kRegister:
			this.register()
		case Request_kLogin:
			this.login()
		case Request_kEnd:
			fallthrough
		default:
			this.badCommand()
		}
	}		
	// encode message
	return this.encode()
}

// bad command
func (this *Login) badCommand() {
	*this.response.Status = Response_kError
	this.response.Error = proto.String("Illegal or Bad command")
}

// pong
func (this *Login) pong() {
	this.response.Pong = proto.String("Pong")
}

// register
func (this *Login) register() {
	this.response.User = this.request.GetRegister()
}

// login
func (this *Login) login() {
	// 
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
			logger: log.New(os.Stdout, 
				"", 
				log.Ldate|log.Lmicroseconds|log.Lshortfile),
			request:  Request{},
			response: Response{},
		}
		// set default response
		proto.SetDefaults(&l.response)
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
