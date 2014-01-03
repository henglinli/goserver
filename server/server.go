package server

import (
	"log"
	"net"
)

type ConnectionHandler interface {
	Handle(net.Conn)
}

type Server struct {
	stop chan int
	address string
	connections chan net.Conn
	clients int64
}

func NewServer(addr string) *Server {
	server := &Server {
		stop: make(chan int, 1),
		address: addr,
		connections: make(chan net.Conn),
		clients: 0,
	}
	
	return server
}

// connected clients
func (server *Server) Clients() int64 {	
	return server.clients
}

func (server *Server) listen(handler ConnectionHandler) {
	go func() {
		for conn := range server.connections {
			handler.Handle(conn)
		}
	}()
}

func (server *Server) Stop() {
	defer close(server.stop)
	//
	log.Println("stopping server...")
	server.stop <- 1
}
//
func (server *Server) Serve(handler ConnectionHandler) error {
	var err error
	err = nil
	//
	server.listen(handler)
	//
	go func() {	
		// close connections
		defer func() {
			for conn := range server.connections {
				conn.Close()
			}
		}()
		// listen
		ln, e := net.Listen("tcp", server.address)
		if err != nil {
			// handle error
			log.Println("net.Listen error: ", err.Error())
			err = e
			return
		}
		// close listener
		defer ln.Close()
		// accept
		for {
			// check should stop
			select {
			case <- server.stop:
				return
			default:
				// continue
			}
			// accept
			conn, err := ln.Accept()
			if err != nil {
				// handle error
				log.Println("net.Listen error: ", err.Error())
				continue
			}
			server.connections <- conn
		}
	}()
	return err
}
