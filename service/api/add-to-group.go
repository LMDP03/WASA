package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasatext/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) AddToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	type Request struct {
		participants []string `json: "participants"`
	}

	var req Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, err, "Couldn't decode the request", ctx)
		return
	}

	for i := range req.participants {
		if len(req.participants[i]) < 3 || len(req.participants[i]) > 16 {
			BadRequest(w, nil, "One or more members have invalid names", ctx)
			return
		}
	}

	err = rt.db.AddParticipants(convId, req.participants)
	if err != nil {
		InternalServerError(w, err, "Error while adding participants to group", ctx)
		return
	}

	dbMembers, err := rt.db.GetParticipants(convId)
	var members = make([]User, len(dbMembers))
	for i := range dbMembers {
		var u User
		err := u.ConvertUser(dbMembers[i])
		if err != nil {
			InternalServerError(w, err, "Error converting users", ctx)
			return
		}
		members[i] = u
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(members); err != nil {
		ctx.Logger.Error("Couldn't encode the response", err)
		http.Error(w, "Couldn't encode the response", http.StatusInternalServerError)
		return
	}

}
