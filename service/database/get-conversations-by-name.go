package database

var queryGetConversationsByName = `SELECT id, name, group, last_message FROM Conversations, Participants WHERE id = convId AND name LIKE '%?%' AND  userId = ? ORDER BY last_message DESC;`

func (db *appdbimpl) GetConversationsbyName(searchname string, userid int) ([]Conversation, error) {

	var conversations []Conversation

	rows, err := db.c.Query(queryGetConversationsByName, searchname, userid)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		var conv Conversation
		if err := rows.Scan(&conv.Id, &conv.Name, &conv.Group, &conv.LastMessage); err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}
	defer func() { err = rows.Close() }()
	return conversations, err
}
