package server

import (
	"../utils"
	"bufio"
	"log"
	"net"
	"testing"
)

//
type Client struct {
	incoming chan string
	outgoing chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

//
func (client *Client) Read() {
	for {
		line, err := client.reader.ReadString('\n')
		if err != nil {
			log.Println("reader.ReadString error: ", err.Error())
			close(client.incoming)
			close(client.outgoing)
			break
		}
		client.incoming <- line
	}
	log.Println("reader.Read done")
}

//
func (client *Client) Write() {
	for data := range client.outgoing {
		client.writer.WriteString(data)
		client.writer.Flush()
	}
	log.Println("writer.Write done")
}

//
func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

//
func NewClient(connection net.Conn) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		incoming: make(chan string),
		outgoing: make(chan string),
		reader:   reader,
		writer:   writer,
	}

	client.Listen()

	return client
}

type EchoHandler struct {
	// nil
}

func (echo *EchoHandler) Handle(conn net.Conn) {
	log.Println("EchoHandler.Serve ", conn.RemoteAddr())
	go func() {
		defer conn.Close()
		client := NewClient(conn)
	loop:
		for {
			log.Println("recv...")
			select {
			case recved, ok := <-client.incoming:
				log.Println(ok, " recv: ", len(recved))
				if ok == false {
					break loop
				}
				client.outgoing <- recved
				// client.outgoing <- <- client.incoming
			}
		}
		log.Println("EchoHandler.Serve done")
	}()
	log.Println("EchoHandler.Serve started")
}

func TestServer(t *testing.T) {
	if false {
		log.Println("testing Server...")

		echo := new(EchoHandler)
		server := NewServer(":9999")
		server.Serve(echo)
		//
		utils.Wait()
	}
}

type EchoMessage struct {
	// nil
}

func (handler *EchoMessage) Handle(message []byte) []byte {
	// nil
	return message
}

func TestDemuxer(t *testing.T) {

	log.Println("testing Server Demuxer...")
	message := new(EchoMessage)
	handler := NewDemuxerHandler(message)
	server := NewServer(":9999")
	server.Serve(handler)
	//
	utils.Wait()
	server.Stop()
}
