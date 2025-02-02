package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/images"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) SetMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	path := images.SetDefaultUserImage(userId)
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
		Image string `json: "image"`
	}
	var response = Response{
		Image: base64.StdEncoding.EncodeToString(data),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		InternalServerError(w, err, "Couldn't encode the response", ctx)
		return
	}

}
