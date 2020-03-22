package db

import (
	"database/sql"
	"os"
)

var CurrentDb *sql.DB

// RemoveDataBase deletes the data base med
func RemoveDataBase() {
	os.Remove("./med.db")
}

// CreateDataBase creates a SqLite3 Database named med.db
func CreateDataBase() (*sql.DB, error) {
	var err error
	CurrentDb, err = sql.Open("sqlite3", "./med.db")
	if err != nil {
		return nil, err
	}

	CreatePeerTable(CurrentDb)
	CreateMedicineTable(CurrentDb)

	return CurrentDb, nil
}
