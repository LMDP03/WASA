package api

import (
	"io"
	"net/http"
	"strconv"

	"wasatext/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) SetMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	if r.Method != http.MethodPut {
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		BadRequest(w, err, "Couldn't read the body", ctx)
		return
	}

	userName := string(body)

	if len(userName) < 3 || len(userName) > 16 {
		BadRequest(w, err, "Invalid username", ctx)
		return
	}

	err = rt.db.SetUserNameById(userName, userId)
	if err != nil {
		BadRequest(w, err, "Username already exists", ctx)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "plain/text")
	_, err = w.Write([]byte(userName))
	if err != nil {
		InternalServerError(w, err, "Couldn't encode the response", ctx)
		return
	}

}
