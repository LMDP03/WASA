package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	//Session routes
	rt.router.POST("/session", rt.wrap(rt.doLogin, false))

	//User routes
	rt.router.GET("/users/:usrId/others", rt.wrap(rt.GetUsers, true))
	rt.router.PUT("/users/:usrId/image", rt.wrap(rt.SetMyPhoto, true))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
