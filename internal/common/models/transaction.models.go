package models

import (
	"time"

	"gorm.io/gorm"
)

// --- Enums & Constants ---

type TransactionType string

const (
	TypeIncome  TransactionType = "INCOME"
	TypeExpense TransactionType = "EXPENSE"
)

// --- Structs ---

type User struct {
	gorm.Model
	PhoneNumber string `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone_number"`
	Name        string `gorm:"type:varchar(100)" json:"name"`

	// Relations (Has Many)
	Budgets      []Budget      `gorm:"foreignKey:UserID" json:"budgets,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
}

type Category struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);not null" json:"name"`

	// Relations
	Budgets      []Budget      `gorm:"foreignKey:CategoryID" json:"budgets,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:CategoryID" json:"transactions,omitempty"`
}

type Budget struct {
	gorm.Model
	UserID     uint    `gorm:"not null;index" json:"user_id"`
	CategoryID uint    `gorm:"not null;index" json:"category_id"`
	Amount     float64 `gorm:"type:decimal(15,2);not null" json:"amount"` // Menggunakan decimal untuk uang
	Month      int     `gorm:"not null" json:"month"`                     // 1-12
	Year       int     `gorm:"not null" json:"year"`                      // 2023, 2024

	// Relations
	User     User     `json:"-"`
	Category Category `json:"category"`
}

type Transaction struct {
	gorm.Model
	UserID          uint            `gorm:"not null;index" json:"user_id"`
	CategoryID      *uint           `gorm:"index" json:"category_id"`
	Amount          float64         `gorm:"type:decimal(15,2);not null" json:"amount"`
	Description     string          `gorm:"type:text" json:"description"`
	TransactionDate time.Time       `gorm:"not null" json:"transaction_date"`
	Type            TransactionType `gorm:"type:varchar(10);not null" json:"type"` // INCOME / EXPENSE

	// Relations
	User     User     `json:"-"`
	Category Category `json:"category"`
}
