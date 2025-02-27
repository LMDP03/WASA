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
	AddParticipants(convid int, names []string) error
	CheckConversationById(convid int) (bool, bool, error)
	CheckMessageById(convid int, msgid int) (bool, error)
	CheckReactionById(convid int, msgid int, senderid int) (bool, error)
	CheckUserById(userid int) (bool, error)
	CheckUserByName(username string) (bool, error)
	CreateConversation(name string, group bool, otherid int, participants []string) (Conversation, error)
	CreateMessage(convid int, senderid int, responseto int, text string, image string) (Message, error)
	CreateReaction(convId int, senderId int, msgId int, emoji string) (Reaction, error)
	CreateUser(username string) (User, error)
	DeleteMessage(convid int, msgid int, senderid int) error
	DeleteParticipant(convid int, userid int) error
	DeleteReaction(convid int, msgid int, senderid int) error
	GetConversationById(convid int, userid int) (Conversation, error)
	GetConversations(userid int, searchname string) ([]Preview, error)
	GetMessageById(convId int, msgId int) (Message, error)
	GetMessages(convId int) ([]Message, error)
	GetOtherParticipant(convid int, userid int) (User, error)
	GetParticipants(convid int) ([]User, error)
	GetReactionById(convid int, msgid int, senderid int) (Reaction, error)
	GetReactions(convid int, msgid int) ([]Reaction, error)
	GetUserById(userid int) (User, error)
	GetUserByName(username string) (User, error)
	GetUsersByName(searchname string, userid int) ([]User, error)
	IsParticipant(convid int, userid int) (bool, error)
	SetGroupNameById(convid int, name string) error
	SetUserNameById(username string, userid int) error
	UpdateReaction(convId int, senderId int, msgId int, emoji string) (Reaction, error)
	UpdateLastMessage(convid int, msgid int) error
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
	var tables uint8
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
