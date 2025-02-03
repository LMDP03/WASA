package database

var queryGetReaction = `SELECT emoji FROM Reactions WHERE convId = ? AND msgId = ? AND senderId = ?;`

func (db *appdbimpl) GetReactionById(convid int, msgid int, senderid int) (Reaction, error) {

	var reac Reaction
	err := db.c.QueryRow(queryGetReactions, convid, msgid, senderid).Scan(&reac.Emoji)
	if err != nil {
		return reac, err
	}
	u, err := db.GetUserById(senderid)
	if err != nil {
		return reac, err
	}
	reac.SenderName = u.Name
	return reac, nil
}
