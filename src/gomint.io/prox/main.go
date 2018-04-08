package main

import (
	"gomint.io/prox/config"
	"gomint.io/prox/raknet"
	"gomint.io/prox/log"
)

func main() {
	// Get config
	log.Info("Going to bind %v", config.Config.Listener)

	// Create network
	socket := raknet.ConstructUDPSocket()
	socket.Listen(config.Config.Listener.IP, config.Config.Listener.Port)
}
