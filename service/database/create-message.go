package database

import (
	"database/sql"
	"errors"
)

var queryAddMessage = `INSERT INTO Messages (convId, senderId, msgId, text, image, responseTo, checkMark) VALUES (?, ?, ?, ?, ?, ?, "received");`

var queryGetLastMessageId = `SELECT last_message FROM Conversations WHERE id = ?;`

func (db *appdbimpl) CreateMessage(convid int, senderid int, responseto int, text string, image string) (Message, error) {

	var msg Message

	var max_id int
	err := db.c.QueryRow(queryGetLastMessageId, convid).Scan(&max_id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return msg, err
	}
	max_id += 1

	_, err = db.c.Exec(queryAddMessage, convid, senderid, max_id, text, image, responseto)
	if err != nil {
		return msg, err
	}
	err = db.UpdateLastMessage(convid)
	if err != nil {
		return msg, err
	}
	msg, err = db.GetMessageById(convid, max_id)

	return msg, err
}
