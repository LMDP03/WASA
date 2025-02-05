package database

import (
	"fmt"
	"io"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/images"
)

var queryAddConversation = `INSERT INTO Conversations (id, name, groupFlag, last_message) VALUES (?, ?, ?, 0);`

var queryGetConversationId = `SELECT MAX(id) AS new_id FROM Conversations;`

func (db *appdbimpl) CreateConversation(name string, group bool, otherid int, participants []string) (Conversation, error) {

	var conv Conversation

	var new_id int

	var flag int

	if group {

		flag = 1

		err := db.c.QueryRow(queryGetConversationId).Scan(&new_id)
		if err != nil {
			return conv, err
		}

		new_id += 1

		_, err = db.c.Exec(queryAddConversation, new_id, name, flag)
		if err != nil {
			return conv, err
		}

		err = db.AddParticipants(new_id, participants)
		if err != nil {
			return conv, err
		}

		path := fmt.Sprintf("./storage/groups/%d", new_id)

		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return conv, err
		}
		source, err := os.Open("./storage/default_image.jpg") // Open the img file
		if err != nil {
			return conv, err
		}
		defer source.Close()
		destination, err := os.Create(images.SetDefaultGroupImage(new_id)) // Create the path where the photo will be saved
		if err != nil {
			return conv, err
		}
		defer destination.Close()
		_, err = io.Copy(destination, source)
		if err != nil {
			return conv, err
		}

		conv, err = db.GetConversationById(new_id, 0)

	} else {
		flag = 0

		err := db.c.QueryRow(queryGetConversationId).Scan(&new_id)
		if err != nil {
			return conv, err
		}

		new_id += 1

		_, err = db.c.Exec(queryAddConversation, new_id, "", flag)
		if err != nil {
			return conv, err
		}

		err = db.AddParticipants(new_id, participants)
		if err != nil {
			return conv, err
		}

		conv, err = db.GetConversationById(new_id, otherid)

	}

	return conv, nil
}
