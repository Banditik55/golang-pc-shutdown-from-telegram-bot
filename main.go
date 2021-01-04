package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	env "github.com/joho/godotenv"
)

var (
	count    = 0
	started  = false
	access   = false
	adminID  int
	bot      *tgbotapi.BotAPI
	keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("+15min"),
			tgbotapi.NewKeyboardButton("+30min"),
			tgbotapi.NewKeyboardButton("+60min"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("check"),
			tgbotapi.NewKeyboardButton("destroy"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("shutdown now"),
		),
	)
)

func main() {
	go interval()
	go checkAccess()
	renderBot()
}

func renderBot() {
	err := env.Load(".env")
	if err != nil {
		panic("error loading .env file")
	}

	token := os.Getenv("TG_TOKEN")
	_bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	bot = _bot
	fmt.Println("Connected to telegram bot", bot.Self.UserName)
	// bot.Debug = true

	_adminID, err := strconv.Atoi(os.Getenv("ADMIN_ID"))
	if err != nil {
		panic(err)
	}
	adminID = _adminID

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := createMessage(int(update.Message.Chat.ID), "")
		if update.Message.From.ID != adminID {
			msg.Text = "sorry, this is a private bot ¯\\_(ツ)_/¯"
			bot.Send(msg)
			continue
		}

		cmd := getCommand(update)
		switch cmd {
		case "start":
			msg.Text = "welcome"
			msg.ReplyMarkup = keyboard
		default:
			text := update.Message.Text
			switch text {
			case "+15min":
				msg.Text = "+15min"
				started = true
				count = count + (60 * 15)
			case "+30min":
				msg.Text = "+30min"
				started = true
				count = count + (60 * 30)
			case "+60min":
				msg.Text = "+60min"
				started = true
				count = count + (60 * 60)
			case "check":
				newCount := int(math.Floor(float64(count) / 60))
				msg.Text = strconv.Itoa(newCount) + " min. (" + strconv.Itoa(count) + " sec.)"
			case "destroy":
				msg.Text = "destroyed"
				started = false
				count = 0
			case "shutdown now":
				started = false
				count = 0
				shutdown()
			default:
				msg.Text = "use /start"
			}
		}
		bot.Send(msg)
	}
}

func shutdown() {
	if access != true {
		bot.Send(createMessage(adminID, "access denied"))
		log.Println("access denied")
		return
	}
	bot.Send(createMessage(adminID, "shutdown"))
	if err := exec.Command("sudo", "shutdown", "-h", "+1").Run(); err != nil {
		panic(err)
	}
}

func interval() {
	for {
		if started && count > 0 {
			count--
			if count == 0 {
				started = false
				count = 0
				shutdown()
			} else if count == 60 {
				bot.Send(createMessage(adminID, "1 min left"))
			} else if count == 300 {
				bot.Send(createMessage(adminID, "5 min left"))
			} else if count == 600 {
				bot.Send(createMessage(adminID, "10 min left"))
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func checkAccess() {
	time.Sleep(time.Second * 3)
	access = true
}

func createMessage(id int, text string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(int64(id), text)
}

func getCommand(update tgbotapi.Update) string {
	return strings.ToLower(update.Message.Command())
}
