package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) CheckReactionById(convid int, msgid int, senderid int) (bool, error) {
	var reac Reaction
	err := db.c.QueryRow(queryGetReactions, convid, msgid, senderid).Scan(&reac.Emoji)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, err
}
