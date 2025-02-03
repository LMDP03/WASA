package api

import (
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) DeleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	if r.Method != http.MethodDelete {
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

	err = rt.db.DeleteMessage(convId, msgId)
	if err != nil {
		InternalServerError(w, err, "Error while deleting the message", ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
