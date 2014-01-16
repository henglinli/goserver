package client

import (
	"../message"
	"../server"
	"code.google.com/p/goprotobuf/proto"
	//"code.google.com/p/leveldb-go/leveldb"
	"errors"
	"log"
	"os"
)

// impl server.MessageHandler
type ProtobufClient struct {
	logger   *log.Logger
	request  *message.Request
	response message.Response
	session  server.Session
}

// check message type
func (this *ProtobufClient) CheckType(desired, expected uint32) error {
	if expected != desired {
		if expected > message.KNil && expected < message.KEnd {
			return errors.New("Bad message type")
		} else {
			return errors.New("Illegal message type")
		}
	}
	return nil
}

// handle
func (this *ProtobufClient) Handle(in interface{}, s server.Session) interface{} {
	this.logger.Println("session:", s.Name())
	// get request
	var ok bool
	this.request, ok = in.(*message.Request)
	if ok != true {
		*this.response.Status = message.Response_kError
		this.response.Error = proto.String("Bad message.Request type")
		return &this.response
	}
	// check type
	if false {
		err := this.CheckType(message.KLoginRequest,
			this.request.GetType())
		//
		if err != nil {
			*this.response.Status = message.Response_kError
			this.response.Error = proto.String(err.Error())
			return &this.response
		}
	}
	// check command
	command := this.request.GetCommand()

	this.logger.Println(command)
	*this.response.Status = message.Response_kError
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
	//
	return &this.response
}

// bad command
func (this *ProtobufClient) badCommand() {
	this.response.Error = proto.String("Illegal or Bad command")
}

// pong
func (this *ProtobufClient) pong() {
	// do nothing
	*this.response.Status = message.Response_kOk
}

// veryfy
func (this *ProtobufClient) veryfy() {

}

// register
func (this *ProtobufClient) register() {

	// has extension
	if proto.HasExtension(this.request, message.E_Register_User) {
		data, err := proto.GetExtension(this.request,
			message.E_Register_User)
		if err != nil {
			this.response.Error = proto.String(err.Error())
			return
		}
		user, ok := data.(*message.User)
		if ok != true {
			*this.response.Status = message.Response_kError
			this.response.Error =
				proto.String("Bad Register.User type")
			return
		}
		account := user.GetAccount()
		name := account.GetName()
		profile := user.GetProfile()
		log.Println("Register", name, account, profile)
		if false {
			user_byte, err := proto.Marshal(user)
			if err != nil {
				this.response.Error = proto.String(err.Error())
				return
			}
			//err = this.session.db.Set([]byte(name), user_byte)
			if err != nil {
				this.response.Error = proto.String(err.Error())
				return
			}
			log.Println(user_byte)
		}

		//this.session.user = user;
		//
		captcha := &message.Captcha{}
		captcha.Code = proto.String("captcha")
		captcha.Image = []byte("captcha")
		*this.response.Status = message.Response_kOk
		err = proto.SetExtension(&this.response,
			message.E_Register_Captcha,
			captcha)
		return
	}
	// not have extension
	this.response.Error = proto.String("Request need extentsion 9")
}

// login
func (this *ProtobufClient) login() {
	//
}

// login manager
type ProtobufClientManager struct {
	logger  *log.Logger
	clients chan *ProtobufClient
	coder   message.Coder
}

//
func NewProtobufClientManager() *ProtobufClientManager {
	manager := &ProtobufClientManager{
		logger:  nil,
		clients: make(chan *ProtobufClient, server.KMaxOnlineClients),
		coder:   message.NewProtobufCoder(),
	}
	manager.logger = log.New(os.Stdout,
		"",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)
	//
	return manager
}

//
func (this *ProtobufClientManager) GetClient() *ProtobufClient {
	select {
	case c := <-this.clients:
		return c
	default:
		c := &ProtobufClient{
			logger:   this.logger,
			request:  nil,
			response: message.Response{},
			session:  nil,
		}
		// set default response
		proto.SetDefaults(&c.response)
		return c
	}
}

//
func (this *ProtobufClientManager) PutClient(in *ProtobufClient) {
	select {
	case this.clients <- in:
		// nil
	default:
		// nil
	}
}

//
func (this *ProtobufClientManager) Handle(in []byte, s server.Session) []byte {
	this.logger.Println("decode request ...")
	// decode
	request, err := this.coder.Decode(in)
	if err != nil {
		return []byte(err.Error())
	}
	// get
	c := this.GetClient()
	// handle
	this.logger.Println("hanle request ...")
	out := c.Handle(request, s)
	// recycle login
	this.PutClient(c)
	// encode
	this.logger.Println("encode reponse ...")
	response, err := this.coder.Encode(out)
	if err != nil {
		return []byte(err.Error())
	} else {
		return response
	}
}
