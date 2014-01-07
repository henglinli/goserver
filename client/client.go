package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	//"code.google.com/p/goprotobuf/proto"
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
	incoming chan string
	outgoing chan string
}

func (client *Client) Connect(address string) error {
	log.Println("connecting to ", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
		log.Println(err.Error())
	}
	log.Println("connectd ", conn.RemoteAddr())
	client.handle(conn)
	return nil
}

func (client *Client) read(conn net.Conn) {
	defer close(client.incoming)
	defer conn.Close()
	//
	bondary := [4]byte{}
	var messageSize uint32
	var err error
	var readed int
	var buffer []byte
	reader := bufio.NewReader(conn)
	//
	//message := proto.NewBuffer(nil)
	//
loop:
	for {
		// wait cann recv
		<-client.incoming
		//
		log.Println("receiving...")
		// read boundary
		readed, err = io.ReadAtLeast(reader, bondary[0:4], 4)
		if err != nil || readed < 4 {
			log.Println("read bondary:", err.Error())
			break loop
		}
		// get size and type
		messageSize = GetSize(bondary[0:4])
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
	client.outgoing <- "closed"
}

func (client *Client) write(conn net.Conn) {
	defer close(client.outgoing)
	//
	var line string
	writer := bufio.NewWriter(conn)
	bondary := [8]byte{}
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
		case <-client.outgoing:
			break loop
		default:
			// nil
		}
		size := len(line)
		log.Println("sending bondary...")
		//
		SetSize(uint32(size), bondary[0:4])
		//
		n, err = writer.Write(bondary[0:4])
		if err != nil || n != 4 {
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
		client.incoming <- "yes"
	}
}

//
func (client *Client) handle(conn net.Conn) {
	go client.read(conn)
	client.write(conn)
}

//
func NewClient() *Client {
	client := &Client{
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
