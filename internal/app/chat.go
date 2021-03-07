package app

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8877", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request, settings *settings) {

	log.Println(r.URL)

	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	for _, cookie := range r.Cookies() {
		fmt.Println("Found a cookie named: ", cookie.Name, cookie.Value)
	}

	generator := rand.New(rand.NewSource(time.Now().UnixNano()))

	settings.roomId = "room-1"
	settings.userId = generator.Int()

	http.ServeFile(w, r, "./web/index.html")
}

type settings struct {
	roomId string
	userId int
}

func Start() {

	settings := &settings{}

	go h.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHome(w, r, settings)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, settings)
	})

	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}