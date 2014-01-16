package server

//
import (
	"code.google.com/p/leveldb-go/leveldb"
	"net"
)

// ===========================================
// interface
// message handler
type MessageHandler interface {
	// handle the mesage
	//Handle([]byte) []byte
	Handle([]byte, Session) []byte
}

// session
type Session interface {
	// check session
	IsLogin() bool
	// session name
	Name() string
	// handle session, read and write connection
	Handle([]byte)
	//
	GetDB() *leveldb.DB
}

// header coder
type HeaderCoder interface {
	// get message size
	GetSize() uint32
	// set message size
	SetSize(uint32)
}

// session manager
type SessionManager interface {
	// called by Server.Serve
	// create new session
	NewSession(net.Conn, *leveldb.DB)
	// reset all connections
	Reset()
	// called by demuxer
	// login
	Login(Session)
	// size
	Sessions() int
	// is login
	IsLogin(string) bool
	// logout
	Logout(Session)
}

// server interface
type Server interface {
	// new SessionManager
	NewSessionManager(MessageHandler)
	// open db
	OpenDB(string) error
	// serve connections
	Serve()
	// close db
	CloseDB()
	// stop Server
	Stop()
}
