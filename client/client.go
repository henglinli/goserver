package main

import (
	"flag"
	"log"
	"net"
	"fmt"
	"io"	
	"bufio"
	"encoding/binary"
)
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

type Client struct {
	incoming chan string
	outgoing chan string
}

func (client *Client) Connect (address string) error {
	log.Println("connecting to ", address)
	conn , err := net.Dial("tcp", address)
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
	bondary := [8]byte{}
	var message_size, message_type uint32
	var err error
	var readed int
	var buffer []byte
	reader := bufio.NewReader(conn)
	//
	loop:	for {
		// wait cann recv
		<- client.incoming 
		//
		log.Println("receiving...")
		// read boundary
		readed, err = io.ReadAtLeast(reader, bondary[0:8], 8)
		if err != nil || readed < 8 {
			log.Println("read bondary:", err.Error())
			break loop			
		}
		// get size and type
		message_size, message_type = GetBondary(bondary[0:8])
		log.Println("message size:", message_size)
		log.Println("message type:", message_type)
		// get buffer
		needed := len(buffer) - int(message_size)
		if needed < 0 {
			log.Println("Growing buffer...")
			buffer = make([]byte, message_size)
		}
		// get message
		readed, err = io.ReadAtLeast(reader, buffer, int(message_size))
		if err != nil || readed < int(message_size) {
			log.Println("read message:", err.Error())
			break loop		
		}
		//
		log.Println("message: ", string(buffer[0:readed]), readed)
	}
	client.outgoing <- "closed"
}

func (client *Client) write(conn net.Conn) {
	defer close(client.outgoing)
	//
	var line string 
	writer := bufio.NewWriter(conn)
	bondary := [8]byte{}
	loop: for {	
		//	
		n, err := fmt.Scanln(&line)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		//
		select {
		case <- client.outgoing:
			break loop
		default:
			// nil
		}
		size := len(line)
		log.Println("sending bondary...")
		SetBondary(uint32(size), 0, bondary[0:8])
		n, err = writer.Write(bondary[0:8])
		if err != nil || n != 8 {
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

func (client *Client) handle(conn net.Conn) {
	go client.read(conn)
	client.write(conn)
}

func NewClient() *Client{
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
