package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) CommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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
		Forbidden(w, nil, "The user isn't a member of this group", ctx)
		return
	}

	msgId, err := strconv.Atoi(ps.ByName("msgId"))
	if err != nil {
		BadRequest(w, err, "Invalid msgId", ctx)
		return
	}
	ok, err = rt.db.CheckMessageById(convId, msgId)
	if err != nil {
		InternalServerError(w, err, "Error while checking the original message", ctx)
		return
	}
	if !ok {
		BadRequest(w, err, "The message doesn't exist or isn't from this conversation", ctx)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		BadRequest(w, err, "Couldn't read the body", ctx)
		return
	}

	emoji := string(body)

	if len(emoji) != 1 || !CheckEmoji(emoji) {
		BadRequest(w, err, "Invalid emoji", ctx)
		return
	}

	exists, err = rt.db.CheckReactionById(convId, msgId, userId)
	if err != nil {
		InternalServerError(w, err, "Error While checking your comment", ctx)
		return
	}

	var reac Reaction
	if !exists {
		dbReac, err := rt.db.CreateReaction(convId, userId, msgId, emoji)
		if err != nil {
			InternalServerError(w, err, "Couldn't comment the message", ctx)
			return
		}
		reac.ConvertReaction(dbReac)
		w.WriteHeader(http.StatusCreated)
	} else {
		dbReac, err := rt.db.UpdateReaction(convId, userId, msgId, emoji)
		if err != nil {
			InternalServerError(w, err, "Couldn't update the comment", ctx)
			return
		}
		reac.ConvertReaction(dbReac)
		w.WriteHeader(http.StatusOK)
	}
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(reac); err != nil {
		InternalServerError(w, err, "Error encoding response", ctx)
		return
	}

}
