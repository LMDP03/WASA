package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	convId, err := strconv.Atoi(ps.ByName("convId"))
	if err != nil {
		BadRequest(w, err, "Invalid convId", ctx)
		return
	}

	dbConv, err := rt.db.GetConversationById(convId)
	if err != nil {
		InternalServerError(w, err, "Couldn't get the conversation from the database", ctx)
		return 
	}

	var conv Conversation
	err = conv.ConvertConversation(dbConv)
	if err != nil {
		InternalServerError(w, err, "Error while converting the conversation", ctx)
		return
	}

	dbMessages, err :=
}