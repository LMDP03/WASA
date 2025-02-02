package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) CreateGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	if r.Method != http.MethodGet {
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

	type RequestConv struct {
		Name         string `json: "name"`
		Group        bool   `json: "group"`
		participants []string `json: "participants"`
	}

	var conv RequestConv
	if err := json.NewDecoder(r.Body).Decode(&conv); err != nil {
		BadRequest(w, err, "Couldn't decode the request", ctx)
		return
	}

	if !conv.Group {
		if len(conv.Name) < 3 || len(conv.Name) > 16 || len(conv.participants) != 2 {
			BadRequest(w, nil, "Invalid username", ctx)
			return
		} 
	}

	if conv.Group{
		if len(conv.Name) < 1 || len(conv.Name) > 20 ||  len(conv.participants) < 3 || len(conv.participants) > 50  {
		BadRequest(w, nil, "Invalid username", ctx)
		return
	}

	for i := range conv.participants {
		if len(conv.participants[i]) < 3 || len(conv.participants[i]) > 16 {
			BadRequest(w, nil, "One or more members have invalid names", ctx)
			return
		}
	}

	dbConv, err := rt.db.CreateConversation(conv.Name, conv.Group)
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
	dbParticipants, err := rt.db.AddParticipants(dbConv.Id, conv.participants)
	if err != nil {
		InternalServerError(w, err, "Error while getting the participants", ctx)
		return
	}
	conversation.Participants = dbParticipants
	

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		ctx.Logger.Error("Couldn't encode the response", err)
		http.Error(w, "Couldn't encode the response", http.StatusInternalServerError)
		return
	}
	

}
