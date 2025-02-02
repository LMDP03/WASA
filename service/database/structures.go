package database

type User struct {
	Id   int    `json: "id"`
	Name string `json: "name"`
}

type Conversation struct {
	Id          int    `json: "id"`
	Name        string `json: "name"`
	Group       bool   `json: "group"`
	LastMessage int    `json: "lastmessage"`
}

type Message struct {
	convId   int
	senderId int
	msgId    int
	normal   bool
	text     string
}
