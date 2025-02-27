package database

var queryDeleteParticipant = `DELETE FROM Participants WHERE convId = ? AND userId = ?;`

var queryDeleteConversation = `DELETE FROM Conversations WHERE id = ? ;`

func (db *appdbimpl) DeleteParticipant(convid int, userid int) error {

	_, err := db.c.Exec(queryDeleteParticipant, convid, userid)
	if err != nil {
		return err
	}

	remaining_memebers, err := db.GetParticipants(convid)
	if err != nil {
		return err
	}
	if len(remaining_memebers) == 0 {
		_, err = db.c.Exec(queryDeleteConversation, convid)
	}
	return err
}
