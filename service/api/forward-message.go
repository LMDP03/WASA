package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasatext/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) ForwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	dbMsgOrig, err := rt.db.GetMessageById(convId, msgId)
	if err != nil {
		InternalServerError(w, err, "Error while getting the original message", ctx)
		return
	}

	type Destination struct {
		Id     int  `json:"id"`
		IsConv bool `json:"isConv"`
	}

	var destinations []Destination
	if err := json.NewDecoder(r.Body).Decode(&destinations); err != nil {
		BadRequest(w, err, "Couldn't decode the request", ctx)
		return
	}

	for i := range destinations {

		if destinations[i].IsConv {
			exists, _, err := rt.db.CheckConversationById(destinations[i].Id)
			if err != nil {
				InternalServerError(w, err, "Error checking the destination", ctx)
				return
			}
			if !exists {
				BadRequest(w, err, "The destination doesn't exists", ctx)
				return
			}
			ok, err := rt.db.IsParticipant(destinations[i].Id, userId)
			if err != nil {
				InternalServerError(w, err, "Couldn't check destination existance", ctx)
				return
			}
			if !ok {
				Forbidden(w, nil, "The user isn't a member of at least one conversation", ctx)
				return
			}

			_, err = rt.db.CreateMessage(destinations[i].Id, userId, 0, dbMsgOrig.Text, dbMsgOrig.Image, true)
			if err != nil {
				InternalServerError(w, err, "Error while sending the message", ctx)
				return
			}

		} else {
			exists, err := rt.db.CheckUserById(userId)
			if err != nil {
				InternalServerError(w, err, "Error while checking the receiver", ctx)
				return
			}
			if !exists {
				BadRequest(w, err, "Receiver doesn't exists", ctx)
				return
			}
			var participants []string
			you, err := rt.db.GetUserById(userId)
			if err != nil {
				InternalServerError(w, err, "Couldn't get user", ctx)
				return
			}
			participants = append(participants, you.Name)

			other, err := rt.db.GetUserById(destinations[i].Id)
			if err != nil {
				InternalServerError(w, err, "Couldn't get other user", ctx)
				return
			}
			participants = append(participants, other.Name)

			dbConv, err := rt.db.CreateConversation("", false, userId, participants)
			if err != nil {
				InternalServerError(w, err, "Error while creating the conversation", ctx)
				return
			}

			_, err = rt.db.CreateMessage(dbConv.Id, userId, 0, dbMsgOrig.Text, dbMsgOrig.Image, true)
			if err != nil {
				InternalServerError(w, err, "Couldn't send the message", ctx)
				return
			}

			destinations[i].Id = dbConv.Id
			destinations[i].IsConv = true
		}

	}

	k := len(destinations) - 1

	dbConv, err := rt.db.GetConversationById(destinations[k].Id, userId)
	if err != nil {
		InternalServerError(w, err, "Error while getting the destination conv", ctx)
		return
	}

	var conv Conversation
	err = conv.ConvertConversation(dbConv)
	if err != nil {
		InternalServerError(w, err, "Error while converting the destination conv", ctx)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(conv); err != nil {
		InternalServerError(w, err, "Error encoding response", ctx)
		return
	}

}
