package database

var queryGetMessages = `SELECT senderId, msgId, text, COALESCE(image, ""), timeStamp, responseTo, checkMark FROM Messages WHERE convId = ? ORDER BY msgId DESC;`

func (db *appdbimpl) GetMessages(convId int) ([]Message, error) {

	var messages []Message

	rows, err := db.c.Query(queryGetMessages, convId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		var msg Message
		msg.ConvId = convId
		err = rows.Scan(&msg.SenderId, &msg.MsgId, &msg.Text, &msg.Image, &msg.Timestamp, &msg.ResponseTo, &msg.Checkmark)
		if err != nil {
			return nil, err
		}
		reactions, err := db.GetReactions(convId, msg.MsgId)
		if err != nil {
			return nil, err
		}
		msg.Reactions = reactions
		messages = append(messages, msg)
	}

	return messages, nil

}
