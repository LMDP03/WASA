package database

var queryGetMessages = `SELECT senderId, msgId, text, COALESCE(image, ""), timeStamp, responseTo, checkMark, forwarded FROM Messages WHERE convId = ? ORDER BY msgId;`

func (db *appdbimpl) GetMessages(convId int) ([]Message, error) {

	var messages []Message

	rows, err := db.c.Query(queryGetMessages, convId)
	if err != nil {
		return nil, err
	}
	defer func() { rows.Close() }()

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		var msg Message
		msg.ConvId = convId
		var senderId int
		var responseTo int
		var forwarded int
		err = rows.Scan(&senderId, &msg.MsgId, &msg.Text, &msg.Image, &msg.Timestamp, &responseTo, &msg.Checkmark, &forwarded)
		if err != nil {
			return nil, err
		}
		msg.Sender, err = db.GetUserById(senderId)
		if err != nil {
			return nil, err
		}
		msg.Reactions, err = db.GetReactions(convId, msg.MsgId)
		if err != nil {
			return nil, err
		}
		msg.ResponseTo.MsgId = responseTo
		if responseTo != 0 {
			var responseSender int
			err = db.c.QueryRow(queryGetResponse, convId, responseTo).Scan(&responseSender, &msg.ResponseTo.Text, &msg.ResponseTo.Image)
			if err != nil {
				return nil, err
			}
			msg.ResponseTo.Sender, err = db.GetUserById(responseSender)
			if err != nil {
				return nil, err
			}
		}
		if forwarded == 1 {
			msg.Forwarded = true
		} else {
			msg.Forwarded = false
		}
		messages = append(messages, msg)
	}

	return messages, nil

}
