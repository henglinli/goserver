package server

import ( 
	"log"
	"net"
	"bufio"
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
	forward chan string
}

func (echo *EchoHandler) Serve(conn net.Conn) {
	log.Println("EchoHandler.Serve")
	go func() {
		defer conn.Close()
		client := NewClient(conn)
	loop:   for {
			log.Println("recv...")
			select {
			case recv, ok := <- client.incoming:
				log.Println(ok, " recv: ", len(recv))
				if ok == false {
					break loop
				}
				client.outgoing <- recv
				// client.outgoing <- <- client.incoming
			}
			log.Println("re handle...")
		}
		log.Println("EchoHandler.Serve done")
	}()
}

func NewEchoHandler() *EchoHandler {
	echo := &EchoHandler {
		forward: make(chan string),
	}

	return echo
}

func TestServer(t *testing.T) {
	log.Println("testing Server...")
	
	echo := NewEchoHandler()
	server := NewServer(":9999")
	server.Serve(echo)
}
