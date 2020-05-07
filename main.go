package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"./database"
	"./ethermine"
	"./models"
	"./repositories"
	"./services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
)

var minerKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📊 Dashboard", "/dashboard"),
		tgbotapi.NewInlineKeyboardButtonData("📈 Current Stats", "/currentStats"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🗒 History", "/history"),
		tgbotapi.NewInlineKeyboardButtonData("💵 Payouts", "/payouts"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔄 Rounds", "/rounds"),
		tgbotapi.NewInlineKeyboardButtonData("⚙️ Settings", "/settings"),
	),
)

func createUser(db *gorm.DB, telegramUser *tgbotapi.User, token string) {
	userRepository := repositories.NewUserRepository(db)

	if !userExist(db, telegramUser) {
		var user models.User
		user.FirstName = telegramUser.FirstName
		user.LastName = telegramUser.LastName
		user.Username = telegramUser.UserName
		user.TelegramID = strconv.Itoa(telegramUser.ID)
		user.EtherminerToken = token

		services.CreateUser(&user, *userRepository)
	}
}

func userExist(db *gorm.DB, telegramUser *tgbotapi.User) bool {
	userRepository := repositories.NewUserRepository(db)

	resp := services.FindOneUserByUsername(telegramUser.UserName, *userRepository)

	return resp.Success
}

func getUser(db *gorm.DB, telegramUser *tgbotapi.User) *models.User {
	userRepository := repositories.NewUserRepository(db)

	resp := services.FindOneUserByUsername(telegramUser.UserName, *userRepository)

	return resp.Data.(*models.User)
}

func minerCallbacks(db *gorm.DB, data string, telegramUser *tgbotapi.User) string {
	msg := ""

	switch data {
	case "/dashboard":
		ethermine.Wallet = getUser(db, telegramUser).EtherminerToken
		dashboard := ethermine.GetDashboard()
		dashboardStr := "🛠 Workers\n"
		dashboardStr += "- " + dashboard.Data.Workers[0].Worker
		msg = dashboardStr
	}

	return msg
}

func databaseConnection() *gorm.DB {
	db, err := database.ConnectToDB("postgres", "pass", "eithermine_bot")
	if err != nil {
		panic(err)
	}

	err = db.DB().Ping()
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.User{})

	return db
}

func sendMessage(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err) // You should add better error handling than this!
	}

	db := databaseConnection()

	bot.Debug = true // Has the library display every request and response.

	log.Printf("Started '%s'", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	waitingForToken := false

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		if update.Message != nil {
			if waitingForToken {
				ethermine.Wallet = update.Message.Text
				dashboard := ethermine.GetDashboard()

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
					ChatID:    update.Message.Chat.ID,
					MessageID: update.Message.MessageID,
				})

				if strings.ToLower(dashboard.Status) != "error" {
					waitingForToken = false
					createUser(db, update.Message.From, update.Message.Text)
					msg.Text = "⛏ Miner registered, use /miner to start"
				} else {
					msg.Text = "😔 This token is not available. Write your Ethermine token again"
				}

				sendMessage(msg, bot)
			} else {
				switch update.Message.Text {
				case "/start":
					bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
						ChatID:    update.Message.Chat.ID,
						MessageID: update.Message.MessageID,
					})

					if !userExist(db, update.Message.From) {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
						msg.Text = "💻 Welcome to Eitherminer. To start getting info about your Ethermine account, " +
							"write provide me your Ethermine token (sample: 716B383fA19526Lh73sd44353B3655e0339b513d)"
						waitingForToken = true

						sendMessage(msg, bot)
					}
				case "/miner":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
						ChatID:    update.Message.Chat.ID,
						MessageID: update.Message.MessageID,
					})
					msg.Text = "⛏ Choose one Miner option:"
					msg.ReplyMarkup = minerKeyboard

					sendMessage(msg, bot)
				}
			}
		} else if callback := update.CallbackQuery; callback != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

			msg.Text = minerCallbacks(db, callback.Data, update.CallbackQuery.From)

			sendMessage(msg, bot)
		}
	}
}
