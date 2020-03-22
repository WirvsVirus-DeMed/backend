package main

import "github.com/WirvsVirus-DeMed/backend/api"

func main() {
	// log.SetOutput(os.Stdout)

	// _, _ = db.CreateDataBase()

	// client, _ := strconv.Atoi(os.Args[1])

	// fmt.Println("Hello World!")

	// n := node.Node{}
	// n.Init(
	// 	"certs/client"+strconv.Itoa(client)+".crt",
	// 	"certs/client"+strconv.Itoa(client)+".key",
	// 	"certs/rootCA.crt",
	// 	"client"+strconv.Itoa(client),
	// 	uint32(10011+client),
	// 	100,
	// 	100,
	// 	8,
	// 	10*time.Second)

	// go n.Listen()
	// go n.HandleMessages()
	// go n.BroadcastSender()
	// go n.PeerDiscovery()
	// go n.TidyMedicineOffersRoutine()

	// if client != 0 {
	// 	n.Connect(net.IPv4(127, 0, 0, 1), 10011)
	// }

	// for {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	text, _ := reader.ReadString('\n')
	// 	if strings.ToLower(text) == "d\n" {
	// 		log.Printf("[d] CurrentConnections: %v\n", n.Clients.Len())
	// 	} else if strings.ToLower(text) == "unban\n" {

	// 		for e := n.PeerBlackList.Front(); e != nil; e = e.Next() {
	// 			n.PeerBlackList.Remove(e)
	// 		}
	// 	} else if strings.ToLower(text) == "q\n" {
	// 		return
	// 	}
	// }
	api.Api()
}
