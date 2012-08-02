package handlers

import (
	. "github.com/paulbellamy/mango"
	"github.com/sunfmin/mangotemplate"
	"net/http"
)

type RenderData struct {
	Email         string
	WebSocketHost string
}

func Home(env Env) (status Status, headers Headers, body Body) {
	mangotemplate.ForRender(env, "chats/home", nil)
	headers = Headers{}
	return
}

func Join(env Env) (status Status, headers Headers, body Body) {
	email := env.Request().FormValue("email")
	if email == "" {
		return Redirect(http.StatusFound, "/")
	}

	r := env.Request()
	mangotemplate.ForRender(env, "chats/room", &RenderData{Email: email, WebSocketHost: r.Host})
	headers = Headers{}
	return
}
