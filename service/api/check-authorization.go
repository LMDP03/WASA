package api

import (
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
)

func checkAuthorization(w http.ResponseWriter, ctx reqcontext.RequestContext, userId int) error {
	if userId != ctx.UserId {
		Forbidden(w, nil, "Unauthorized user", ctx)
		return errors.New("Unauthorized user")
	}
	return nil
}
