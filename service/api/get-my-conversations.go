package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasatext/service/api/reqcontext"

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

	exists, err := rt.db.CheckUserById(userId)
	if err != nil {
		InternalServerError(w, err, "Error while checking the user", ctx)
		return
	}
	if !exists {
		BadRequest(w, err, "User doesn't exists", ctx)
		return
	}

	var searchName string
	if !r.URL.Query().Has("srcName") {
		searchName = ""
	} else {
		searchName = r.URL.Query().Get("srcName")
		if len(searchName) > 20 {
			BadRequest(w, nil, "the name can be at most 20 characters", ctx)
			return
		}
	}

	dbPrevs, err := rt.db.GetConversations(userId, searchName)
	if err != nil {
		ctx.Logger.Error("Couldn't find conversations for this user", err)
		http.Error(w, "Couldn't find conversations for this user", http.StatusInternalServerError)
		return
	}

	previews := make([]Preview, len(dbPrevs))

	for i, dbPrev := range dbPrevs {
		var prev Preview
		err := prev.ConvertPreview(dbPrev)
		if err != nil {
			InternalServerError(w, err, "Couldn't transform previews", ctx)
			return
		}
		previews[i] = prev
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(previews); err != nil {
		ctx.Logger.Error("Couldn't encode the response", err)
		http.Error(w, "Couldn't encode the response", http.StatusInternalServerError)
		return
	}
}
