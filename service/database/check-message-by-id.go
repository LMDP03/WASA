package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) CheckMessageById(convid int, msgid int) (bool, error) {
	var msg Message
	msg.ConvId = convid
	msg.MsgId = msgid
	var senderId int
	var responseTo int
	err := db.c.QueryRow(queryGetMessageById, convid, msgid).Scan(&senderId, &msg.Text, &msg.Image, &msg.Timestamp, &responseTo, &msg.Checkmark)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, err
}
