package database

var queryGetConversationById = `SELECT id, name, groupFlag, last_message FROM Conversations WHERE id = ?;`
var queryGetLastSeenMessage = `SELECT lastRead FROM Participants WHERE convId = ? AND userId = ?;`
var queryUpdateMessageCounters = `UPDATE Messages SET accessCount = accessCount + 1 WHERE convId = ? AND msgId > ?;`
var queryUpdateCheckmarks = `UPDATE Messages SET checkMark = "read" WHERE convId = ? AND accessCount >= ?;`
var queryUpdateLastRead = `UPDATE Participants SET lastRead = ? WHERE convId = ? AND userId = ?;`

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

	var lastRead int
	err = db.c.QueryRow(queryGetLastSeenMessage, convid, userid).Scan(&lastRead)
	if err != nil {
		return conv, err
	}
	if lastRead < lastmsg {
		_, err = db.c.Exec(queryUpdateMessageCounters, convid, lastRead)
		if err != nil {
			return conv, err
		}
		_, err = db.c.Exec(queryUpdateCheckmarks, convid, len(participants))
		if err != nil {
			return conv, err
		}
		_, err = db.c.Exec(queryUpdateLastRead, lastmsg, convid, userid)
		if err != nil {
			return conv, err
		}
	}

	if lastmsg != 0 {
		conv.Messages, err = db.GetMessages(convid)
		if err != nil {
			return conv, err
		}
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
