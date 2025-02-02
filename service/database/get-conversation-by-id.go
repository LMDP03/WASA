package database

var queryGetConversationById = `SELECT id, name, group, last_message FROM Conversations WHERE id = ?;`

func (db *appdbimpl) GetConversationById(id int) (Conversation, error) {

	var conv Conversation
	err := db.c.QueryRow(queryGetConversationById, id).Scan(&conv.Id, &conv.Name, &conv.Group, &conv.LastMessage)
	if err != nil {
		return Conversation{}, err
	}
	return conv, nil
}
