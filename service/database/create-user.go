package database

import (
	"fmt"
	"io"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/images"
)

var queryAddUser = `INSERT INTO Users (name) VALUES ?;`

func (db *appdbimpl) CreateUser(username string) (User, error) {

	var user User
	_, err := db.c.Exec(queryAddUser, username)
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
