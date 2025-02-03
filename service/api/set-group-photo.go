package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/images"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) SetGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userId, err := strconv.Atoi(ps.ByName("usrId"))
	if err != nil {
		BadRequest(w, err, "Invalid userId", ctx)
		return
	}

	if checkAuthorization(w, ctx, userId) != nil {
		return
	}

	exists, err := rt.db.CheckUserById(userId)
	if err != nil {
		InternalServerError(w, err, "Error while checking the user", ctx)
		return
	}
	if !exists {
		BadRequest(w, err, "User doesn't exists", ctx)
		return
	}

	convId, err := strconv.Atoi(ps.ByName("convId"))
	if err != nil {
		BadRequest(w, err, "Invalid convId", ctx)
		return
	}

	exists, is_group, err := rt.db.CheckConversationById(convId)
	if err != nil {
		InternalServerError(w, err, "Error checking the conversation", ctx)
		return
	}
	if !exists {
		BadRequest(w, err, "The conversation doesn't exists", ctx)
		return
	}
	ok, err := rt.db.IsParticipant(convId, userId)
	if err != nil {
		InternalServerError(w, err, "Couldn't check Group existance", ctx)
		return
	}
	if !ok || !is_group {
		Forbidden(w, nil, "The user isn't a member of this group or the conversation isn't a group", ctx)
		return
	}

	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		BadRequest(w, err, "The file is too big", ctx)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		BadRequest(w, err, "Couldn't access the image file from the request", ctx)
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		InternalServerError(w, err, "Couldn't read the image file from the request", ctx)
		return
	}
	defer func() { err = file.Close() }()

	filetype := http.DetectContentType(data)
	if filetype != "image/jpeg" {
		http.Error(w, "Can only use .jpeg images", http.StatusBadRequest)
		return
	}
	defer func() { err = file.Close() }()

	path := images.SetDefaultGroupImage(userId)
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		InternalServerError(w, err, "Couldn't copy the new image", ctx)
		return
	}

	err = images.SaveImage(path, 250, 250)
	if err != nil {
		InternalServerError(w, err, "Couldn't save the image", ctx)
	}

	type Response struct {
		image string `json: "image"`
	}
	var response Response
	img, err := images.ConvertToBase64(path)
	if err != nil {
		BadRequest(w, err, "Error taking image from storage", ctx)
		return
	}
	response.image = img

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		InternalServerError(w, err, "Couldn't encode the response", ctx)
		return
	}
}
