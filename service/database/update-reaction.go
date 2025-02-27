package database

var queryUpdateReaction = `UPDATE Reactions SET emoji = ? WHERE convId = ? AND senderId = ? AND msgId = ?;`

func (db *appdbimpl) UpdateReaction(convId int, senderId int, msgId int, emoji string) (Reaction, error) {

	var reac Reaction
	_, err := db.c.Exec(queryUpdateReaction, emoji, convId, senderId, msgId)
	if err != nil {
		return reac, err
	}
	reac, err = db.GetReactionById(convId, msgId, senderId)
	if err != nil {
		return reac, err
	}
	return reac, nil
}
