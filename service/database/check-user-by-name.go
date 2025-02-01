package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) CheckUserByName(username string) (bool, error) {
	var name string
	err := db.c.QueryRow(queryGetUserByName, username).Scan(&name)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, err
}
