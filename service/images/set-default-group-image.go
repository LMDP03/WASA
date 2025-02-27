package images

import "fmt"

func SetDefaultGroupImage(convId int) string {
	return fmt.Sprintf("./storage/groups/%d/chat_image.jpg", convId)
}
