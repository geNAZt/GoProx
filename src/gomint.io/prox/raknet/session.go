package raknet

import (
	"net"
	"bytes"
	"encoding/binary"
	"gomint.io/prox/mcpe"
	"strconv"
	"gomint.io/prox/log"
)

type Session struct {
	addr        *net.UDPAddr
	conn        *net.UDPConn
	packetInput chan []byte
	guid        []byte
	state       int

	mtu int
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
				log.Warn("Could not get packet id: %v", err)
				continue
			}

			switch packetId {
			case OPEN_CONNECTION_REQUEST_1:
				self.handleConnectionRequest(buf)
				break
			case UNCONNECTED_PING:
				self.handleUnconnectedPing(buf)
				break
			}

			log.Debug("Got data: %v", data)
		}
	}
}

func (self *Session) write(buf *bytes.Buffer) {
	wroteBytes, oobBytes, err := self.conn.WriteMsgUDP(buf.Bytes(), []byte{}, self.addr)
	if err != nil {
		log.Warn("Could not write data: %v", err)
	}

	log.Debug("Wrote %v bytes %v oob", wroteBytes, oobBytes)
}

func isRaknetMagic(buf *bytes.Buffer) bool {
	// The rest should be the raknet magic
	rakNetMagic := make([]byte, 16)
	buf.Read(rakNetMagic)
	return bytes.Equal(rakNetMagic, OFFLINE_MAGIC)
}

func (self *Session) handleUnconnectedPing(buf *bytes.Buffer) {
	// Only valid when unconnected
	if self.state != STATE_UNCONNECTED {
		return
	}

	// Unconnected ping contains a ping time (long) and the raknet magic
	longBuf := make([]byte, 8)
	buf.Read(longBuf)
	ping := binary.BigEndian.Uint64(longBuf)

	if !isRaknetMagic(buf) {
		return
	}

	motd := "MCPE;§cGoProx §f- §6Golden performance;" + strconv.Itoa(mcpe.NETWORK_PROTOCOL) + ";" +
		mcpe.MINECRAFT_VERSION + ";0;200"

	// We answer with 0x1C (UNCONNECTED_PONG with MOTD) contains out of packetID + 2 longs + OFFLINE_ID + short + motd bytes
	answerBytes := make([]byte, 0, 35+len(motd))
	answer := bytes.NewBuffer(answerBytes)
	answer.WriteByte(UNCONNECTED_PONG_WITH_MOTD)
	answer.Write(longBuf)
	answer.Write(self.guid)
	answer.Write(OFFLINE_MAGIC)

	motdLengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(motdLengthBytes, uint16(len(motd)))
	answer.Write(motdLengthBytes)
	answer.WriteString(motd)

	log.Debug("Found ping %v", ping)
	self.write(answer)
}

func (self *Session) handleConnectionRequest(buf *bytes.Buffer) {
	// Only valid when unconnected
	if self.state != STATE_UNCONNECTED {
		return
	}

	self.state = STATE_INIT_CONNECTION

	// Needs to be a raknet magic next
	if !isRaknetMagic(buf) {
		return
	}

	protocolVersion, err := buf.ReadByte()
	if err != nil {
		log.Warn("Incompatible version of raknet connected: %v", err)
		self.sendIncompatibleVersion()
		self.state = STATE_UNCONNECTED
		return
	}

	if protocolVersion != PROTOCOL_VERSION {
		log.Warn("Incompatible version of raknet connected: %v", protocolVersion)
		self.sendIncompatibleVersion()
		self.state = STATE_UNCONNECTED
		return
	}

	self.mtu = buf.Len() + 18

	// Cap MTU
	if self.mtu > MAX_MTU {
		self.mtu = MAX_MTU
	}

	if self.mtu < MIN_MTU {
		self.mtu = MIN_MTU
	}

	log.Debug("MTU size: %v", self.mtu)
	self.sendConnectionReply1()
}

func (self *Session) sendConnectionReply1() {
	answerBytes := make([]byte, 0, 28)
	answer := bytes.NewBuffer(answerBytes)
	answer.WriteByte(OPEN_CONNECTION_REPLY_1)
	answer.Write(OFFLINE_MAGIC)
	answer.Write(self.guid)
	answer.WriteByte(0) // no libcat security

	mtuBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(mtuBytes, uint16(self.mtu))
	answer.Write(mtuBytes)

	self.write(answer)
}

func (self *Session) sendIncompatibleVersion() {
	answerBytes := make([]byte, 0, 26)
	answer := bytes.NewBuffer(answerBytes)
	answer.WriteByte(INCOMPATIBLE_VERSION)
	answer.WriteByte(8)
	answer.Write(OFFLINE_MAGIC)
	answer.Write(self.guid)
	self.write(answer)
}
