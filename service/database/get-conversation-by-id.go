package database

var queryGetConversationById = `SELECT id, name, groupFlag, last_message FROM Conversations WHERE id = ?;`

func (db *appdbimpl) GetConversationById(convid int, userid int) (Conversation, error) {

	var conv Conversation
	var lastmsg int
	var flag int

	err := db.c.QueryRow(queryGetConversationById, convid).Scan(&conv.Id, &conv.Name, &flag, &lastmsg)
	if err != nil {
		return conv, err
	}

	conv.Group = true

	participants, err := db.GetParticipants(convid)
	if err != nil {
		return conv, err
	}
	conv.Participants = participants

	if lastmsg != 0 {
		messages, err := db.GetMessages(convid)
		if err != nil {
			return conv, err
		}
		conv.Messages = messages
	}

	if flag == 0 {
		conv.Group = false
		other, err := db.GetOtherParticipant(convid, userid)
		if err != nil {
			return conv, err
		}
		conv.Name = other.Name
		conv.UserId = other.Id
	}

	return conv, nil
}
