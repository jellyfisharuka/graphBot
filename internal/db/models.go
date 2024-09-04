package db

import "gorm.io/gorm"

type FinancialData struct {
	gorm.Model
	Month   string  `gorm:"type:varchar(50);not null"`
	Income  float64 `gorm:"not null"`
	Expense float64 `gorm:"not null"`
	Profit  float64 `gorm:"not null"`
	Tax     float64 `gorm:"not null"`
}