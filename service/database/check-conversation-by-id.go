package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) CheckConversationById(convid int) (bool, bool, error) {
	var conv Conversation
	var lastmsg int

	err := db.c.QueryRow(queryGetConversationById, convid).Scan(&conv.Id, &conv.Name, &conv.Group, &lastmsg)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, conv.Group, nil
	}
	return true, conv.Group, err
}
