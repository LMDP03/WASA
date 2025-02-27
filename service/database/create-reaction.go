package database

var queryCreateReaction = `INSERT INTO Reactions (convId, senderId, msgId, emoji) VALUES (?, ?, ?, ?);`

func (db *appdbimpl) CreateReaction(convId int, senderId int, msgId int, emoji string) (Reaction, error) {

	var reac Reaction
	_, err := db.c.Exec(queryCreateReaction, convId, senderId, msgId, emoji)
	if err != nil {
		return reac, err
	}
	reac, err = db.GetReactionById(convId, msgId, senderId)
	if err != nil {
		return reac, err
	}
	return reac, nil
}
