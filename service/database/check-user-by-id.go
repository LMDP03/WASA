package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) CheckUserById(userid int) (bool, error) {
	var user User
	err := db.c.QueryRow(queryGetUserById, userid).Scan(&user.Id, &user.Name)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, err
}
