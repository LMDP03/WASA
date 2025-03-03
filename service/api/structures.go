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
	if user.Id == 0 {
		u.Image = ""
	} else {
		image, err := images.ConvertToBase64(images.SetDefaultUserImage(user.Id))
		if err != nil {
			return err
		}
		u.Image = image
	}
	return nil
}

type Reaction struct {
	Sender User   `json: "sender"`
	Emoji  string `json: "emoji"`
}

func (r *Reaction) ConvertReaction(reac database.Reaction) error {
	err := r.Sender.ConvertUser(reac.Sender)
	r.Emoji = reac.Emoji
	return err
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

type Response struct {
	Sender User
	MsgId  int
	Text   string
	Image  string
}

func (r *Response) ConvertResponse(res database.Response) error {
	err := r.Sender.ConvertUser(res.Sender)
	if err != nil {
		return err
	}
	r.MsgId = res.MsgId
	r.Text = res.Text
	r.Image = res.Image
	return nil
}

type Message struct {
	ConvId     int        `json: "conversation"`
	Sender     User       `json: "sender"`
	MsgId      int        `json: "id"`
	ResponseTo Response   `json: "responseTo"`
	Text       string     `json: "text"`
	Image      string     `json: "image"`
	Timestamp  time.Time  `json: "timestamp"`
	Checkmark  string     `json: "checkmark"`
	Reactions  []Reaction `json: "reactions"`
}

func (m *Message) ConvertMessage(msg database.Message) error {
	m.ConvId = msg.ConvId
	err := m.Sender.ConvertUser(msg.Sender)
	if err != nil {
		return err
	}
	m.MsgId = msg.MsgId
	if m.ResponseTo.MsgId != 0 {
		err = m.ResponseTo.ConvertResponse(msg.ResponseTo)
		if err != nil {
			return err
		}
	}

	m.Text = msg.Text
	m.Image = msg.Image
	m.Timestamp = msg.Timestamp
	m.Checkmark = msg.Checkmark
	m.Reactions = make([]Reaction, len(msg.Reactions))

	for i := range msg.Reactions {
		var reac Reaction
		err = reac.ConvertReaction(msg.Reactions[i])
		if err != nil {
			return err
		}
		m.Reactions[i] = reac
	}
	return nil
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
	c.Participants = make([]User, len(conv.Participants))
	c.Messages = make([]Message, len(conv.Messages))

	for i := range conv.Participants {
		var u User
		err := u.ConvertUser(conv.Participants[i])
		if err != nil {
			return err
		}
		c.Participants[i] = u
	}

	for i := range conv.Messages {
		var m Message
		err := m.ConvertMessage(conv.Messages[i])
		if err != nil {
			return err
		}
		c.Messages[i] = m
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
	err := msg.ConvertMessage(prev.LastMessage)
	if err != nil {
		return err
	}
	p.LastMessage = msg

	if p.Group {
		image, err := images.ConvertToBase64(images.SetDefaultGroupImage(prev.Id))
		if err != nil {
			return err
		}
		p.Image = image
	} else {
		if prev.UserId == 0 {
			p.Image = "No messages sent yet."
		} else {
			p.Image, err = images.ConvertToBase64(images.SetDefaultUserImage(prev.UserId))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
