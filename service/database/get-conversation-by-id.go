package database

var queryGetConversationById = `SELECT id, name, group, last_message FROM Conversations WHERE id = ?;`

func (db *appdbimpl) GetConversationById(convid int, userid int) (Conversation, error) {

	var conv Conversation
	var lastmsg int

	err := db.c.QueryRow(queryGetConversationById, convid).Scan(&conv.Id, &conv.Name, &conv.Group, &lastmsg)
	if err != nil {
		return conv, err
	}

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

	if !conv.Group {
		other, err := db.GetOtherParticipant(convid, userid)
		if err != nil {
			return conv, err
		}
		conv.Name = other.Name
		conv.UserId = other.Id
	}

	return conv, nil
}
