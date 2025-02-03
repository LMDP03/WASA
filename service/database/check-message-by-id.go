package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) CheckMessageById(convid int, msgid int) (bool, error) {
	var msg Message
	msg.ConvId = convid
	msg.MsgId = msgid
	err := db.c.QueryRow(queryGetMessageById, convid, msgid).Scan(&msg.SenderId, &msg.Text, &msg.Image, &msg.Timestamp, &msg.ResponseTo, &msg.Checkmark)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, err
}
