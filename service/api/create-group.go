package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasatext/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) CreateGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	type RequestConv struct {
		Name         string   `json: "name"`
		participants []string `json: "participants"`
	}

	var conv RequestConv
	if err := json.NewDecoder(r.Body).Decode(&conv); err != nil {
		BadRequest(w, err, "Couldn't decode the request", ctx)
		return
	}

	if len(conv.Name) < 1 || len(conv.Name) > 20 {
		BadRequest(w, nil, "Invalid username", ctx)
		return
	}

	for i := range conv.participants {
		if len(conv.participants[i]) < 3 || len(conv.participants[i]) > 16 {
			BadRequest(w, nil, "One or more members have invalid names", ctx)
			return
		}
	}

	dbConv, err := rt.db.CreateConversation(conv.Name, true, 0, conv.participants)
	if err != nil {
		InternalServerError(w, err, "Error while creating the conversation", ctx)
		return
	}

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
