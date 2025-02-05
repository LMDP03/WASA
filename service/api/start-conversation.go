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

func (rt *_router) StartConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	receiverid, err := strconv.Atoi(r.URL.Query().Get("rcvId"))
	if err != nil {
		BadRequest(w, err, "Couldn't reid the receiverId", ctx)
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

	exists, err = rt.db.CheckUserById(receiverid)
	if err != nil {
		InternalServerError(w, err, "Couldn't check receiver's existance", ctx)
		return
	}
	if !exists {
		BadRequest(w, nil, "Receiver doesn't exists", ctx)
		return
	}

	var participants []string
	you, err := rt.db.GetUserById(userId)
	if err != nil {
		InternalServerError(w, err, "Couldn't get user", ctx)
		return
	}
	participants = append(participants, you.Name)

	other, err := rt.db.GetUserById(receiverid)
	if err != nil {
		InternalServerError(w, err, "Couldn't get other user", ctx)
		return
	}
	participants = append(participants, other.Name)

	dbConv, err := rt.db.CreateConversation("", false, receiverid, participants)
	if err != nil {
		InternalServerError(w, err, "Error while creating the conversation", ctx)
		return
	}

	dbMsg, err := rt.db.CreateMessage(dbConv.Id, userId, 0, req.text, req.image)
	if err != nil {
		InternalServerError(w, err, "Couldn't send the message", ctx)
		return
	}
	dbConv.Messages = append(dbConv.Messages, dbMsg)

	var conversation Conversation
	err = conversation.ConvertConversation(dbConv)
	if err != nil {
		InternalServerError(w, err, "Error while converting the conversation", ctx)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		ctx.Logger.Error("Couldn't encode the response", err)
		http.Error(w, "Couldn't encode the response", http.StatusInternalServerError)
		return
	}

}
