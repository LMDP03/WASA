package api

import (
	"time"
	"unicode/utf8"

	"wasatext/service/database"
	"wasatext/service/images"
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

type Reaction struct {
	SenderName string `json: "sender"`
	Emoji      string `json: "emoji"`
}

func (r *Reaction) ConvertReaction(reac database.Reaction) {
	r.SenderName = reac.SenderName
	r.Emoji = reac.Emoji

}

func CheckEmoji(emoji string) bool {
	e, _ := utf8.DecodeRuneInString(emoji)
	if e == utf8.RuneError {
		return false
	}
	return (e >= 0x1F600 && e <= 0x1F64F) || // Emoticon
		(e >= 0x1F300 && e <= 0x1F5FF) || // Simboli e pittogrammi vari
		(e >= 0x1F680 && e <= 0x1F6FF) || // Trasporti e simboli mappa
		(e >= 0x1F700 && e <= 0x1F77F) || // Simboli alchemici
		(e >= 0x2600 && e <= 0x26FF) || // Simboli vari
		(e >= 0x2700 && e <= 0x27BF) || // Dingbats
		(e >= 0xFE00 && e <= 0xFE0F) || // Variazioni selettori
		(e >= 0x1F900 && e <= 0x1F9FF) || // Simboli supplementari
		(e >= 0x1FA70 && e <= 0x1FAFF) || // Emoji aggiuntivi
		(e >= 0x1F1E6 && e <= 0x1F1FF) // Bandiere (regioni)
}

type Message struct {
	ConvId     int        `json: "conversation"`
	SenderId   int        `json: "sender"`
	msgId      int        `json: "id"`
	ResponseTo int        `json: "responseTo"`
	Text       string     `json: "text"`
	Image      string     `json: "image"`
	Timestamp  time.Time  `json: "timestamp"`
	Checkmark  string     `json: "checkmark"`
	Reactions  []Reaction `json: "reactions"`
}

func (m *Message) ConvertMessage(msg database.Message) {
	m.ConvId = msg.ConvId
	m.SenderId = msg.SenderId
	m.msgId = msg.MsgId
	m.ResponseTo = m.ResponseTo
	m.Text = msg.Text
	m.Image = msg.Image
	m.Timestamp = msg.Timestamp
	m.Checkmark = msg.Checkmark

	for i := range msg.Reactions {
		var reac Reaction
		reac.ConvertReaction(msg.Reactions[i])
		m.Reactions = append(m.Reactions, reac)
	}

}

type Conversation struct {
	Id           int       `json: "id"`
	Name         string    `json: "name"`
	Image        string    `json: "image"`
	Group        bool      `json: "group"`
	Participants []User    `json: "participants"`
	Messages     []Message `json: "messages"`
}

func (c *Conversation) ConvertConversation(conv database.Conversation) error {
	c.Id = conv.Id
	c.Name = conv.Name
	c.Group = conv.Group

	for i := range conv.Participants {
		var u User
		err := u.ConvertUser(conv.Participants[i])
		if err != nil {
			return err
		}
		c.Participants = append(c.Participants, u)
	}

	for i := range conv.Messages {
		var m Message
		m.ConvertMessage(conv.Messages[i])
		c.Messages = append(c.Messages, m)
	}

	if c.Group {
		image, err := images.ConvertToBase64(images.SetDefaultGroupImage(conv.Id))
		if err != nil {
			return err
		}
		c.Image = image
	} else {
		image, err := images.ConvertToBase64(images.SetDefaultUserImage(conv.UserId))
		if err != nil {
			return err
		}
		c.Image = image
	}

	return nil
}

type Preview struct {
	Id          int     `json: "id"`
	Name        string  `json: "name"`
	Image       string  `json: "image"`
	Group       bool    `json: "group"`
	LastMessage Message `json: "lastmessage"`
}

func (p *Preview) ConvertPreview(prev database.Preview) error {
	p.Id = prev.Id
	p.Name = prev.Name
	p.Group = prev.Group
	var msg Message
	msg.ConvertMessage(prev.LastMessage)
	p.LastMessage = msg

	if p.Group {
		image, err := images.ConvertToBase64(images.SetDefaultGroupImage(prev.Id))
		if err != nil {
			return err
		}
		p.Image = image
	} else {
		image, err := images.ConvertToBase64(images.SetDefaultUserImage(prev.UserId))
		if err != nil {
			return err
		}
		p.Image = image
	}

	return nil
}
