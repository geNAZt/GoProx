package raknet

import (
	"net"
	"log"
	"bytes"
	"encoding/binary"
)

type Session struct {
	addr        *net.UDPAddr
	conn        *net.UDPConn
	packetInput chan []byte
	guid        []byte
}

func ConstructSession(guid []byte, conn *net.UDPConn, addr *net.UDPAddr) *Session {
	// Create new session and return
	session := &Session{
		guid:        guid,
		addr:        addr,
		conn:        conn,
		packetInput: make(chan []byte, 10),
	}

	go session.ProcessPackets()
	return session
}

func (self *Session) PushDatagram(bytes []byte) {
	self.packetInput <- bytes
}

func (self *Session) ProcessPackets() {
	for {
		select {
		case data := <-self.packetInput:
			buf := bytes.NewBuffer(data)

			// Check for packet id
			packetId, err := buf.ReadByte()
			if err != nil {
				log.Printf("Could not get packet id: %v\n", err)
				continue
			}

			switch packetId {
			case UNCONNECTED_PING:
				self.handleUnconnectedPing(buf)
			}

			log.Printf("Got data: %v\n", data)
		}
	}
}

func (self *Session) handleUnconnectedPing(buf *bytes.Buffer) {
	// Unconnected ping contains a ping time (long) and the raknet magic
	longBuf := make([]byte, 8)
	buf.Read(longBuf)
	ping := binary.BigEndian.Uint64(longBuf)

	// The rest should be the raknet magic
	rakNetMagic := make([]byte, 16)
	buf.Read(rakNetMagic)
	if !bytes.Equal(rakNetMagic, OFFLINE_MAGIC) {
		return
	}

	// We answer with 0x1C (UNCONNECTED_PONG with MOTD) contains out of packetID + 2 longs + OFFLINE_ID + short + motd bytes
	motd := "§cGoProx §f- §6Golden performance"
	answerBytes := make([]byte, 0, 35+len(motd))
	answer := bytes.NewBuffer(answerBytes)
	answer.WriteByte(UNCONNECTED_PONG_WITH_MOTD)
	answer.Write(longBuf)
	answer.Write(self.guid)

	motdLengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(motdLengthBytes, uint16(len(motd)))
	answer.Write(motdLengthBytes)
	answer.Write([]byte(motd))

	log.Printf("Found ping %v\n", ping)

	self.conn.WriteMsgUDP(answer.Bytes(), []byte{}, self.addr)
}
