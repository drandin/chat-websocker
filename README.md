# Пример - шаблон чата на websocket

### Использование

http://localhost:8877

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