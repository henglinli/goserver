package server

import ( 
	"log"
	"net"
	"bufio"
	"testing"
)
//
type EchoHandler struct {
	incoming chan string
	outgoing chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

//
func (echo *EchoHandler) Read() {
	for {
		line, _ := echo.reader.ReadString('\n')
		echo.incoming <- line
	}
}

//
func (echo *EchoHandler) Write() {
	for data := range echo.outgoing {
		echo.writer.WriteString(data)
		echo.writer.Flush()
	}
}

//
func (echo *EchoHandler) Listen() {
	go echo.Read()
	go echo.Write()
}

//
func NewEchoHandler(connection net.Conn) *EchoHandler {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	echo := &EchoHandler{
		incoming: make(chan string),
		outgoing: make(chan string),
		reader:   reader,
		writer:   writer,
	}
	
	echo.Listen()
	
	return echo
}

func (echo *EchoHandler) Handle(conn net.Conn) {

	echo.Listen()
}

func TestServer(t *testing.T) {
	log.Println("testing Server...")
}
