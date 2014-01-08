package server

import (
	"../message"
	"bufio"
	"io"
	"log"
	"net"
)

//
type MessageHandler interface {
	Handle([]byte) []byte
}

// 1, recv header
// 2, recv message
type Demuxer struct {
	// reader to writer channel
	forward chan []byte
	reader  *bufio.Reader
	writer  *bufio.Writer
	buffer  chan []byte
	handler MessageHandler
}

//
func (this *Demuxer) ValidSize(expected uint32) bool {
	if expected > 1024*1024*10 {
		return false
	}
	return true
}

// read message
func (this *Demuxer) read(buffer []byte) {
	//
	defer close(this.forward)
	//
	header := message.NewTinyHeader()
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
			header[:],
			len(header))

		if err != nil || readed != len(header) {
			log.Println("read header:", err.Error())
			break loop
		}
		// check maigc
		ok := header.CheckMagic()
		if ok != true {
			log.Println("*Inllegal client, magic:", header.GetMagic())
			break loop
		}
		// get size
		messageSize = header.GetSize()
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
		output := this.handler.Handle(input)
		this.forward <- output
	}
	log.Println("Demuxer.Read done")
}

// wirite message
func (this *Demuxer) write() {
	header := message.NewTinyHeader()
	var err error
	var writen int
	var size int
	//
loop:
	for data := range this.forward {
		size = len(data)
		header.SetSize(uint32(size))
		//
		writen, err = this.writer.Write(header[:])
		if err != nil || writen != len(header) {
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
	log.Println("Demuxer.Write done")
}

// no used now
func (this *Demuxer) listen() {
	var buffer []byte
	go this.read(buffer)
	go this.write()
}

//
func (this *Demuxer) Demux(buffer []byte) {
	go this.write()
	this.read(buffer)
}

//
func NewDemuxer(conn net.Conn, h MessageHandler) *Demuxer {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	demuxer := &Demuxer{
		forward: make(chan []byte),
		reader:  reader,
		writer:  writer,
		buffer:  make(chan []byte),
		handler: h,
	}
	//this.listen()

	return demuxer
}

const (
	kMaxOnlineClients = 2048
	kBufferSize       = 2048
)

// impl server.ConnectionHandler
type DemuxerHandler struct {
	buffers        chan []byte
	messageHandler MessageHandler
}

func (this *DemuxerHandler) Handle(conn net.Conn) {
	log.Println("DemuxerHandler.Handle", conn.RemoteAddr())
	//
	go func() {
		//
		defer conn.Close()
		// new demuxer
		demuxer := NewDemuxer(conn, this.messageHandler)
		// select buffer
		var buffer []byte
		select {
		case buffer = <-this.buffers:
			log.Println("reuse buffer", len(buffer))
			// nil
		default:
			buffer = make([]byte, kBufferSize)
			log.Println("make buffer")
		}
		//
		demuxer.Demux(buffer)
		// give back buffer back into free list
		this.buffers <- buffer
		//
		log.Println("DemuxerHandler.Handle done")
	}()
	log.Println("DemuxerHandler.Handle started")
}

func NewDemuxerHandler(h MessageHandler) *DemuxerHandler {
	handler := &DemuxerHandler{
		buffers:        make(chan []byte, kMaxOnlineClients*2),
		messageHandler: h,
	}

	return handler
}
