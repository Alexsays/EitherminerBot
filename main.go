package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Alexsays/EithermineBot/coinbase"
	"github.com/Alexsays/EithermineBot/database"
	"github.com/Alexsays/EithermineBot/ethermine"
	"github.com/Alexsays/EithermineBot/models"
	"github.com/Alexsays/EithermineBot/repositories"
	"github.com/Alexsays/EithermineBot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var minerKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📊 Dashboard", "/dashboard"),
		tgbotapi.NewInlineKeyboardButtonData("💵 Payouts", "/payouts"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🛑 Change Token", "/changeToken"),
		tgbotapi.NewInlineKeyboardButtonData("📈 Current Stats", "/currentStats"),
		// tgbotapi.NewInlineKeyboardButtonData("🗒 History", "/history"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("💰 ETH price", "/price"),
	),
	// tgbotapi.NewInlineKeyboardRow(
	// 	tgbotapi.NewInlineKeyboardButtonData("🔄 Rounds", "/rounds"),
	// 	tgbotapi.NewInlineKeyboardButtonData("⚙️ Settings", "/settings"),
	// ),
)

var waitingForToken = false

func createRefreshKeyboard(action string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Refresh", action),
		),
	)
}

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

func createMessageForCommand(command string, db *gorm.DB, telegramUser *tgbotapi.User) string {
	msg := ""

	switch command {
	case "/dashboard":
		ethermine.Wallet = getUser(db, telegramUser).EtherminerToken
		dashboard := ethermine.GetDashboard()
		dashboardStr := "📊 Dashboard\n\n"
		dashboardStr += "💰 Credits\n"
		unpaid := dashboard.Data.CurrentStatistics.Unpaid / math.Pow(1000, 6)
		dashboardStr += "- Unpaid: " + fmt.Sprintf("%.5f", unpaid) + " ETH\n"
		validShares := dashboard.Data.CurrentStatistics.ValidShares
		invalidShares := dashboard.Data.CurrentStatistics.InvalidShares
		staleShares := dashboard.Data.CurrentStatistics.StaleShares
		dashboardStr += fmt.Sprint("- Valid, Stale, Invalid Shares: ", validShares, " / ", staleShares, " / ", invalidShares, "\n")
		currentHashrate := fmt.Sprintf("%.1f", dashboard.Data.CurrentStatistics.CurrentHashrate/math.Pow(1000, 2))
		averageHashrate := fmt.Sprintf("%.1f", dashboard.Data.CurrentStatistics.AverageHashrate/math.Pow(1000, 2))
		reportedHashrate := fmt.Sprintf("%.1f", dashboard.Data.CurrentStatistics.ReportedHashrate/math.Pow(1000, 2))
		dashboardStr += fmt.Sprint("- Current, Average, Reported Hashrate: \n", currentHashrate, "MH/s / ", averageHashrate, "MH/s / ", reportedHashrate, "MH/s\n")
		dashboardStr += "\n🛠 Workers\n"
		dashboardStr += fmt.Sprint("- Active workers: ", dashboard.Data.CurrentStatistics.ActiveWorkers, " / ", len(dashboard.Data.Workers), "\n")
		for i := 0; i < len(dashboard.Data.Workers); i++ {
			dashboardStr += "* " + dashboard.Data.Workers[i].Worker + "\n"
			lastSeen := dashboard.Data.Workers[i].LastSeen
			lastSeenTime := time.Unix(int64(lastSeen), 0).Format("02-01-2006 15:04")
			dashboardStr += "* * Last seen: " + lastSeenTime + "\n"
		}
		msg = dashboardStr
	case "/payouts":
		ethermine.Wallet = getUser(db, telegramUser).EtherminerToken
		payouts := ethermine.GetPayouts()
		payoutsStr := "💵 Payouts\n\n"
		totalPayouts := 0.0
		for i := 0; i < len(payouts.Data); i++ {
			totalPayouts += payouts.Data[0].Amount
		}
		payoutsStr += "Total payouts: " + fmt.Sprintf("%.5f", totalPayouts/math.Pow(1000, 6)) + " ETH \n\n"
		for i := 0; i < len(payouts.Data); i++ {
			paidOn := payouts.Data[i].PaidOn
			paidOnTime := time.Unix(int64(paidOn), 0).Format("02-01-2006 15:04")
			payoutsStr += "- Paid on: " + paidOnTime + "\n"
			amount := payouts.Data[i].Amount / math.Pow(1000, 6)
			payoutsStr += "- - Amount: " + fmt.Sprintf("%.5f", amount) + " ETH\n"
			payoutsStr += "- - From / To Block: " + fmt.Sprintf("%.0f", payouts.Data[i].Start) + " / " + fmt.Sprintf("%.0f", payouts.Data[i].End) + "\n"
			payoutsStr += "- - Tx Hash: " + payouts.Data[i].TxHash + "\n\n"
		}
		msg = payoutsStr
	case "/currentStats":
		ethermine.Wallet = getUser(db, telegramUser).EtherminerToken
		currentStats := ethermine.GetCurrentStats()
		currentStatsStr := "📈 Current Statistics\n\n"
		dataTime := currentStats.Data.Time
		currentTime := time.Unix(int64(dataTime), 0).Format("02-01-2006 15:04")
		currentHashrate := fmt.Sprintf("%.1f", currentStats.Data.CurrentHashrate/math.Pow(1000, 2))
		averageHashrate := fmt.Sprintf("%.1f", currentStats.Data.AverageHashrate/math.Pow(1000, 2))
		reportedHashrate := fmt.Sprintf("%.1f", currentStats.Data.ReportedHashrate/math.Pow(1000, 2))
		currentStatsStr += fmt.Sprint("- Current, Average, Reported Hashrate: \n", currentHashrate, "MH/s / ", averageHashrate, "MH/s / ", reportedHashrate, "MH/s\n")
		validShares := currentStats.Data.ValidShares
		invalidShares := currentStats.Data.InvalidShares
		staleShares := currentStats.Data.StaleShares
		currentStatsStr += fmt.Sprint("- Valid, Stale, Invalid Shares: ", validShares, " / ", staleShares, " / ", invalidShares, "\n")
		unpaid := currentStats.Data.Unpaid / math.Pow(1000, 6)
		currentStatsStr += "- Unpaid: " + fmt.Sprintf("%.5f", unpaid) + " ETH\n"
		unconfirmed := currentStats.Data.Unconfirmed / math.Pow(1000, 6)
		currentStatsStr += "- Unconfirmed: " + fmt.Sprintf("%.5f", unconfirmed) + " ETH\n"
		coinsPerMin := currentStats.Data.CoinsPerMin
		currentStatsStr += "- Coins per minute: " + fmt.Sprintf("%.8f", coinsPerMin) + " ETH\n"
		usdPerMin := currentStats.Data.UsdPerMin
		currentStatsStr += "- USD per minute: " + fmt.Sprintf("%.8f", usdPerMin) + " $\n"
		btcPerMin := currentStats.Data.BtcPerMin
		currentStatsStr += "- Bitcoin per minute: " + fmt.Sprintf("%.8f", btcPerMin) + " ₿\n"
		currentStatsStr += "- Current Time: " + currentTime + "\n"
		msg = currentStatsStr
	case "/changeToken":
		waitingForToken = true
		changeTokenStr := "🛑 Write now your Ethermine token or something else to skip this"
		msg = changeTokenStr
	case "/price":
		price := coinbase.GetPrice()
		priceStr := "💰 Price \n\n"
		priceStr += "- Amount: " + price.Data.Amount + " €"
		msg = priceStr
	default:
		msg = "🏗 Not available"
	}

	return msg
}

func databaseConnection() *gorm.DB {
	db, err := database.ConnectToDB("postgres", os.Getenv("POSTGRES_PASSWORD"), "eithermine_bot")
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
	err := godotenv.Load(".env")
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
							"provide me your Ethermine token (sample: 716B383fA19526Lh73sd44353B3655e0339b513d)"
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

			bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
				ChatID:    update.CallbackQuery.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.MessageID,
			})

			msg.Text = createMessageForCommand(callback.Data, db, update.CallbackQuery.From)
			msg.ReplyMarkup = createRefreshKeyboard(callback.Data)

			sendMessage(msg, bot)
		}
	}
}
