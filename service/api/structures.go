package api

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/images"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
)

type User struct {
	Id    int    `json: "id"`
	Name  string `json: "name"`
	Image string `json: "image"`
}

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
