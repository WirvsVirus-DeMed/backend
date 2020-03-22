package main

import (
	"github.com/WirvsVirus-DeMed/backend/api"
	"github.com/WirvsVirus-DeMed/backend/db"
)

func main() {
	// log.SetOutput(os.Stdout)

	// client, _ := strconv.Atoi(os.Args[1])

	// fmt.Println("Hello World!")

	// n := node.Node{}
	// n.Init(
	// 	"certs/client"+strconv.Itoa(client)+".crt",
	// 	"certs/client"+strconv.Itoa(client)+".key",
	// 	"certs/rootCA.crt",
	// 	"client"+strconv.Itoa(client),
	// 	uint32(5000+client),
	// 	100,
	// 	100)

	// go n.Listen()
	// go n.HandleMessages()
	// go n.BroadcastSender()

	// n.Connect(net.IPv4(127, 0, 0, 1), 5000)

	// for {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	text, _ := reader.ReadString('\n')
	// 	if strings.ToLower(text) == "d\n" {
	// 		log.Printf("[d] CurrentConnections: %v\n", n.Clients.Len())
	// 	} else if strings.ToLower(text) == "q\n" {
	// 		return
	// 	}
	// }

	database, _ := db.CreateDataBase()
	db.CreateMedicineTable(database)
	med := &db.Medicine{"1", "21", "1", db.Peer{}, 1, 1}
	med2 := &db.Medicine{"2", "31", "1", db.Peer{}, 1, 2}

	med.Add(database)
	med2.Add(database)
	// med = &db.Medicine{"1", "2", "1", time.Now(), "1", 1, 1}
	// med.Update(database)
	// db.GetAll(database)
	// db.Get(database, "3")

	api.Api(database)
	database.Close()
}
