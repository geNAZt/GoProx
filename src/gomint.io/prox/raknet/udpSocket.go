package raknet

import (
	"net"
	"strconv"
	"log"
	"crypto/rand"
)

type UDPSocket struct {
	sessionManager *SessionManager
	guid           []byte
}

func ConstructUDPSocket() *UDPSocket {
	guid := make([]byte, 16)
	rand.Read(guid)

	return &UDPSocket{
		guid: guid,
		sessionManager: ConstructSessionManager(guid),
	}
}

func (self *UDPSocket) Listen(ip string, port int) {
	udpAddr, err := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Could not resolve UDP address: %v", err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("Could not lsiten to UDP address: %v", err)
	}

	for {
		buf := make([]byte, MAX_MTU)
		bytesRead, addr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error while reading from UDP socket: %v\n", err)
			continue
		}

		copyBuf := buf[0:bytesRead]
		log.Printf("Addr: %v, Content: %v\n", addr, copyBuf)

		session := self.sessionManager.GetSession(udpConn, addr)
		if session != nil {
			log.Printf("Found session %v\n", session)
			session.PushDatagram(copyBuf)
		}
		// net.DialUDP("udp", udpAddr, addr)
	}
}
