package main

import (
	"log"
	"time"

	"github.com/WirvsVirus-DeMed/backend/db"
)

func main() {
	// fmt.Println("DB TEST")

	// med := &db.Medicine{"1", "21", "1", time.Now(), "1", 1, 1}
	// med2 := &db.Medicine{"2", "31", "1", time.Now(), "1", 1, 2}

	// database, err := db.CreateDataBase()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// db.CreateMedicineTable(database)
	// med.Add(database)
	// med2.Add(database)
	// med = &db.Medicine{"1", "2", "1", time.Now(), "1", 1, 1}
	// med.Update(database)
	// // db.GetAll(database)
	// db.Get(database, "3")
	// database.Close()

	var a []byte = []byte("1")
	var b []byte = []byte("2")
	var c []byte = []byte("3")

	peer1 := &db.Peer{'1', a, 443, time.Now()}
	peer2 := &db.Peer{'2', b, 80, time.Now()}

	database, err := db.CreateDataBase()
	if err != nil {
		log.Fatal(err)
		return
	}

	db.CreatePeerTable(database)
	peer1.Add(database)
	peer2.Add(database)
	peer := &db.Peer{'1', c, 443, time.Now()}
	peer.Update(database)
	db.GetAllPeers(database)
	// db.Get(database, "1")
	database.Close()
}
