package database

import (
	"database/sql"
	"errors"
)

var queryGetParticipant = `SELECT convId, userId FROM  Participants WHERE convId = ? AND userId = ?;`

func (db *appdbimpl) IsParticipant(convid int, userid int) (bool, error) {
	type Participant struct {
		convId int
		userId int
	}
	var part Participant
	err := db.c.QueryRow(queryGetParticipant, convid, userid).Scan(&part.convId, &part.userId)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, err
}
