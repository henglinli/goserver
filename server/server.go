package server

import (
	"log"
	"net"
)

type Handler interface {
	Serve(net.Conn)
}

const (	
	kGoRoutines = 4
)

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

func (server *Server) listen(handler Handler) {
	for i := 0; i < kGoRoutines; i++ {
		go func() {
			for conn := range server.connections {
				log.Println("serve... ", conn)
				handler.Serve(conn)
			}
		}()
	}
}

func (server *Server) Serve(handler Handler) error {
	// worker queue
	server.listen(handler)
	// listen
	ln, err := net.Listen("tcp", server.address)
	if err != nil {
		// handle error
		log.Println("net.Listen error: ", err.Error())
		return err
	}
	defer ln.Close()
	// accept
	//	go func() {
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			log.Println("net.Listen error: ", err.Error())
			continue
		}
		server.connections <- conn
	}
	//	}()

	return nil
}
