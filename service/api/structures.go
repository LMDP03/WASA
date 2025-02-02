package api

import (
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/images"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
)

type User struct {
	Id    int    `json: "id"`
	Name  string `json: "name"`
	Image string `json: "image"`
}

// convert user from db structure to api structure
func (u *User) ConvertUser(user database.User) error {
	u.Id = user.Id
	u.Name = user.Name
	image, err := images.ConvertToBase64(images.SetDefaultUserImage(user.Id))
	if err != nil {
		return err
	}
	u.Image = image
	return nil
}

type Message struct{}

type Conversation struct {
	Id           int       `json: "id"`
	Name         string    `json: "name"`
	Image        string    `json: "image"`
	Group        bool      `json: "group"`
	LastMessage  int       `json: "lastmessage"`
	Participants []User    `json: "participants"`
	Messages     []Message `json: "messages"`
}

func (c *Conversation) ConvertConversation(conv database.Conversation) error {
	c.Id = conv.Id
	c.Name = conv.Name
	image, err := images.ConvertToBase64(images.SetDefaultConversationImage(conv.Id))
	if err != nil {
		return err
	}
	c.Image = image
	c.Group = conv.Group
	c.LastMessage = conv.LastMessage
	return nil
}

type Snippet struct {
	sender       User
	timestamp    time.Time
	checkmark    string
	conversation int
	id           int
	contenttype  bool
	content      string
}
