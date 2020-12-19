![size](https://img.shields.io/github/size/banditik55/golang-pc-shutdown-from-telegram-bot/main.go?style=plastic)

<p align="center">
    <a href="https://github.com/Banditik55/golang-pc-shutdown-from-telegram-bot">
        <img width="600px" src="./img/gopher.png"/>
    </a>
    <h1 align="center">Golang PC shutdown from telegram bot</h2>
</p>

## Installation

```bash
go get github.com/go-telegram-bot-api/telegram-bot-api github.com/joho/godotenv
```

Create .env and write "TG_TOKEN" & "ADMIN_ID"

"TG_TOKEN" is your telegram bot token

"ADMIN_ID" is your telegram user id

## Build and run

```bash
go build main.go
sudo ./main
```

You can uncomment "bot.Debug = true" (line 57) if you want see debug

Tested on GNU/Linux

This is my first public repository and I just started learning golang. Do not judge strictly :)