package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"

	"wasatext/service/images"
)

var queryAddUser = `INSERT INTO Users (id, name) VALUES (?, ?);`

var queryGetUserId = `SELECT MAX(id) AS new_id FROM Users;`

func (db *appdbimpl) CreateUser(username string) (User, error) {

	var user User
	var new_id int
	var max_id = sql.NullInt64{Int64: 0, Valid: false}
	err := db.c.QueryRow(queryGetUserId).Scan(&max_id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return user, err
	}
	if !max_id.Valid {
		new_id = 1
	} else {
		new_id = int(max_id.Int64) + 1
	}

	_, err = db.c.Exec(queryAddUser, new_id, username)
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
	fmt.Println(user.Id, new_id)
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
