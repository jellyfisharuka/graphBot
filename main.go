package main

import (
	"lilumaBot/internal/db"
	"lilumaBot/internal/telegram"

	"log"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func main() {
	db.ConnectDB()
	bot, err := tgbotapi.NewBotAPI("7468786356:AAEPEUg14ZO46FKEDLzdlIMCRFGx95m3OP0")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				// Отправка клавиатуры с выбором месяца
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose a month")
				msg.ReplyMarkup = telegram.CreateMonthKeyboard()
				bot.Send(msg)
			}
		} else if update.CallbackQuery != nil {
			data := update.CallbackQuery.Data
			var filePath string

			// Обработка выбора месяца
			if telegram.IsMonth(data) {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "You chose: "+data)
				msg.ReplyMarkup = telegram.CreateInfoKeyboard()
				bot.Send(msg)
			} else {
				// Получение финансовых данных
				data, err := telegram.FetchFinancialData(data)
				if err != nil {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Error fetching data: "+err.Error())
					bot.Send(msg)
					continue
				}

				// Создание графика и отправка
				filePath, err = telegram.CreateChart(data)
				if err != nil {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Error generating chart: "+err.Error())
					bot.Send(msg)
					continue
				}

				// Отправка графика
				photo := tgbotapi.NewPhotoUpload(update.CallbackQuery.Message.Chat.ID, filePath)
				_, err = bot.Send(photo)
				if err != nil {
					log.Println("Error sending photo:", err)
			}
		}
	}
}
}
