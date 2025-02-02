package images

import "fmt"

func SetDefaultConversationImage(convId int) string {
	return fmt.Sprintf("./storage/conversations/%d/chat_image.jpg", convId)
}
