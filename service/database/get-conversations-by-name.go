package database

import "strings"

var queryGetConversationsByName = `SELECT id, name, group, last_message FROM Conversations, Participants WHERE id = convId AND userId = ? ORDER BY last_message DESC;`

func (db *appdbimpl) GetConversationsbyName(searchname string, userid int) ([]Preview, error) {

	var previews []Preview

	rows, err := db.c.Query(queryGetConversationsByName, userid)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		var prev Preview
		var lastmsg int
		var flag int
		if err := rows.Scan(&prev.Id, &prev.Name, &flag, &lastmsg); err != nil {
			return nil, err
		}
		if flag == 0 {
			prev.Group = false
		} else {
			prev.Group = true
		}
		if lastmsg != 0 {
			message, err := db.GetMessageById(prev.Id, lastmsg)
			if err != nil {
				return nil, err
			}
			prev.LastMessage = message
		}
		if !prev.Group {
			other, err := db.GetOtherParticipant(prev.Id, userid)
			if err != nil {
				return nil, err
			}
			prev.Name = other.Name
			prev.UserId = other.Id
		}
		if strings.Contains(prev.Name, searchname) {
			previews = append(previews, prev)
		}
	}
	defer func() { err = rows.Close() }()
	return previews, err
}
