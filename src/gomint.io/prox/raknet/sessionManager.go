package raknet

import (
	"sync"
	"net"
)

type SessionManager struct {
	sessions *sync.Map
	guid []byte
}

func ConstructSessionManager(guid []byte) *SessionManager {
	return &SessionManager{
		guid: guid,
		sessions: &sync.Map{

		},
	}
}

func (self *SessionManager) GetSession(conn *net.UDPConn, addr *net.UDPAddr) *Session {
	// Check if we already have a connection
	if val, ok := self.sessions.Load(addr); ok {
		return val.(*Session)
	}

	// Construct new session
	session := ConstructSession(self.guid, conn, addr)
	if session != nil {
		self.sessions.Store(addr, session)
	}

	return session
}
