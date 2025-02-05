package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasatext/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userid, err := strconv.Atoi(ps.ByName("usrId"))
	if err != nil {
		BadRequest(w, err, "Invalid userId", ctx)
		return
	}
	if checkAuthorization(w, ctx, userid) != nil {
		return
	}

	exists, err := rt.db.CheckUserById(userid)
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
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	dbUsers, err := rt.db.GetUsersByName(searchName, userid)
	if err != nil {
		ctx.Logger.Error("Couldn't find users", err)
		http.Error(w, "Couldn't find users", http.StatusInternalServerError)
		return
	}

	users := make([]User, len(dbUsers))

	for i, dbUser := range dbUsers {
		var user User
		err := user.ConvertUser(dbUser)
		if err != nil {
			ctx.Logger.Error("Couldn't load users properly", err)
			http.Error(w, "Couldn't load users properly", http.StatusInternalServerError)
			return
		}
		users[i] = user
	}

	type Response struct {
		users []User `json: "users"`
	}

	var res = Response{
		users: users,
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		ctx.Logger.Error("Couldn't encode the response", err)
		http.Error(w, "Couldn't encode the response", http.StatusInternalServerError)
		return
	}
}
