package database

var queryAddMessage = `INSERT INTO Messages (convId, senderId, text, image, responseTo, checkMark) VALUES (?, ?, ?, ?, ?, "sent");`

var queryGetMessageId = `SELECT MAX(id) AS new_id FROM Messages;`

func (db *appdbimpl) CreateMessage(convid int, senderid int, responseto int, text string, image string) (Message, error) {

	var msg Message
	_, err := db.c.Exec(queryAddMessage, convid, senderid, text, image, responseto)
	if err != nil {
		return msg, err
	}

	var new_id int
	err = db.c.QueryRow(queryGetConversationId).Scan(&new_id)
	if err != nil {
		return msg, err
	}

	err = db.UpdateLastMessage(convid, new_id)
	if err != nil {
		return msg, err
	}

	msg, err = db.GetMessageById(convid, new_id)

	return msg, err
}
