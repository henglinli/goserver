package server

//
import (
	"log"
	"net"
)

// ===========================================
// server options
const (
	kMaxOnlineClients = 2048
	kBufferSize       = 2048
)

// ===========================================
// interface
// session
type Session interface {
	// check session
	IsLogin() bool
	// session name
	Name() string
}

//
type SessionManager interface {
	// new
	NewSession(string) Session
	// login
	Login(Session)
	// size
	Sessions() int
	// is login
	IsLogin(Session) bool
	// logout
	Logout(Session)
}

// message handler
type MessageHandler interface {
	// handle the mesage
	//Handle([]byte) []byte
	Handle([]byte, Session) []byte
}

// connection handler
type ConnectionHandler interface {
	// handle the connection
	Handle(net.Conn)
}

// ===========================================
//
type Server struct {
	stop        chan int
	address     string
	connections chan net.Conn
	manager     SessionManager
}

//
func NewServer(addr string) *Server {
	server := &Server{
		stop:        make(chan int, 1),
		address:     addr,
		connections: make(chan net.Conn),
	}
	//
	return server
}

//
func (this *Server) listen(handler ConnectionHandler) {
	go func() {
		for conn := range this.connections {
			handler.Handle(conn)
		}
	}()
}

//
func (this *Server) Stop() {
	defer close(this.stop)
	//
	log.Println("stopping server...")
	this.stop <- 1
}

//
func (this *Server) Serve(handler ConnectionHandler) error {
	var err error
	err = nil
	// handle connections
	this.listen(handler)
	//
	go func() {
		// close connections
		defer func() {
			for conn := range this.connections {
				conn.Close()
			}
		}()
		// listen
		ln, e := net.Listen("tcp", this.address)
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
			select {
			// check should stop
			case <-this.stop:
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
			//
			this.connections <- conn
		}
	}()
	//
	return err
}
