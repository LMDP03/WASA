package database

import (
	"database/sql"
	"errors"
)

var queryDeleteMessage = `DELETE FROM Messages WHERE msgId = ? AND convId = ? AND senderId = ?;`
var queryUpdateLaterResponses = `UPDATE Messages SET responseTo = 0 WHERE convId = ? AND msgId > ? AND responseTo = ?;`
var queryDeleteReactions = `DELETE FROM Reactions WHERE convId = ? AND msgId = ?;`

func (db *appdbimpl) DeleteMessage(convid int, msgid int, senderid int) error {

	var max_id int
	err := db.c.QueryRow(queryGetLastMessageId, convid).Scan(&max_id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if msgid < max_id {
		_, err = db.c.Exec(queryUpdateLaterResponses, convid, msgid, msgid)
		if err != nil {
			return err
		}
	} else {
		err = db.UpdateLastMessage(convid, max_id-1)
		if err != nil {
			return err
		}
	}

	_, err = db.c.Exec(queryDeleteMessage, convid, msgid, senderid)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(queryDeleteReactions, convid, msgid)
	if err != nil {
		return err
	}

	return nil
}
