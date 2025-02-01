package api

import (
	"encoding/json"
	"net/http"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request) {

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(user.Name) < 3 || len(user.Name) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, err := rt.db.CheckUserByName(user.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		dbUser, err := rt.db.CreateUser(user.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = user.ConvertUser(dbUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	} else {
		dbUser, err := rt.db.GetUserByName(user.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = user.ConvertUser(dbUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(user)

}
