package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"wasatext/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) SendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	if r.Method != http.MethodPost {
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

	exists, _, err = rt.db.CheckConversationById(convId)
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
	if !ok {
		Forbidden(w, nil, "The user isn't a member of this conversation", ctx)
		return
	}

	type Request struct {
		text  string `json: "text"`
		image string `json: "image"`
	}
	var req Request

	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		BadRequest(w, err, "The file is too big", ctx)
		return
	}

	req.text = r.FormValue("text")

	file, _, err := r.FormFile("image")

	if err == nil {
		data, err := io.ReadAll(file)
		if err != nil {
			InternalServerError(w, err, "Couldn't read image file", ctx)
			return
		}

		filetype := http.DetectContentType(data)
		if filetype != "image/jpeg" {
			http.Error(w, "Can only use .jpeg images", http.StatusBadRequest)
			return
		}
		defer func() { err = file.Close() }()
		req.image = base64.StdEncoding.EncodeToString(data)
	}

	var responseid int
	if !r.URL.Query().Has("responseTo") {
		responseid, err = strconv.Atoi(r.URL.Query().Get("responseTo"))
		if err != nil {
			BadRequest(w, err, "Couldn't read the responseId", ctx)
			return
		}
		ok, err = rt.db.CheckMessageById(convId, responseid)
		if err != nil {
			InternalServerError(w, err, "Error while checking the message to respond to", ctx)
			return
		}
		if !ok {
			BadRequest(w, err, "The message doesn't exist or isn't from this conversation", ctx)
			return
		}
	} else {
		responseid = 0
	}

	dbMsg, err := rt.db.CreateMessage(convId, userId, responseid, req.text, req.image)
	if err != nil {
		InternalServerError(w, err, "Error while sending the message", ctx)
		return
	}

	var msg Message
	msg.ConvertMessage(dbMsg)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		InternalServerError(w, err, "Error encoding response", ctx)
		return
	}

}
