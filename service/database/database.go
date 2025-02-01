/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	CheckUserByName(username string) (bool, error)
	CreateUser(username string) (User, error)
	GetUserById(userid int32) (User, error)
	GetUserByName(username string) (User, error)
	GetUsersByName(username string, userid int) ([]User, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tables int
	err := db.QueryRow(`SELECT COUNT(name) FROM sqlite_master WHERE type='table'`).Scan(&tables)
	if err != nil {
		return nil, fmt.Errorf("error checking if database is empty: %w", err)
	}

	if tables != 5 {
		_, err = db.Exec(sql_USERS)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure Users: %w", err)
		}

		_, err = db.Exec(sql_CONVERSATIONS)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure Conversations: %w", err)
		}

		_, err = db.Exec(sql_PARTICIPANTS)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure Participants: %w", err)
		}

		_, err = db.Exec(sql_MESSAGES)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure Messages: %w", err)
		}

		_, err = db.Exec(sql_REACTIONS)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure Reactions: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
