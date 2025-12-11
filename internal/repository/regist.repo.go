package repository

import (
	"pannypal/internal/repository/analytics"
	"pannypal/internal/repository/bot"
	"pannypal/internal/repository/budget"
	"pannypal/internal/repository/category"
	logdata "pannypal/internal/repository/log-data"
	"pannypal/internal/repository/transaction"
	"pannypal/internal/repository/user"
)

type IRepository struct {
	User        user.IRepository
	Category    category.IRepository
	Budget      budget.IRepository
	Transaction transaction.IRepository
	Analytics   analytics.IRepository
	LogData     logdata.IRepository
	Bot         bot.IRepository
}
