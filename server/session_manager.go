package server

import (
	"code.google.com/p/leveldb-go/leveldb"
	"log"
	"net"
	"sync"
)

// impl SessionManager
type AddrSessionManager struct {
	rwlock     sync.RWMutex
	sessionmap map[string]Session
	buffers    *BufferManager
	handler    MessageHandler
}

// impl SessionManager
func (this *AddrSessionManager) NewSession(conn net.Conn, db *leveldb.DB) {
	log.Println("AddrSession.NewSession", conn.RemoteAddr())
	go func() {
		// new session
		s := NewAddrSession(conn, db, this, this.handler)
		// login session
		this.Login(s)
		// get buffer
		buffer := this.buffers.Get()
		// user buffer
		s.Handle(buffer)
		// put buffer back
		this.buffers.Put(buffer)
		// lgoout session
		this.Logout(s)
	}()
	log.Println("AddrSession.NewSession started")
}

//
func (this *AddrSessionManager) Reset() {
	this.rwlock.Lock()
	for session := range this.sessionmap {
		delete(this.sessionmap, session)
	}
	this.rwlock.Unlock()
}

//
func (this *AddrSessionManager) Login(s Session) {
	this.rwlock.Lock()
	this.sessionmap[s.Name()] = s
	this.rwlock.Unlock()
}

//
func (this *AddrSessionManager) Sessions() int {
	this.rwlock.Lock()
	sessions := len(this.sessionmap)
	this.rwlock.Unlock()
	return sessions
}

//
func (this *AddrSessionManager) IsLogin(name string) bool {
	this.rwlock.RLock()
	_, ok := this.sessionmap[name]
	this.rwlock.RUnlock()
	return ok
}

//
func (this *AddrSessionManager) Logout(s Session) {
	this.rwlock.Lock()
	delete(this.sessionmap, s.Name())
	this.rwlock.Unlock()
}

//
func NewAddrSessionManager(handler MessageHandler) *AddrSessionManager {
	m := &AddrSessionManager{
		rwlock:     sync.RWMutex{},
		sessionmap: make(map[string]Session, KMaxOnlineClients),
		buffers:    NewBufferManager(),
		handler:    handler,
	}

	return m
}
