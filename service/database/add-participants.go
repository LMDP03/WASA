package database

var queryAddParticipant = `INSERT INTO Participants (convId, userId, lastRead) VALUES (?, ?, 0);`

func (db *appdbimpl) AddParticipants(convid int, names []string) error {

	for i := range names {
		user, err := db.GetUserByName(names[i])
		if err != nil {
			return err
		}

		_, err = db.c.Exec(queryAddParticipant, convid, user.Id)
		if err != nil {
			return err
		}
	}

	return nil

}
