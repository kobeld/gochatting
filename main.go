package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/kobeld/gochatting/accounts"
	"github.com/kobeld/gochatting/chat"
	. "github.com/paulbellamy/mango"
	"net/http"
)

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[len("/"):])
}

func main() {
	l, r := accounts.LayoutAndRenderer()
	s := new(Stack)
	s.Middleware(l, r)

	http.HandleFunc("/", s.HandlerFunc(chat.Index))
	http.Handle("/chat", websocket.Handler(chat.ReceiveChatMessage))
	http.HandleFunc("/public/", assetsHandler)

	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
