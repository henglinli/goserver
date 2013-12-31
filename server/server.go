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
	handlers chan Handler
	connections chan net.Conn
	counter chan int64
	clients int64
}

func NewServer(addr string) *Server {
	server := &Server {
		address: addr,
		handlers: make(chan Handler),
		connections: make(chan net.Conn),
		counter: make(chan int64),
		clients: 0,
	}
	
	return server
}

// connected clients
func (server *Server) Clients() int64 {
	for change := range server.counter {
		server.clients += change
	}
	return server.clients
}

func (server *Server) Serve() error {
	// worker queue
	for i := 0; i < kGoRoutines; i++ {
		go func() {
			for conn := range server.connections {
				for handler := range server.handlers {
					server.counter <- 1
					handler.Serve(conn)
					server.counter <- -1
				}
			}
		}()
	}
	// listen
	ln, err := net.Listen("tcp", server.address)
	if err != nil {
		// handle error
		log.Println("net.Listen error: ", err.Error())
		return err
	}
	// accept
	go func() {
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

	return nil
}
