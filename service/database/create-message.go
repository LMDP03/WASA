package database

import (
	"database/sql"
	"errors"
)

var queryAddMessage = `INSERT INTO Messages (convId, senderId, msgId, text, image, responseTo, checkMark) VALUES (?, ?, ?, ?, ?, ?, "received");`

var queryGetMessageId = `SELECT MAX(msgId) AS new_id FROM Messages WHERE convId = ?;`

func (db *appdbimpl) CreateMessage(convid int, senderid int, responseto int, text string, image string) (Message, error) {

	var msg Message

	var new_id int
	var max_id = sql.NullInt64{Int64: 0, Valid: false}
	err := db.c.QueryRow(queryGetUserId).Scan(&max_id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return msg, err
	}
	if !max_id.Valid {
		new_id = 0
	} else {
		new_id = int(max_id.Int64)
	}
	new_id += 1

	_, err = db.c.Exec(queryAddMessage, convid, senderid, new_id, text, image, responseto)
	if err != nil {
		return msg, err
	}

	err = db.UpdateLastMessage(convid, new_id)
	if err != nil {
		return msg, err
	}

	msg, err = db.GetMessageById(convid, new_id)

	return msg, err
}
