package database

var queryAddMessage = `INSERT INTO Messages (convId, senderId, msgId, text, image, responseTo, checkMark) VALUES (?, ?, ?, ?, ?, ?, "received");`

var queryGetMessageId = `SELECT MAX(id) AS new_id FROM Messages;`

func (db *appdbimpl) CreateMessage(convid int, senderid int, responseto int, text string, image string) (Message, error) {

	var msg Message

	var new_id int

	err := db.c.QueryRow(queryGetMessageId).Scan(&new_id)
	if err != nil {
		return msg, err
	}

	new_id += 1

	_, err = db.c.Exec(queryAddMessage, convid, senderid, new_id, text, image, responseto)
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
