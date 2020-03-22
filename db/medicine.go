package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Needs to be like this
)

// Medicine is the stuff that makes you healthy
type Medicine struct {
	UUID        string `json:"uuid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Owner       Peer   `json:"owner"`
	Amount      int    `json:"amount"`
	Pzn         int    `json:"pzn"`
}

// Add adds a Medicine Object to the Database
func (med *Medicine) Add(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}


	stmt, err := tx.Prepare("insert into med(id, title, description, amount, pzn) values(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(med.UUID, med.Title, med.Description, med.Amount, med.Pzn)
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

// get wrapper to find a set of results
func get(db *sql.DB, query string, searchStr string) ([]*Medicine, error) {
	// SELECT * FROM table WHERE instr(title, searchStr) > 0 OR instr(description, searchStr) > 0 OR searchStr == CAST(pzn as text)
	// SELECT * FROM med WHERE instr(title, '3') > 0 OR instr(description, '3') > 0 OR '3' == CAST(pzn as text);
	// PZN != ID
	rows, err := db.Query(query, searchStr, searchStr, searchStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meds []*Medicine

	for rows.Next() {
		var id string
		var title string
		var description string
		var amount int
		var pzn int

		err = rows.Scan(&id, &title, &description, &amount, &pzn)
		if err != nil {
			log.Fatal(err)
		}

		med := &Medicine{id, title, description, Peer{}, amount, pzn}
		meds = append(meds, med)
		fmt.Println(id, title, description, amount, pzn)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return meds, nil
}

// Get a specific Medicine from the Database based on the search String or the pzn
func Get(db *sql.DB, searchStr string) ([]*Medicine, error) {
	return get(db, "select * from med where instr(title, ?) > 0 OR instr(description, ?) > 0 or ? == cast(pzn as text)", searchStr)
}

// GetAll all the rows of the Database
func GetAll(db *sql.DB) ([]*Medicine, error) {
	return get(db, "select id, title, description, amount, pzn from med", "")
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
	create table med (id integer not null primary key, title text, description text, amount integer, pzn integer);`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
