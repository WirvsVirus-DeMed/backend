package main

import (
	"fmt"
	"log"
	"time"

	"github.com/WirvsVirus-DeMed/backend/db"
)

func main() {
	fmt.Println("DB TEST")

	med := &db.Medicine{"1", "1", "1", time.Now(), "1"}
	med2 := &db.Medicine{"2", "1", "1", time.Now(), "1"}

	database, err := db.CreateDataBase()
	if err != nil {
		log.Fatal(err)
		return
	}

	db.CreateMedicineTable(database)
	med.Add(database)
	med2.Add(database)
	med = &db.Medicine{"1", "2", "1", time.Now(), "1"}
	med.Update(database)
	med.Get(database)
	database.Close()
}
