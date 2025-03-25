package database

var queryGetParticipants = `SELECT userId, lastRead FROM Participants WHERE convId = ?;`

func (db *appdbimpl) GetParticipants(convid int) ([]User, error) {

	rows, err := db.c.Query(queryGetParticipants, convid)
	if err != nil {
		return nil, err
	}

	var participants []User

	for rows.Next() {

		if rows.Err() != nil {
			return nil, err
		}

		var userid int
		var lastRead int
		err = rows.Scan(&userid, &lastRead)
		if err != nil {
			return nil, err
		}
		user, err := db.GetUserById(userid)
		if err != nil {
			return nil, err
		}

		participants = append(participants, user)
	}

	return participants, nil
}
