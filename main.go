package main

import (
	"github.com/WirvsVirus-DeMed/backend/api"
	"github.com/WirvsVirus-DeMed/backend/db"
)

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
	database, _ := db.CreateDataBase()
	defer database.Close()

	db.CreateMedicineTable(database)
	// med := &db.Medicine{"1", "21", "1", db.Peer{}, 1, 1}
	// med2 := &db.Medicine{"2", "31", "1", db.Peer{}, 1, 2}

	// med.Add(database)
	// med2.Add(database)
	// med = &db.Medicine{"1", "2", "1", time.Now(), "1", 1, 1}
	// med.Update(database)
	// db.GetAll(database)
	// db.Get(database, "3")

	api.Api(database)
}
