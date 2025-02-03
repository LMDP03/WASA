package database

var queryUpdateLastMessage = `UPDATE Conversations SET last_message = ? WHERE id = ?;`

func (db *appdbimpl) UpdateLastMessage(convid int, msgid int) error {

	_, err := db.c.Exec(queryUpdateLastMessage, msgid, convid)
	if err != nil {
		return err
	}
	return nil
}
