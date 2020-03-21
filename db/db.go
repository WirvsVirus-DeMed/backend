package db

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3" // Needs to be like this
)

// Medicine is the stuff that makes you healthy
type Medicine struct {
	UUID       string    `json:"uuid"`
	Title      string    `json:"title"`
	Desciption string    `json:"desciption"`
	CreatedAt  time.Time `json:"createdAt"`
	Owner      string    `json:"owner"`
}

// Packet for transmitting between Peers
type Packet struct {
	UUID string
	Type string
	data []byte
}

// Peer for caching the Doctors office
type Peer struct {
	IP       net.IP
	Port     uint32
	LastSeen time.Time
}

// AddMedicine adds a Medicine Object to the Database
func (med *Medicine) AddMedicine(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into med(id, title, desciption, createdAt, owner) values(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// _, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
	_, err = stmt.Exec(med.UUID, med.Title, med.Desciption, med.CreatedAt, med.Owner)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

// GetMedicine ???
func (med *Medicine) GetMedicine(db *sql.DB) ([]*Medicine, error) {
	rows, err := db.Query("select id, title, desciption, createdAt, owner from med")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meds []*Medicine

	for rows.Next() {
		var id string
		var title string
		var desciption string
		var createdAt time.Time
		var owner string

		err = rows.Scan(&id, &title, &desciption, &createdAt, &owner)
		if err != nil {
			log.Fatal(err)
		}

		med := &Medicine{id, title, desciption, createdAt, owner}
		meds = append(meds, med)
		fmt.Println(id, title, desciption, createdAt, owner)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return meds, nil
}

// DeleteMedicine deletes Medicine from Database
func (med *Medicine) DeleteMedicine(db *sql.DB) {

}

// TODO: was macht 'delete'
func (med *Medicine) CreateMedicineTable(db *sql.DB) {
	sqlStmt := `
	create table med (id integer not null primary key, title text, desciption text, createdAt timestamp, owner text);
	delete from med;
	`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

// CreateDataBase creates a SqLite Database
func CreateDataBase() (*sql.DB, error) {
	// Für Debug
	os.Remove("./med.db")
	db, err := sql.Open("sqlite3", "./med.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}

// func trash() {
// 	stmt, err = db.Prepare("select name from med where id = ?")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()
// 	var name string
// 	err = stmt.QueryRow("3").Scan(&name)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(name)

// 	_, err = db.Exec("delete from med")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	_, err = db.Exec("insert into med(id, name) values(1, 'med'), (2, 'bar'), (3, 'baz')")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	rows, err = db.Query("select id, name from med")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var id int
// 		var name string
// 		err = rows.Scan(&id, &name)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(id, name)
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
