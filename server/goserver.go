package server

//
import (
	"../utils"
	"code.google.com/p/leveldb-go/leveldb"
	"log"
	"net"
	"time"
)

// go server
type GoServer struct {
	err         error
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
	log.Println("error:", this.err.Error())
	log.Println("stopping server ...")
	this.stop <- 1
}

//
func (this *GoServer) SetDefaultOptions(conn *net.TCPConn) {
	conn.SetKeepAlive(true)
	interval, err := time.ParseDuration("45s")
	if err != nil {
		conn.SetKeepAlivePeriod(interval)
	}
	conn.SetNoDelay(true)
}

// serve
func (this *GoServer) Serve() {
	go func() {
		// open db
		this.err = this.OpenDB(kDataBasePath)
		if this.err != nil {
			return
		}
		// handle connections
		this.listen()
		// clear
		defer this.clear()
		// listen remote
		laddr, err1 := net.ResolveTCPAddr("tcp", this.address)
		if err1 != nil {
			log.Println(err1.Error())
			this.err = err1
			return
		}
		// listen
		listener, err2 := net.ListenTCP("tcp", laddr)
		if err1 != nil {
			log.Println(err2.Error())
			this.err = err2
			return
		}
		// close listener
		defer listener.Close()
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
			conn, err := listener.AcceptTCP()
			if err != nil {
				// handle error
				log.Println(err.Error())
				this.err = err
				continue
			}
			//
			this.SetDefaultOptions(conn)
			//
			this.connections <- conn
		}
	}()
}
