package database

import (
	"database/sql"
	"errors"
)

var queryDeleteMessage = `DELETE FROM Messages WHERE msgId = ? AND convId = ?;`
var queryUpdateLaterResponses = `UPDATE Messages SET responseTo = 0 WHERE convId = ? AND msgId > ? AND responseTo = ?;`
var queryDeleteReactions = `DELETE FROM Reactions WHERE convId = ? AND msgId = ?;`
var queryUpdateLastReadAll = `UPDATE Participants SET lastRead = (SELECT last_message FROM Conversations WHERE id = ?) WHERE convId = ?;`

func (db *appdbimpl) DeleteMessage(convid int, msgid int) error {

	var max_id int
	err := db.c.QueryRow(queryGetLastMessageId, convid).Scan(&max_id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = db.c.Exec(queryDeleteMessage, msgid, convid)
	if err != nil {
		return err
	}

	if msgid < max_id {
		_, err = db.c.Exec(queryUpdateLaterResponses, convid, msgid, msgid)
		if err != nil {
			return err
		}
	} else {
		err = db.UpdateLastMessage(convid)
		if err != nil {
			return err
		}

		_, err = db.c.Exec(queryUpdateLastReadAll, convid, convid)
		if err != nil {
			return err
		}

	}

	_, err = db.c.Exec(queryDeleteReactions, convid, msgid)
	if err != nil {
		return err
	}
	return nil
}
