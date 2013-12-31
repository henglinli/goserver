package server
//
import (
	"bufio"
	"net"
	"log"
	"client"
)
//
const (
	kRecvBufLen = 2048
)

//
type Handler interface {
	Read(message string) error
	Write(message string) error
}
//
type Server struct {
	clients []*client.Client
	joins chan net.Conn
	incoming chan string
	outgoing chan string
}
//
func (service *Server) response(message string) {
	data
}
//
func (service *Server) join(conn net.Conn) {
	
}
//
func (service *Server) listen() error {
	go func() {
		for {
			select {
			case data := service.incoming;
				service.response(data)
			case conn := <-service.jions:
				serice.join(conn)
			}
		}
	}()
}
//
func (service *Server) Serve() error {
	err := service.listen()
	if err {
		return err
	}
	listener, err := net.Listen("tcp", service.address)
	if err {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err {
			log.Println("Accept error: ", err.Error())
		}
		service.joins <-conn
	}
}
//
func NewServer(string address) *Server {
	service := &Server{
		clients: make([]*chat.Client, 0)
		joins: make(chan net.Conn)
		incoming: make(chan string)
		outgoing: make(chan string)
	}
	return service
}
