package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"../message"
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
	rheader *message.DefaultHeader
	wheader *message.DefaultHeader
	incoming chan string
	outgoing chan string
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
	reader := bufio.NewReader(conn)
	//
loop:
	for {
		// wait cann recv
		<-this.incoming
		//
		log.Println("receiving...")
		// read boundary
		readed, err = io.ReadAtLeast(reader, 
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
		readed, err = io.ReadAtLeast(reader, buffer, int(messageSize))
		if err != nil || readed < int(messageSize) {
			log.Println("read message:", err.Error())
			break loop
		}
		//
		// message.DebugPrint("message", buffer[0:readed])
		log.Println(string(buffer[0:readed]))
	}
	this.outgoing <- "closed"
}

func (this *Client) write(conn net.Conn) {
	defer close(this.outgoing)
	//
	var line string
	writer := bufio.NewWriter(conn)
loop:
	for {
		//
		n, err := fmt.Scanln(&line)
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
		size := len(line)
		log.Println("sending header...")
		//
		this.wheader.SetSize(uint32(size))
		log.Println(this.wheader)
		//
		n, err = writer.Write(this.wheader[:])
		if err != nil || n != len(this.wheader) {
			log.Println(err.Error())
			break loop
		}
		log.Println(n, "bytes sent")
		//
		log.Println("sending message...")
		n, err = writer.WriteString(line)
		if err != nil || n != size {
			log.Println(err.Error())
			break loop
		}
		log.Println(n, "bytes sent")
		err = writer.Flush()
		if err != nil {
			log.Println(err.Error())
			break loop
		}
		//
		this.incoming <- "yes"
	}
}

//
func (this *Client) handle(conn net.Conn) {
	go this.read(conn)
	this.write(conn)
}

//
func NewClient() *Client {
	client := &Client{
		rheader: message.NewDefaultHeader(),
		wheader: message.NewDefaultHeader(),
		incoming: make(chan string),
		outgoing: make(chan string),
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
