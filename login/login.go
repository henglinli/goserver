package login

import (
	"../message"
	"../server"
	"os"
	"log"
	"errors"
	"code.google.com/p/leveldb-go/leveldb"
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
	request  message.Request
	response message.Response
	db *leveldb.DB
}

// decode message
func (this *Login) decode(in []byte) error {
	err := proto.Unmarshal(in, &this.request)
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
func (this *Login) Handle(in []byte, s server.Session) []byte {
	log.Println("session:" , s.Name())
	// var
	var err error
	// decode message
	err = this.decode(in) 
	if err != nil {
		return []byte(err.Error())	
	}
	// message type check
	err = this.check(message.KLoginRequest)
	if err != nil {		
		*this.response.Status = message.Response_kError
		this.response.Error = proto.String(err.Error())
	} else {
	// command
		command := this.request.GetCommand()
		switch command {
		case message.Request_kPing:
			this.pong()
		case message.Request_kVeryfy:
			this.veryfy()
		case message.Request_kRegister:
			this.register()
		case message.Request_kLogin:
			this.login()
		case message.Request_kEnd:
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
	*this.response.Status = message.Response_kError
	this.response.Error = proto.String("Illegal or Bad command")
}

// pong
func (this *Login) pong() {
	// do nothing
}
// veryfy
func (this *Login) veryfy() {

}
// register
func (this *Login) register() {
	// has extension
	if proto.HasExtension(&this.request, message.E_Register_User) {
		data, err := proto.GetExtension(&this.request, 
			message.E_Register_User)
		if err != nil {
			user := data.(*message.User)
			account := user.GetAccount()
			profile := user.GetProfile()
			log.Println(account, profile)
		}
	}
	// not have extension
	*this.response.Status = message.Response_kError
	this.response.Error = proto.String("Request need extentsion 9")
}

// login
func (this *Login) login() {
	// 
}

// login manager
type LoginManager struct {
	logins chan Login
	db *leveldb.DB 
}

//
func NewLoginManager() *LoginManager {
	manager := &LoginManager{
		logins: make(chan Login, kMaxOnlineClients),
		db: nil,
	}
	//
	return manager
}

// open db
func (manager *LoginManager) OpenDB(name string) error {
	var err error
	manager.db, err = leveldb.Open(name, nil)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
// close db
func (manager *LoginManager) CloseDB() error {
	if manager.db != nil {
		return manager.db.Close()
	}
	return nil
}
//
func (manager *LoginManager) Handle(in []byte, s server.Session) (out []byte) {
	select {
	// get login
	case l := <-manager.logins:
		// handle
		out = l.Handle(in, s)
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
			request:  message.Request{},
			response: message.Response{},
			db : manager.db,
		}
		// set default response
		proto.SetDefaults(&l.response)
		// handle
		out = l.Handle(in, s)
		// recycle login
		select {
		case manager.logins <- l:
			// nil
		default:
			// nil
		}
	}
	//
	return out
}
