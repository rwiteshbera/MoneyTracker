package config

import (
	"database/sql"
	"errors"
)

func ConnectDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./user.db")
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return database, nil
}
