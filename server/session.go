package server

import (
	"bufio"
	"code.google.com/p/leveldb-go/leveldb"
	"io"
	"log"
	"net"
)

// impl Session
type AddrSession struct {
	address string
	manager SessionManager
	db      *leveldb.DB
	forward chan []byte
	reader  *bufio.Reader
	writer  *bufio.Writer
	handler MessageHandler
	rheader *TinyHeader
	wheader *TinyHeader
	user    interface{}
	captcha interface{}
}

//
func NewAddrSession(conn net.Conn,
	db *leveldb.DB,
	m SessionManager,
	h MessageHandler) *AddrSession {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	//
	addr := conn.RemoteAddr().String()
	//
	session := &AddrSession{
		address: addr,
		manager: m,
		db:      db,
		forward: make(chan []byte),
		reader:  reader,
		writer:  writer,
		rheader: NewTinyHeader(),
		wheader: NewTinyHeader(),
		handler: h,
		user:    nil,
		captcha: nil,
	}

	return session
}

//
func (this *AddrSession) IsLogin() bool {
	return this.manager.IsLogin(this.address)
}

//
func (this *AddrSession) Name() string {
	return this.address
}

//
func (this *AddrSession) GetDB() *leveldb.DB {
	return this.db
}

func (this *AddrSession) Handle(buffer []byte) {
	go this.write()
	this.read(buffer)
}

//
func (this *AddrSession) ValidSize(expected uint32) bool {
	if expected > 1024*1024*10 {
		return false
	}
	return true
}

// read message
func (this *AddrSession) read(buffer []byte) {
	//
	defer close(this.forward)
	//
	var messageSize uint32
	var err error
	var readed int
	//
loop:
	for {
		log.Println("receiving...")
		// read header
		readed, err = io.ReadAtLeast(this.reader,
			this.rheader[:],
			len(this.rheader))

		if err != nil || readed != len(this.rheader) {
			log.Println("read header:", err.Error())
			break loop
		}
		/*
			// check maigc
			ok := rheader.CheckMagic()
			if ok != true {
				break loop
			}
		*/
		// get size
		messageSize = this.rheader.GetSize()
		// check size
		if this.ValidSize(messageSize) == false {
			log.Println("Message size too big:", messageSize)
			break loop
		}
		// get buffer
		needed := len(buffer) - int(messageSize)
		if needed < 0 {
			log.Println("Growing buffer...")
			buffer = make([]byte, messageSize)
		}
		// get message
		readed, err = io.ReadAtLeast(this.reader, buffer, int(messageSize))
		if err != nil || readed < int(messageSize) {
			log.Println("read message:", err.Error())
			break loop
		}
		// handle message
		input := buffer[0:messageSize]
		this.forward <- this.handler.Handle(input, this)
	}
	log.Println("Session.Read done")
}

// wirite message
func (this *AddrSession) write() {
	var err error
	var writen int
	var size int
	//
loop:
	for data := range this.forward {
		size = len(data)
		this.wheader.SetSize(uint32(size))
		//
		writen, err = this.writer.Write(this.wheader[:])
		if err != nil || writen != len(this.wheader) {
			log.Println(err.Error())
			break loop
		}
		//
		writen, err = this.writer.Write(data)
		if err != nil || writen != len(data) {
			log.Println(err.Error())
			break loop
		}
		//
		err = this.writer.Flush()
		if err != nil {
			log.Println(err.Error())
			break loop
		}
	}
	log.Println("Session.Write done")
}
