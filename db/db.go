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
	Amount     int       `json:"amount"`
	Pzn        int       `json:"pzn"`
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

// Add adds a Medicine Object to the Database
func (med *Medicine) Add(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into med(id, title, desciption, createdAt, owner, amount, pzn) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(med.UUID, med.Title, med.Desciption, med.CreatedAt, med.Owner, med.Amount, med.Pzn)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

// Delete Medicine from Database
func (med *Medicine) Delete(db *sql.DB) {
	_, err := db.Exec("delete from med where id=?", med.UUID)
	if err != nil {
		log.Fatal(err)
	}
}

// Update the Medicine object in the Database
func (med *Medicine) Update(db *sql.DB) {
	med.Delete(db)
	med.Add(db)
}

// func get(db *sql.DB, searchStr)

func Get(db *sql.DB, searchStr string) ([]*Medicine, error) {
	// SELECT * FROM table WHERE instr(title, searchStr) > 0 OR instr(description, searchStr) > 0 OR searchStr == CAST(pzn as text)
	// PZN != ID
	rows, err := db.Query("select * from med where instr(title, ?) > 0 OR instr(description, ?) > 0 or ? == cast(id as text)", searchStr)
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
		var amount int
		var pzn int

		err = rows.Scan(&id, &title, &desciption, &createdAt, &owner, &amount, &pzn)
		if err != nil {
			log.Fatal(err)
		}

		med := &Medicine{id, title, desciption, createdAt, owner, amount, pzn}
		meds = append(meds, med)
		fmt.Println(id, title, desciption, createdAt, owner, amount, pzn)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return meds, nil
}

// GetAll all the rows of the Database
func GetAll(db *sql.DB) ([]*Medicine, error) {
	rows, err := db.Query("select id, title, desciption, createdAt, owner, amount, pzn from med")
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
		var amount int
		var pzn int

		err = rows.Scan(&id, &title, &desciption, &createdAt, &owner, &amount, &pzn)
		if err != nil {
			log.Fatal(err)
		}

		med := &Medicine{id, title, desciption, createdAt, owner, amount, pzn}
		meds = append(meds, med)
		fmt.Println(id, title, desciption, createdAt, owner, amount, pzn)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return meds, nil
}

// DeleteMedicineTable deletes Medicine from Database
func DeleteMedicineTable(db *sql.DB) {
	_, err := db.Exec("delete from med")
	if err != nil {
		log.Fatal(err)
	}
}

// CreateMedicineTable creates the Medicine Table
func CreateMedicineTable(db *sql.DB) {
	sqlStmt := `
	create table med (id integer not null primary key, title text, desciption text, createdAt timestamp, owner text, amount integer, pzn integer);`

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
