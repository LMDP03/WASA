package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) CheckUserByName(username string) (bool, error) {
	var user User
	err := db.c.QueryRow(queryGetUserByName, username).Scan(&user.Id, &user.Name)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, err
}
