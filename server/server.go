package server

import (
	"log"
	"net"
)

type ConnectionHandler interface {
	Handle(net.Conn)
}

type Server struct {
	address string
	connections chan net.Conn
	clients int64
}

func NewServer(addr string) *Server {
	server := &Server {
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

func (server *Server) Serve(handler ConnectionHandler) error {
	var err error
	err = nil
	
	go func() {
		server.listen(handler)
	// listen
		ln, e := net.Listen("tcp", server.address)
		if err != nil {
			// handle error
			log.Println("net.Listen error: ", err.Error())
			err = e
			return
		}
		defer ln.Close()
		// accept
		for {
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
