package main

import (
	"fmt"

	"github.com/WirvsVirus-DeMed/backend/db"
)

func main() {
	fmt.Println("DB TEST")

	db.CreateMedicine()
}
