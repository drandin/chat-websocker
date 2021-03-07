# Пример - шаблон чата на websocket

### Использование

http://localhost:8877

Можно задать **roomId** и **userId**. 

Смотрите файл _internal/app/chat.go_. 

```go
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
```

### Настройка Nginx (если нужна)

Проксирование запрсов, если HTTPS:

wss://your-domain.ru/websocket ---> http://localhost:8877/ws

В секци **Server** добавить:

```
    location /websocket {
      proxy_pass http://localhost:8877/ws;
      proxy_http_version 1.1;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection "Upgrade";
   }
```