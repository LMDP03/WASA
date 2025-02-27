package database

var queryGetOtherParticipant = `SELECT userId FROM Participants WHERE convId = ? AND NOT userId = ?;`

func (db *appdbimpl) GetOtherParticipant(convid int, userid int) (User, error) {

	// give your user id to find the other user of a private chat

	var other User

	var otherid int

	err := db.c.QueryRow(queryGetOtherParticipant, convid, userid).Scan(&otherid)
	if err != nil {
		return other, err
	}

	other, err = db.GetUserById(otherid)
	if err != nil {
		return other, err
	}
	return other, err
}
