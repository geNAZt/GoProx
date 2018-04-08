package main

import (
	"log"
	"gomint.io/prox/config"
	"gomint.io/prox/raknet"
)

func main() {
	// Get config
	log.Printf("Going to bind %v\n", config.Config.Listener)

	// Create network
	socket := raknet.ConstructUDPSocket()
	socket.Listen(config.Config.Listener.IP, config.Config.Listener.Port)
}
