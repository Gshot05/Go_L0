Первым делом, нужно скачать nats-сервер по этой ссылке: https://github.com/nats-io/nats-server/releases/tag/v2.10.11

Чтобы запустить сервер, переходим в папку cmd, затем в server, затем 
```sh
go run main.go
```

Чтобы запустить nats, пепреходим в папку cmd, затем в nats, затем 
```sh
go run nats.go
```

Затем заходим в браузер, и переходим по http://localhost:8080/web, и мы попадём на главную страницу
