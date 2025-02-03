package database

var queryGetMessageById = `SELECT senderId, text, COALESCE(image, ""), timeStamp, responseTo, checkMark FROM Messages WHERE convId = ? AND msgId = ?;`

func (db *appdbimpl) GetMessageById(convId int, msgId int) (Message, error) {
	var msg Message
	msg.ConvId = convId
	msg.MsgId = msgId
	err := db.c.QueryRow(queryGetMessageById, convId, msgId).Scan(&msg.SenderId, &msg.Text, &msg.Image, &msg.Timestamp, &msg.ResponseTo, &msg.Checkmark)
	msg.Reactions, err = db.GetReactions(convId, msgId)
	return msg, err
}
