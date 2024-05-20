package main

import (
	"fmt"
	"net"
	"netcat/server"
	"os"
	"sync"
)

type Client struct {
	Name string
	Conn net.Conn
	// Other relevant client information
}

var connectedClients = make(map[*Client]bool)
var clientsMutex sync.Mutex // to synchronize access to the map

func main() {
	args := os.Args
	port := "8989"

	if len(args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	if len(args) == 2 {
		port = args[1]
	}

	address := server.GetIP() + ":" + port
	serverInstance := &server.Server{}
	serverInstance.StartServer(address)

	serverInstance.CloseServer()
	//server.BroadCastMessage
}
