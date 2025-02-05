package database

import (
	"fmt"
	"io"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/images"
)

var queryAddUser = `INSERT INTO Users (id, name) VALUES (?, ?);`

var queryGetUserId = `SELECT MAX(id) AS new_id FROM Messages;`

func (db *appdbimpl) CreateUser(username string) (User, error) {

	var user User

	var new_id int
	err := db.c.QueryRow(queryGetUserId).Scan(&new_id)
	if err != nil {
		return user, err
	}

	_, err = db.c.Exec(queryAddUser, new_id+1, username)
	if err != nil {
		return user, err
	}
	user, err = db.GetUserByName(username)
	if err != nil {
		return user, err
	}

	// Creation of user folder
	path := fmt.Sprintf("./storage/users/%d/conversations", user.Id)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return user, err
	}
	// Set default photo profile
	source, err := os.Open("./storage/default_image.jpg") // Open the img file
	if err != nil {
		return user, err
	}
	defer source.Close()

	destination, err := os.Create(images.SetDefaultUserImage(user.Id)) // Create the path where the photo will be saved
	if err != nil {
		return user, err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source) // Copy the photo in the user folder
	if err != nil {
		return user, err
	}

	return user, nil
}
