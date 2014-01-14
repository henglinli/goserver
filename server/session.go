package server

import (
	"sync"
)

// impl Session
type AddrSession struct {
	address string
	manager SessionManager
}

//
func (this *AddrSession) IsLogin() bool {
	return this.manager.IsLogin(this)
}

//
func (this *AddrSession) Name() string {
	return this.address
}

// impl SessionManager
type AddrSessionManager struct {
	rwlock     sync.RWMutex
	sessionmap map[string]Session
}

//
func (this *AddrSessionManager) NewSession(in string) Session {
	s := &AddrSession{
		address: in,
		manager: this,
	}

	return s
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
func (this *AddrSessionManager) IsLogin(s Session) bool {
	this.rwlock.RLock()
	_, ok := this.sessionmap[s.Name()]
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
func NewAddrSessionManager() *AddrSessionManager {
	m := &AddrSessionManager{
		rwlock:     sync.RWMutex{},
		sessionmap: make(map[string]Session, kMaxOnlineClients),
	}

	return m
}
