package database

var queryDeleteMessage = `DELETE FROM Messages WHERE msgId = ? AND convId = ? AND senderId = ?;`

func (db *appdbimpl) DeleteMessage(convid int, msgid int, senderid int) error {

	_, err := db.c.Exec(queryDeleteMessage, convid, msgid, senderid)
	if err != nil {
		return err
	}

	err = db.UpdateLastMessage(convid, msgid-1)

	return err
}
