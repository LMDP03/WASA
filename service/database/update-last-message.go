package database

var queryUpdateLastMessage = `UPDATE Conversations SET last_message = (SELECT MAX(msgId) FROM Messages WHERE convId = ?) WHERE id = ?;`

func (db *appdbimpl) UpdateLastMessage(convid int) error {

	_, err := db.c.Exec(queryUpdateLastMessage, convid, convid)
	if err != nil {
		return err
	}
	return nil
}
