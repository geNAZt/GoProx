package raknet

const (
	// MTU stuff
	MAX_MTU = 1464 	// Seems to be the maximum value MCPE supports
	MIN_MTU = 548	// Seems to be the minimum value MCPE supports

	PROTOCOL_VERSION = 8

	// Packet IDs
	UNCONNECTED_PING = 0x01
	OPEN_CONNECTION_REQUEST_1 = 0x05
	OPEN_CONNECTION_REPLY_1 = 0x06
	INCOMPATIBLE_VERSION = 0x19
	UNCONNECTED_PONG_WITH_MOTD = 0x1C
)

var (
	// Raknet offline ID
	OFFLINE_MAGIC = []byte{0x00, 0xFF, 0xFF, 0x00, 0xFE, 0xFE, 0xFE, 0xFE, 0xFD, 0xFD, 0xFD, 0xFD, 0x12, 0x34, 0x56, 0x78}
)
