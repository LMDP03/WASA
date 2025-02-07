package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"wasatext/service/api/reqcontext"
	"wasatext/service/images"

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

	exists, err := rt.db.CheckUserById(userId)
	if err != nil {
		InternalServerError(w, err, "Error while checking the user", ctx)
		return
	}
	if !exists {
		BadRequest(w, err, "User doesn't exists", ctx)
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
		return
	}

	dbUser, err := rt.db.GetUserById(userId)
	if err != nil {
		InternalServerError(w, err, "Couldn't get the user", ctx)
		return
	}
	var user User
	err = user.ConvertUser(dbUser)
	if err != nil {
		InternalServerError(w, err, "Couldn't convert the user", ctx)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		InternalServerError(w, err, "Couldn't encode the response", ctx)
		return
	}

}
