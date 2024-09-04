package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := "user=postgres password=1234 dbname=tz sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&FinancialData{})
	if err != nil {
		log.Panic("Failed to migrate database:", err)
	}
	data := []FinancialData{
		{Month: "January", Income: 120000, Expense: 80000, Profit: 40000, Tax: 10000},
		{Month: "February", Income: 130000, Expense: 85000, Profit: 45000, Tax: 11250},
		{Month: "March", Income: 140000, Expense: 90000, Profit: 50000, Tax: 12500},
		{Month: "April", Income: 150000, Expense: 95000, Profit: 55000, Tax: 13750},
	}

	for _, entry := range data {
		if !recordExists(DB, entry.Month) {
			result := DB.Create(&entry)
			if result.Error != nil {
				log.Println("Error inserting data:", result.Error)
			}
		}
	}

	log.Println("Data successfully inserted.")

}
func GetDB() *gorm.DB {
    return DB
}

func recordExists(db *gorm.DB, month string) bool {
	var count int64
	db.Model(&FinancialData{}).Where("month = ?", month).Count(&count)
	return count > 0
}
