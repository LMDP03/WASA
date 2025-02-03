package database

var queryDeleteMessage = `DELETE FROM Messages WHERE msgId = ? AND convId = ?;`

func (db *appdbimpl) DeleteMessage(convid int, msgid int) error {

	_, err := db.c.Exec(queryDeleteMessage, convid, msgid)
	return err
}
