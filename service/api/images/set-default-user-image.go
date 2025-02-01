package images

import (
	"fmt"
)

func SetDefaultUserImage(userId int) string {
	return fmt.Sprintf("./storage/profiles/%d/profile_image.jpg", userId)
}
