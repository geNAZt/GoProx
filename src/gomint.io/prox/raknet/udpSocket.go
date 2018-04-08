package raknet

import (
	"net"
	"strconv"
	"crypto/rand"
	"gomint.io/prox/log"
)

type UDPSocket struct {
	sessionManager *SessionManager
	guid           []byte
}

func ConstructUDPSocket() *UDPSocket {
	guid := make([]byte, 8)
	rand.Read(guid)

	return &UDPSocket{
		guid: guid,
		sessionManager: ConstructSessionManager(guid),
	}
}

func (self *UDPSocket) Listen(ip string, port int) {
	udpAddr, err := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("Could not resolve UDP address: %v", err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Could not lsiten to UDP address: %v", err)
	}

	for {
		buf := make([]byte, 65565)
		bytesRead, addr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Warn("Error while reading from UDP socket: %v", err)
			continue
		}

		copyBuf := buf[0:bytesRead]
		log.Debug("Addr: %v, Content: %v", addr, copyBuf)

		session := self.sessionManager.GetSession(udpConn, addr)
		if session != nil {
			log.Debug("Found session %v", session)
			session.PushDatagram(copyBuf)
		}
	}
}
