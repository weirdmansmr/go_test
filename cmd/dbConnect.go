package cmd

import (
	"database/sql"
	"fmt"
)

func ConnectDB() (*sql.DB, error) {
	connStr := "postgres://postgres:Saykhanov01@localhost/bank_test?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging the database: %v", err)
	}
	return db, nil
}