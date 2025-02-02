package database

import (
	"fmt"
	"io"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/images"
)

var queryAddConversation = `INSERT INTO Conversations (name, group, last_message) VALUES (?, ?, 0);`

var queryGetConversationId = `SELECT MAX(id) AS new_id FROM Conversations;`

func (db *appdbimpl) CreateConversation(name string, group bool) (Conversation, error) {

	var conv Conversation

	_, err := db.c.Exec(queryAddConversation, name, group)
	if err != nil {
		return conv, err
	}

	var new_id int
	err = db.c.QueryRow(queryGetConversationId).Scan(&new_id)
	if err != nil {
		return conv, err
	}

	path := fmt.Sprintf("./storage/conversations/%d", new_id)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return conv, err
	}

	var source *os.File

	if group {
		source, err := os.Open("./storage/default_image.jpg") // Open the img file
		if err != nil {
			return conv, err
		}
		defer source.Close()
	} else {
		user, err := db.GetUserByName(name)
		if err != nil {
			return conv, err
		}
		source, err := os.Open(images.SetDefaultUserImage(user.Id)) // Open the img file
		if err != nil {
			return conv, err
		}
		defer source.Close()
	}

	destination, err := os.Create(images.SetDefaultConversationImage(new_id)) // Create the path where the photo will be saved
	if err != nil {
		return conv, err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return conv, err
	}

	conv.Id = new_id
	conv.Name = name
	conv.Group = group
	conv.LastMessage = 0

	return conv, nil
}
