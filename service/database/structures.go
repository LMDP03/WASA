package database

import "time"

type User struct {
	Id   int
	Name string
}

type Preview struct {
	Id          int
	Name        string
	UserId      int
	Group       bool
	LastMessage Message
}

type Message struct {
	ConvId     int
	SenderId   int
	MsgId      int
	ResponseTo int
	Text       string
	Image      string
	Timestamp  time.Time
	Checkmark  string
	Reactions  []Reaction
}

type Conversation struct {
	Id           int
	Name         string
	UserId       int
	Group        bool
	Participants []User
	Messages     []Message
}

type Reaction struct {
	SenderName string
	Emoji      string
}
