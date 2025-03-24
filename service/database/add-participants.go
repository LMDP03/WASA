package database

var queryAddParticipant = `INSERT INTO Participants (convId, userId, lastRead) VALUES (?, ?, (SELECT last_message FROM Conversations WHERE id = ?));`

func (db *appdbimpl) AddParticipants(convid int, names []string) error {

	for i := range names {
		user, err := db.GetUserByName(names[i])
		if err != nil {
			return err
		}

		_, err = db.c.Exec(queryAddParticipant, convid, user.Id, convid)
		if err != nil {
			return err
		}
	}

	return nil

}
