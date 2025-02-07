package database

var queryGetReaction = `SELECT emoji FROM Reactions WHERE convId = ? AND msgId = ? AND senderId = ?;`

func (db *appdbimpl) GetReactionById(convid int, msgid int, senderid int) (Reaction, error) {

	var reac Reaction
	err := db.c.QueryRow(queryGetReaction, convid, msgid, senderid).Scan(&reac.Emoji)
	if err != nil {
		return reac, err
	}
	reac.Sender, err = db.GetUserById(senderid)
	if err != nil {
		return reac, err
	}
	return reac, nil
}
