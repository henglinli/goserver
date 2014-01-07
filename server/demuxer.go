package server

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
)

//
type MessageHandler interface {
	Handle([]byte) []byte
}

//
func GetSize(b []byte) uint32 {
	return binary.BigEndian.Uint32(b[0:4])
}

//
func SetSize(s uint32, b []byte) {
	binary.BigEndian.PutUint32(b[0:4], s)
}

// 1, recv size
// 2, recv message
type Demuxer struct {
	forward        chan string
	reader         *bufio.Reader
	writer         *bufio.Writer
	buffer         chan []byte
	messageHandler MessageHandler
}

// read message
func (demuxer *Demuxer) read(buffer []byte) {
	//
	defer close(demuxer.forward)
	//
	bondary := []byte{0, 0, 0, 0}
	//
	var messageSize uint32
	var err error
	var readed int
	//
loop:
	for {
		log.Println("receiving...")
		// read boundary
		readed, err = io.ReadAtLeast(demuxer.reader, bondary, 4)
		if err != nil || readed < 4 {
			log.Println("read bondary:", err.Error())
			break loop
		}
		// get size
		messageSize = GetSize(bondary)
		log.Println("message size:", messageSize)
		// get buffer
		needed := len(buffer) - int(messageSize)
		if needed < 0 {
			log.Println("Growing buffer...")
			buffer = make([]byte, messageSize)
		}
		// get message
		readed, err = io.ReadAtLeast(demuxer.reader, buffer, int(messageSize))
		if err != nil || readed < int(messageSize) {
			log.Println("read message:", err.Error())
			break loop
		}
		// handle message
		input := buffer[0:messageSize]
		output := demuxer.messageHandler.Handle(input)
		demuxer.forward <- string(output)
	}
	log.Println("demuxer.Read done")
}

// wirite message
func (demuxer *Demuxer) write() {
	//
	bondary := []byte{0, 0, 0, 0}
	var err error
	var writen int
	var size int
	//
loop:
	for data := range demuxer.forward {
		log.Println("sending...")
		size = len(data)
		SetSize(uint32(size), bondary)
		//
		writen, err = demuxer.writer.Write(bondary)
		if err != nil || writen != 4 {
			log.Println(err.Error())
			break loop
		}
		//
		writen, err = demuxer.writer.WriteString(data)
		if err != nil || writen != len(data) {
			log.Println(err.Error())
			break loop
		}
		//
		err = demuxer.writer.Flush()
		if err != nil {
			log.Println(err.Error())
			break loop
		}
	}
	log.Println("demuxer.Write done")
}

// no used now
func (demuxer *Demuxer) listen() {
	var buffer []byte
	go demuxer.read(buffer)
	go demuxer.write()
}

//
func (demuxer *Demuxer) Demux(buffer []byte) {
	go demuxer.write()
	demuxer.read(buffer)
}

//
func NewDemuxer(conn net.Conn, handler MessageHandler) *Demuxer {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	demuxer := &Demuxer{
		forward:        make(chan string),
		reader:         reader,
		writer:         writer,
		buffer:         make(chan []byte),
		messageHandler: handler,
	}
	//demuxer.listen()

	return demuxer
}

const (
	kMaxOnlineClients = 2048
	kBufferSize       = 2048
)

type Buffer struct {
	buffer []byte
}

// impl server.ConnectionHandler
type DemuxerHandler struct {
	buffers        chan []byte
	messageHandler MessageHandler
}

func (handler *DemuxerHandler) Handle(conn net.Conn) {
	log.Println("DemuxerHandler.Handle", conn.RemoteAddr())
	//
	go func() {
		//
		defer conn.Close()
		// new demuxer
		demuxer := NewDemuxer(conn, handler.messageHandler)
		// select buffer
		var buffer []byte
		select {
		case buffer = <-handler.buffers:
			log.Println("reuse buffer", len(buffer))
			// nil
		default:
			buffer = make([]byte, kBufferSize)
			log.Println("make buffer")
		}
		//
		demuxer.Demux(buffer)
		// give back buffer back into free list
		handler.buffers <- buffer
		/*
			loop:   for {
					select {
					case recved, ok := <- demuxer.incoming:
						if ok == false {
							break loop
						}
						log.Println("redrect recved", len(recved))
						demuxer.outgoing <- recved
						// client.outgoing <- <- client.incoming
					}
				}
		*/
		log.Println("DemuxerHandler.Handle done")
	}()
	log.Println("DemuxerHandler.Handle started")
}

func NewDemuxerHandler(h MessageHandler) *DemuxerHandler {
	handler := &DemuxerHandler{
		buffers:        make(chan []byte, kMaxOnlineClients),
		messageHandler: h,
	}
	/*
		// buffer generator
		go func() {
			for {
				if handler.buffers.Len() == 0 {
					handler.buffers.PushFront(Buffer{buffer:})
				}
				front := handler.buffers.Front()
				select {
				case back := <-handler.put:
					handler.buffers.PushFront(back)
				case handler.get <- front.Value(Buffer):
					handler.buffers.Remove(front)
				}
			}
		}()
	*/
	//
	return handler
}
