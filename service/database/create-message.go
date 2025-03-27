package database

import (
	"database/sql"
	"errors"
)

var queryAddMessage = `INSERT INTO Messages (convId, senderId, msgId, text, image, responseTo, checkMark, accessCount, forwarded) VALUES (?, ?, ?, ?, ?, ?, "received", 1, ?);`
var queryGetLastMessageId = `SELECT last_message FROM Conversations WHERE id = ?;`

func (db *appdbimpl) CreateMessage(convid int, senderid int, responseto int, text string, image string, forwarded bool) (Message, error) {

	var msg Message
	var forwardedInt int
	if forwarded {
		forwardedInt = 1
	} else {
		forwardedInt = 0
	}

	var max_id int
	err := db.c.QueryRow(queryGetLastMessageId, convid).Scan(&max_id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return msg, err
	}
	max_id += 1

	_, err = db.c.Exec(queryAddMessage, convid, senderid, max_id, text, image, responseto, forwardedInt)
	if err != nil {
		return msg, err
	}
	err = db.UpdateLastMessage(convid)
	if err != nil {
		return msg, err
	}
	_, err = db.c.Exec(queryUpdateLastRead, max_id, convid, senderid)
	if err != nil {
		return msg, err
	}

	msg, err = db.GetMessageById(convid, max_id)

	return msg, err
}
