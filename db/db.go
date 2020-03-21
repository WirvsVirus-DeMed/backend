package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Praxis
// Resource

func IsThisRight() {
	fmt.Println("TEST")
}

func CreateMedicine() {
	os.Remove("./meds.db")
	db, err := sql.Open("sqlite3", "./meds.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table med (id integer not null primary key, name text);
	delete from med;
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// tx, err := db.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// stmt, err := tx.Prepare("insert into med(id, name) values(0, 'test')")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
}
