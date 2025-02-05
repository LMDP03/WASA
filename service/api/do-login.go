package api

import (
	"encoding/json"
	"net/http"

	"wasatext/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		BadRequest(w, err, "Couldn't decode the request", ctx)
		return
	}

	if len(user.Name) < 3 || len(user.Name) > 16 {
		BadRequest(w, nil, "Invalid username", ctx)
		return
	}

	exists, err := rt.db.CheckUserByName(user.Name)
	if err != nil {
		InternalServerError(w, err, "Couldn't check user's existance", ctx)
		return
	}

	if !exists {
		dbUser, err := rt.db.CreateUser(user.Name)
		if err != nil {
			InternalServerError(w, err, "Couldn't create user", ctx)
			return
		}
		err = user.ConvertUser(dbUser)
		if err != nil {
			InternalServerError(w, err, "Couldn't create user", ctx)
			return
		}
		w.WriteHeader(http.StatusCreated)

	} else {
		dbUser, err := rt.db.GetUserByName(user.Name)
		if err != nil {
			InternalServerError(w, err, "Error getting user", ctx)
			return
		}
		err = user.ConvertUser(dbUser)
		if err != nil {
			InternalServerError(w, err, "Error converting user", ctx)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		InternalServerError(w, err, "Couldn't encode the response", ctx)
		return
	}
}
