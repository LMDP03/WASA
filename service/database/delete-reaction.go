package database

var queryDeleteReaction = `DELETE FROM Reactions WHERE convId = ? AND senderId = ? AND msgId = ?;`

func (db *appdbimpl) DeleteReaction(convid int, msgid int, senderid int) error {

	_, err := db.c.Exec(queryDeleteMessage, convid, senderid, msgid)
	return err
}
