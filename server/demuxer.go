package server

import (
	"net"
	"log"
	"io"
	"bufio"
	"bytes"
	"encoding/binary"
)

type MessageHandler interface {
	Read()
	Write()
}
// bondary 
type Bondary struct {
	message_size uint32 
	message_type uint32
}

// 1, recv boundary 
// 2, recv message
type Demuxer struct {
	incoming chan string
	outgoing chan string
	reader *bufio.Reader
	writer *bufio.Writer
	buffer chan *bytes.Buffer
	message_handler MessageHandler
}

//
func (demuxer *Demuxer) Read() {
	boundary := []byte{0,0,0,0,0,0,0,0}
	size_buf := bytes.NewReader(boundary[0:4])
	type_buf := bytes.NewReader(boundary[4:8])
	buffer := <- demuxer.buffer
	var message_size uint32
	var message_type uint32
	var err error
	var readed int
	for {
		// read boundary
		readed, err = io.ReadAtLeast(demuxer.reader, boundary, 8)
		if err != nil || readed < 8 {
			log.Println(err.Error())
			close(demuxer.incoming)
			close(demuxer.outgoing)
			break			
		}
		// get size
		err = binary.Read(size_buf, binary.BigEndian, &message_size)
		if err != nil {
			log.Println(err.Error())
		}
		log.Println("message size: ", message_size)
		// get type
		err = binary.Read(type_buf, binary.BigEndian, &message_type)
		if err != nil {
			log.Println(err.Error())
		}
		log.Println("message type: ", message_type)
		// get buffer
		length := buffer.Len()
		needed := length - int(message_size)
		if needed > 0 {
			buffer.Grow(needed)
		}
		// get message
		readed, err = io.ReadAtLeast(demuxer.reader, buffer.Bytes(),int(message_size))
		if err != nil || readed < int(message_size) {
			log.Println(err.Error())
			close(demuxer.incoming)
			close(demuxer.outgoing)
			break			
		}
		demuxer.message_handler.Read()
		// demuxer.incoming <- message.String()
	}
	log.Println("reader.Read done")
}

//
func (demuxer *Demuxer) Write() {
	for data := range demuxer.outgoing {
		demuxer.writer.WriteString(data)
		demuxer.writer.Flush()
	}
	log.Println("writer.Write done")
}

//
func (demuxer *Demuxer) Listen() {
	go demuxer.Read()
	go demuxer.Write()
}

//
func NewDemuxer(conn net.Conn, handler MessageHandler) *Demuxer {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	demuxer := &Demuxer{
		incoming: make(chan string),
		outgoing: make(chan string),
		reader:   reader,
		writer:   writer,
		buffer : make(chan *bytes.Buffer),
		message_handler: handler,
	}

	demuxer.Listen()

	return demuxer
}

const (
	kMaxOnlineClients = 1024
)
// impl server.ConnectionHandler
type DemuxerHandler struct {
	buffers chan *bytes.Buffer
	message_handler MessageHandler
}

func(handler *DemuxerHandler) Handle(conn net.Conn) {
	log.Println("DemuxerHandler.Handle ", conn.RemoteAddr())
	go func() {
		defer conn.Close()
		var buffer *bytes.Buffer
		// select buffer
		select {
		case buffer = <- handler.buffers:
		default:
			buffer = new(bytes.Buffer)
		}
		//
		client := NewDemuxer(conn, handler.message_handler)
		client.buffer <- buffer
	loop:   for {
			log.Println("recv...")
			select {
			case recved, ok := <- client.incoming:
				log.Println(ok, " recv: ", len(recved))
				if ok == false {
					break loop
				}				
				client.outgoing <- recved
				// client.outgoing <- <- client.incoming
			}
		}
		log.Println("DemuxerHandler.Handle done")
		handler.buffers <- buffer
	}()
	log.Println(" EchoHandler.Serve started")
}

func NewDemuxerHandler(h MessageHandler) *DemuxerHandler{
	handler := &DemuxerHandler{
		buffers : make(chan *bytes.Buffer, kMaxOnlineClients),
		message_handler: h,
	}
	
	return handler
}
