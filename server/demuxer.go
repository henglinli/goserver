package server

import (
	"net"
	"log"
	"io"
	"bufio"
	"encoding/binary"
	"container/list"
)

type MessageHandler interface {
	Handle(string) string
}
// bondary 
type Bondary struct {
	message_size uint32 
	message_type uint32
}
//
func  GetBondary(b []byte) (s, t uint32) {
	s = binary.BigEndian.Uint32(b[0:4])
	t = binary.BigEndian.Uint32(b[4:8])
	return s,t
}
//
func  SetBondary(s, t uint32, b []byte) {
	binary.BigEndian.PutUint32(b[0:4], s)
	binary.BigEndian.PutUint32(b[4:8], t)
}
// 1, recv boundary 
// 2, recv message
type Demuxer struct {
	forward chan string
	reader *bufio.Reader
	writer *bufio.Writer
	buffer chan []byte
//	free_buffer chan []byte
	message_handler MessageHandler
}

//
func (demuxer *Demuxer) read(buffer []byte) {
	//	
	defer close(demuxer.forward)	
	//
	bondary := []byte{0,0,0,0,0,0,0,0}
	//
	var message_size uint32
	var message_type uint32
	var err error
	var readed int
	//
	loop:	for {
		log.Println("receiving...")
		// read boundary
		readed, err = io.ReadAtLeast(demuxer.reader, bondary, 8)
		if err != nil || readed < 8 {
			log.Println("read bondary:", err.Error())
			break loop			
		}
		// get size and type
		message_size, message_type = GetBondary(bondary)
		log.Println("message size:", message_size)
		log.Println("message type:", message_type)
		// get buffer
		needed := len(buffer) - int(message_size)
		if needed < 0 {
			log.Println("Growing buffer...")
			buffer = make([]byte, message_size)
		}
		// get message
		readed, err = io.ReadAtLeast(demuxer.reader, buffer, int(message_size))
		if err != nil || readed < int(message_size) {
			log.Println("read message:", err.Error())
			break loop		
		}
		// handle message
		input := string(buffer[0:message_size])
		output := demuxer.message_handler.Handle(input)
		demuxer.forward <- output
	}
	log.Println("demuxer.Read done")
}

//
func (demuxer *Demuxer) write() {
	//
	bondary := []byte{0,0,0,0,0,0,0,0}
	var err error
	var writen int
	var size int
	//
	loop: for data := range demuxer.forward {
		log.Println("sending...")
		size = len(data)
		SetBondary(uint32(size), 0, bondary)
		//
		writen, err = demuxer.writer.Write(bondary)
		if err != nil || writen != 8 {
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
//
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
		forward: make(chan string),
		reader:   reader,
		writer:   writer,
		buffer : make(chan []byte),
	//	free_buffer: make(chan []byte),
		message_handler: handler,
	}

	//demuxer.listen()

	return demuxer
}

const (
	kMaxOnlineClients = 2048
	kBufferSize = 2048
)

type Buffer struct {
	buffer []byte
}

// impl server.ConnectionHandler
type DemuxerHandler struct {
	get, put chan []byte
	buffers chan [2048]byte
	message_handler MessageHandler
}

func(handler *DemuxerHandler) Handle(conn net.Conn) {
	log.Println("DemuxerHandler.Handle", conn.RemoteAddr())
	//
	go func() {
		//
		defer conn.Close()
		//
		demuxer := NewDemuxer(conn, handler.message_handler)	
	        // select buffer
		var buffer []byte
		select {
		case buffer := <- handler.buffers:
			// nil
		default:
			buffer = make([]byte, kBufferSize) 

		}
		demuxer.Demux(buffer)
		// give back buffer back into free list
		handler.put <- buffer
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

func NewDemuxerHandler(h MessageHandler) *DemuxerHandler{
	handler := &DemuxerHandler{
		get: make(chan []byte),
		put: make(chan []byte),
		buffers: make(chan [kBufferSize], kMaxOnlineClients),
		message_handler: h,
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
