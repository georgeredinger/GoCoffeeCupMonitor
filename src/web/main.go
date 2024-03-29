package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"log"
	"net/http"
	"text/template"
	"path"
)

const (
	staticDir = "./"
)

var addr = flag.String("addr", ":1234", "http service address")
var homeTempl = template.Must(template.ParseFiles("home.html"))

func homeHandler(c http.ResponseWriter, req *http.Request) {
	p := req.URL.Path[1:]
	if p == "" {
		homeTempl.Execute(c, req.Host)
		return
	}
	p = path.Join(staticDir, p)
	http.ServeFile(c, req, p)
}



func main() {
	flag.Parse()
	go h.run()
	go getdata()
	http.HandleFunc("/", homeHandler)
	http.Handle("/ws", websocket.Handler(wsHandler))
  if err := http.ListenAndServe(*addr,nil ); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
