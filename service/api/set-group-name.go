package api

import (
	"io"
	"net/http"
	"strconv"

	"wasatext/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) SetGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		BadRequest(w, err, "Couldn't read the body", ctx)
		return
	}

	groupName := string(body)

	if len(groupName) < 1 || len(groupName) > 20 {
		BadRequest(w, err, "Invalid name", ctx)
		return
	}

	err = rt.db.SetGroupNameById(convId, groupName)
	if err != nil {
		InternalServerError(w, err, "Error changing group name", ctx)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "plain/text")
	_, err = w.Write([]byte(groupName))
	if err != nil {
		InternalServerError(w, err, "Couldn't encode the response", ctx)
		return
	}

}
