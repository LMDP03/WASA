package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Session routes
	rt.router.POST("/session", rt.wrap(rt.doLogin, false))

	// Users routes
	rt.router.GET("/users/:usrId/others", rt.wrap(rt.GetUsers, true))
	rt.router.PUT("/users/:usrId/image", rt.wrap(rt.SetMyPhoto, true))
	rt.router.PUT("/users/:usrId/name", rt.wrap(rt.SetMyUserName, true))

	// General conversations routes
	rt.router.GET("/users/:usrId/conversations", rt.wrap(rt.GetMyConversations, true))
	rt.router.POST("/users/:usrId/conversations/private", rt.wrap(rt.StartConversation, true))
	rt.router.POST("/users/:usrId/conversations/group", rt.wrap(rt.CreateGroup, true))

	// Specific conversation routes
	rt.router.GET("/users/:usrId/conversation/:convId", rt.wrap(rt.GetConversation, true))
	rt.router.POST("/users/:usrId/conversation/:convId", rt.wrap(rt.AddToGroup, true))
	rt.router.DELETE("/users/:usrId/conversation/:convId", rt.wrap(rt.LeaveGroup, true))
	rt.router.PUT("/users/:usrId/conversation/:convId/name", rt.wrap(rt.SetGroupName, true))
	rt.router.PUT("/users/:usrId/conversation/:convId/image", rt.wrap(rt.SetGroupPhoto, true))

	// Messages routes
	rt.router.POST("/users/:usrId/conversation/:convId/messages", rt.wrap(rt.SendMessage, true))
	rt.router.POST("/users/:usrId/conversation/:convId/messages/:msgId", rt.wrap(rt.ForwardMessage, true))
	rt.router.DELETE("/users/:usrId/conversation/:convId/messages/:msgId", rt.wrap(rt.DeleteMessage, true))

	// Reactions routes
	rt.router.POST("/users/:usrId/conversation/:convId/messages/:msgId/reactions", rt.wrap(rt.CommentMessage, true))
	rt.router.DELETE("/users/:usrId/conversation/:convId/messages/:msgId/reactions", rt.wrap(rt.UncommentMessage, true))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
