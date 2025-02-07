package database

var queryGetReactions = `SELECT senderId, emoji FROM Reactions WHERE convId = ? AND msgId = ?;`

func (db *appdbimpl) GetReactions(convid int, msgid int) ([]Reaction, error) {

	var reactions []Reaction
	rows, err := db.c.Query(queryGetReactions, convid, msgid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		var reac Reaction
		var senderid int
		err = rows.Scan(&senderid, &reac.Emoji)
		if err != nil {
			return nil, err
		}
		reac.Sender, err = db.GetUserById(senderid)
		if err != nil {
			return nil, err
		}
		reactions = append(reactions, reac)
	}
	return reactions, nil
}
