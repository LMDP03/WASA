package database

import "time"

type User struct {
	Id   int
	Name string
}

type Preview struct {
	Id          int
	Name        string
	Group       bool
	LastMessage Message
	UserId      int
}

type Message struct {
	ConvId     int
	Sender     User
	MsgId      int
	ResponseTo Response
	Text       string
	Image      string
	Timestamp  time.Time
	Checkmark  string
	Reactions  []Reaction
	Forwarded  bool
}

type Response struct {
	Sender User
	MsgId  int
	Text   string
	Image  string
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
	Sender User
	Emoji  string
}
