package main

import (
	"../message"
	"../server"
	"bufio"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

//
func GetSize(b []byte) uint32 {
	return binary.BigEndian.Uint32(b[0:4])
}

//
func SetSize(s uint32, b []byte) {
	binary.BigEndian.PutUint32(b[0:4], s)
}

//
type Client struct {
	rheader  *server.TinyHeader
	wheader  *server.TinyHeader
	incoming chan string
	outgoing chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func (this *Client) Connect(address string) error {
	log.Println("connecting to ", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
		log.Println(err.Error())
	}
	log.Println("connected ", conn.RemoteAddr())
	this.handle(conn)
	return nil
}

func (this *Client) read(conn net.Conn) {
	defer close(this.incoming)
	defer conn.Close()
	//
	var messageSize uint32
	var err error
	var readed int
	var buffer []byte
	this.reader = bufio.NewReader(conn)
	//
loop:
	for {
		// wait cann recv
		<-this.incoming
		//
		log.Println("receiving...")
		// read boundary
		readed, err = io.ReadAtLeast(this.reader,
			this.rheader[:],
			len(this.rheader))
		if err != nil || readed < 4 {
			log.Println("read bondary:", err.Error())
			break loop
		}
		// get size and type
		messageSize = this.rheader.GetSize()
		log.Println("message size:", messageSize)
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
		response := &message.Response{}
		message := buffer[0:readed]
		err := proto.Unmarshal(message, response)
		if err != nil {
			log.Println(err.Error())
		} else {
			log.Println(response.GetStatus())
		}
	}
	this.outgoing <- "closed"
}

func (this *Client) write(conn net.Conn) {
	defer close(this.outgoing)
	//
	var line string
	this.writer = bufio.NewWriter(conn)
loop:
	for {
		//
		_, err := fmt.Scanln(&line)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		//
		select {
		case <-this.outgoing:
			break loop
		default:
			// nil
		}
		var ok bool
		switch line {
		case "register":
			log.Println("register")
		default:
			log.Println("ping")
			ok = this.ping()
		}
		if ok != true {
			log.Println("error")
			return
		}
		this.incoming <- "yes"
	}
}

//
func (this *Client) register() bool {

	return false
}

//
func (this *Client) ping() bool {
	var n int
	request := &message.Request{}
	proto.SetDefaults(request)
	message, err := proto.Marshal(request)
	if err != nil {
		log.Println(err)
		return false
	}
	//size := proto.Size(request)
	size := len(message)
	log.Println("sending header...")
	//
	this.wheader.SetSize(uint32(size))
	log.Println(this.wheader)
	//
	n, err = this.writer.Write(this.wheader[:])
	if err != nil || n != len(this.wheader) {
		log.Println(err.Error())
		return false
	}
	log.Println(n, "bytes sent")
	//
	log.Println("sending message...")
	n, err = this.writer.Write(message)
	if err != nil || n != size {
		log.Println(err.Error())
		return false
	}
	log.Println(n, "bytes sent")
	err = this.writer.Flush()
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

//
func (this *Client) handle(conn net.Conn) {
	go this.read(conn)
	this.write(conn)
}

//
func NewClient() *Client {
	client := &Client{
		rheader:  server.NewTinyHeader(),
		wheader:  server.NewTinyHeader(),
		incoming: make(chan string),
		outgoing: make(chan string),
		reader:   nil,
		writer:   nil,
	}

	return client
}

func main() {
	// flag
	ip := flag.String("ip", "127.0.0.1", "server ip")
	port := flag.String("port", "9999", "server port")
	flag.Parse()
	address := *ip + ":" + *port
	//
	client := NewClient()
	client.Connect(address)
}
