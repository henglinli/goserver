package server

//
import (
	"../utils"
	"code.google.com/p/leveldb-go/leveldb"
	"log"
	"net"
)

// go server
type GoServer struct {
	stop        chan int
	address     string
	connections chan net.Conn
	manager     SessionManager
	db          *leveldb.DB
}

//
func NewGoServer(addr string) Server {
	server := &GoServer{
		stop:        make(chan int, 1),
		address:     addr,
		connections: make(chan net.Conn),
		manager:     nil,
	}
	//
	return server
}

//
func Go(addr string, handler MessageHandler) {
	// new server
	server := NewGoServer(addr)
	// new session manager
	server.NewSessionManager(handler)
	// serve
	server.Serve()
	// wait signal
	utils.Wait()
	// stop
	server.Stop()
}

//
func (this *GoServer) NewSessionManager(handler MessageHandler) {
	this.manager = NewAddrSessionManager(handler)
}

// open db
func (this *GoServer) OpenDB(name string) error {
	log.Println("open db ...", name)
	var err error
	this.db, err = leveldb.Open(name, nil)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// close db
func (this *GoServer) CloseDB() {
	log.Println("close db ...")
	if this.db != nil {
		this.db.Close()
	}
}

//
func (this *GoServer) listen() {
	go func() {
		for conn := range this.connections {
			this.manager.NewSession(conn, this.db)
		}
	}()
}

// clear
func (this *GoServer) clear() {
	// reset sessions
	this.manager.Reset()
	// close connections
	for conn := range this.connections {
		conn.Close()
	}
	// clsoe db
	this.CloseDB()
}

// stop
func (this *GoServer) Stop() {
	defer close(this.stop)
	//
	log.Println("stopping server...")
	this.stop <- 1
}

// serve
func (this *GoServer) Serve() error {
	var err error
	go func() {
		// open db
		err = this.OpenDB(kDataBasePath)
		if err != nil {
			return
		}
		// handle connections
		this.listen()
		// clear
		defer this.clear()
		// listen remote
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
