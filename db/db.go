package db

import (
	"database/sql"
	"os"
)

// RemoveDataBase deletes the data base med
func RemoveDataBase() {
	os.Remove("./med.db")
}

// CreateDataBase creates a SqLite3 Database named med.db
func CreateDataBase() (*sql.DB, error) {
	// Für Debug
	RemoveDataBase()

	db, err := sql.Open("sqlite3", "./med.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}
