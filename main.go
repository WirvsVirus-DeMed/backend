package main

import (
	"WirvsVirus/DeMed/node"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	log.SetOutput(os.Stdout)

	client, _ := strconv.Atoi(os.Args[1])

	fmt.Println("Hello World!")

	n := node.Node{}
	n.Init(
		"certs/client"+strconv.Itoa(client)+".crt",
		"certs/client"+strconv.Itoa(client)+".key",
		"certs/rootCA.crt",
		"client"+strconv.Itoa(client),
		uint32(5020+client),
		100,
		100)

	go n.Listen()
	go n.HandleMessages()
	go n.BroadcastSender()

	n.Connect(net.IPv4(127, 0, 0, 1), 5000)

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if strings.ToLower(text) == "d\n" {
			log.Printf("[d] CurrentConnections: %v\n", n.Clients.Len())
		} else if strings.ToLower(text) == "q\n" {
			return
		}
	}
}
