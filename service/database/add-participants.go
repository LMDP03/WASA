package database

var queryAddParticipant = `INSERT INTO Participants (convId, userId) VALUES (?, ?);`

func (db *appdbimpl) AddParticipants(convid int, names []string) ([]User, error) {

	for i := range names {
		user, err := db.GetUserByName(names[i])
		if err != nil {
			return nil, err
		}

		_, err = db.c.Exec(queryAddParticipant, convid, user.Id)
		if err != nil {
			return nil, err
		}
	}

	participants, err := db.GetParticipants(convid)

	return participants, err

}
