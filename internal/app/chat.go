package app

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8877", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {

	log.Println(r.URL)

	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "./web/index.html")
}

func Start() {

	go h.run()

	roomId := "guest"
	userId := 123

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, roomId, userId)
	})

	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}