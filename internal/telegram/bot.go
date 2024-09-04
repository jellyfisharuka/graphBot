package telegram

import (
	"errors"
	"lilumaBot/internal/db"
	"log"

	"github.com/Syfaro/telegram-bot-api"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gorm.io/gorm"
)


func CreateMonthKeyboard() tgbotapi.InlineKeyboardMarkup {
    months := []string{"January", "February", "March", "April"}
    var buttons []tgbotapi.InlineKeyboardButton
    for _, month := range months {
        button := tgbotapi.NewInlineKeyboardButtonData(month, month)
        buttons = append(buttons, button)
    }
    keyboard := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(buttons...),
    )
    return keyboard
}

func CreateInfoKeyboard() tgbotapi.InlineKeyboardMarkup {
    info := []string{"Income", "Expense", "Profit", "Tax"}
    var buttons []tgbotapi.InlineKeyboardButton
    for _, item := range info {
        button := tgbotapi.NewInlineKeyboardButtonData(item, item)
        buttons = append(buttons, button)
    }
    keyboard := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(buttons...),
    )
    return keyboard
}

func IsMonth(data string) bool {
	months := []string{"January", "February", "March", "April"}
	for _, month := range months {
		if data == month {
			return true
		}
	}
	return false
}

func FetchFinancialData(month string) (*db.FinancialData, error) {
	var data db.FinancialData

	result := db.GetDB().Where("month = ?", month).First(&data)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No record found for month: %s", month)
			return nil, errors.New("record not found")
		}
		log.Printf("Database error: %s", result.Error)
		return nil, result.Error
	}

	return &data, nil
}

func CreateChart(data *db.FinancialData) (string, error) {
    // Создаем новый график
    p := plot.New()
    p.Title.Text = "Financial Data"
    p.X.Label.Text = "Category"
    p.Y.Label.Text = "Value"

    // Определяем категории и значения
    categories := []string{"Income", "Expense", "Profit", "Tax"}
    values := plotter.Values{
        data.Income,
        data.Expense,
        data.Profit,
        data.Tax,
    }

    // Создаем столбиковую диаграмму
    bars, err := plotter.NewBarChart(values, vg.Points(20))
    if err != nil {
        return "", err
    }
    p.Add(bars)
    p.NominalX(categories...)

    // Сохраняем график в файл
    filePath := "chart.png"
    if err := p.Save(6*vg.Inch, 4*vg.Inch, filePath); err != nil {
        return "", err
    }
    return filePath, nil
}