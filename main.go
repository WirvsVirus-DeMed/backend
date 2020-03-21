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

	db, err := db.CreateDataBase()
	if err != nil {
		log.Fatal(err)
		return
	}

	med.CreateMedicineTable(db)
	med.AddMedicine(db)
	med2.AddMedicine(db)
	med.GetMedicine(db)
	db.Close()
}
