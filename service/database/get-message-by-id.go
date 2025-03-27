package database

var queryGetMessageById = `SELECT senderId, text, COALESCE(image, ""), timeStamp, responseTo, checkMark, forwarded FROM Messages WHERE convId = ? AND msgId = ?;`
var queryGetResponse = `SELECT senderId, text, COALESCE(image, "") from Messages WHERE convId = ? AND msgId = ?;`

func (db *appdbimpl) GetMessageById(convId int, msgId int) (Message, error) {
	var msg Message
	msg.ConvId = convId
	msg.MsgId = msgId
	var senderId int
	var responseTo int
	var forwarded int
	err := db.c.QueryRow(queryGetMessageById, convId, msgId).Scan(&senderId, &msg.Text, &msg.Image, &msg.Timestamp, &responseTo, &msg.Checkmark, &forwarded)
	if err != nil {
		return msg, err
	}
	msg.Sender, err = db.GetUserById(senderId)
	if err != nil {
		return msg, err
	}
	msg.ResponseTo.MsgId = responseTo
	if responseTo != 0 {
		var responseSender int
		err = db.c.QueryRow(queryGetResponse, convId, responseTo).Scan(&responseSender, &msg.ResponseTo.Text, &msg.ResponseTo.Image)
		if err != nil {
			return msg, err
		}
		msg.ResponseTo.Sender, err = db.GetUserById(responseSender)
		if err != nil {
			return msg, err
		}
	}
	if forwarded == 1 {
		msg.Forwarded = true
	} else {
		msg.Forwarded = false
	}
	msg.Reactions, err = db.GetReactions(convId, msgId)
	return msg, err
}
