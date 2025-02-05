package api

import (
	"net/http"

	"wasatext/service/api/reqcontext"
)

func InternalServerError(w http.ResponseWriter, err error, msg string, ctx reqcontext.RequestContext) {
	ctx.Logger.WithError(err).Error(msg)
	w.WriteHeader(http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter, err error, msg string, ctx reqcontext.RequestContext) {
	if err != nil {
		http.Error(w, msg+": "+err.Error(), http.StatusBadRequest)
	} else {
		http.Error(w, msg, http.StatusBadRequest)
	}
}

func Forbidden(w http.ResponseWriter, err error, msg string, ctx reqcontext.RequestContext) {
	if err != nil {
		http.Error(w, msg+": "+err.Error(), http.StatusForbidden)
	} else {
		http.Error(w, msg, http.StatusForbidden)
	}
}
