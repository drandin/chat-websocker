package app

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8877", "http service address")

type settings struct {
	roomId string
	userId int
}

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

	setSettings(settings)
	render(w, settings)
}

// Устанавливаем название компаны и идентификатор пользлвателя
// В реальном приложении эти параметры могут быть заданы на основании
// данных, которые храняться в Cookies или как-либо иным способом
func setSettings(s *settings)  {
	s.roomId = "room-1"
	s.userId = rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}

func render(w http.ResponseWriter, s *settings)  {

	view := struct {
		RoomId string
		UserId int
	}{
		RoomId: s.roomId,
		UserId: s.userId,
	}

	tmpl, _ := template.ParseFiles("./web/index.html")
	_ = tmpl.Execute(w, view)

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