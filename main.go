package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/kobeld/gochatting/handlers"
	"github.com/kobeld/gochatting/wscon"
	. "github.com/paulbellamy/mango"
	"net/http"
)

func main() {
	l, r := handlers.LayoutAndRenderer()
	s := new(Stack)
	s.Middleware(l, r)

	http.Handle("/chat", websocket.Handler(wscon.BuildConnection))
	http.HandleFunc("/join", s.HandlerFunc(handlers.Join))
	http.HandleFunc("/", s.HandlerFunc(handlers.Home))
	http.HandleFunc("/public/", assetsHandler)

	go wscon.InitChatRoom()

	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[len("/"):])
}
