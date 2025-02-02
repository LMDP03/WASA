package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	var searchName string
	if !r.URL.Query().Has("srcName") {
		searchName = ""
	} else {
		searchName = r.URL.Query().Get("srcName")
		if len(searchName) > 20 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	dbConvs, err := rt.db.GetConversationsbyName(searchName, userId)
	if err != nil {
		ctx.Logger.Error("Couldn't find conversations for this user", err)
		http.Error(w, "Couldn't find conversations for this user", http.StatusInternalServerError)
		return
	}

	conversations := make([]Conversation, len(dbConvs))

	for i, dbConv := range dbConvs {
		var conv Conversation
		err := conv.ConvertConversation(dbConv)
		if err != nil {
			ctx.Logger.Error("Couldn't load conversations properly", err)
			http.Error(w, "Couldn't load conversations properly", http.StatusInternalServerError)
			return
		}
		conversations[i] = conv
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(conversations); err != nil {
		ctx.Logger.Error("Couldn't encode the response", err)
		http.Error(w, "Couldn't encode the response", http.StatusInternalServerError)
		return
	}
}
