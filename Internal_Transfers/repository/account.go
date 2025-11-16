package repository

import (
	"errors"

	database "internal-transfers/db"
	"internal-transfers/db/models"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateAccount(accountId int64, balance decimal.Decimal) error {
	db, err := database.GetDb()
	if err != nil {
		zap.L().Error("Failed to get database connection", zap.Error(err))
		return  fiber.ErrInternalServerError
	}
	account := &models.Account{
		AccountId: accountId,
		Balance: balance,
	}

	result := db.Create(account)
	if err := result.Error; err != nil  {
		zap.L().Error("Error occurred while trying to insert record in account table", zap.Error(err))
		if database.IsDuplicateKeyError(err) {
			return fiber.ErrConflict
		}
		return fiber.ErrInternalServerError
	}
	return nil
}

func GetAccount(accountId int64) (*models.Account, error) {
	db, err := database.GetDb()
	if err != nil {
		zap.L().Error("Could not connect to db", zap.Error(err), zap.Any("accountId", accountId))
		return nil, fiber.ErrInternalServerError
	}
	account := new(models.Account)
	err = db.Where("account_id = ?", accountId).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		return nil, fiber.ErrInternalServerError
	}
	return account, nil
}

func DebitFunds(accountId int64, amount decimal.Decimal) error {
	db, err := database.GetDb()
	if err != nil {
		zap.L().Error("Could not connect to db", zap.Error(err), zap.Any("accountId", accountId))
		return fiber.ErrInternalServerError
	}
	results := db.Clauses(clause.Locking{ // we will block other queries from updating or deleting the row
		Strength: "SHARE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).
		Model(&models.Account{}).Where("account_id = ?", accountId).Update("balance",
		gorm.Expr("balance - ?", amount))
	return results.Error
}

func CreditFunds(accountId int64, amount decimal.Decimal) error {
	db, err := database.GetDb()
	if err != nil {
		zap.L().Error("Could not connect to db", zap.Error(err), zap.Any("accountId", accountId))
		return fiber.ErrInternalServerError
	}
	results := db.Clauses(clause.Locking{ // we will block other queries from updating or deleting the row
		Strength: "SHARE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).
		Model(&models.Account{}).Where("account_id = ?", accountId).Update("balance",
		gorm.Expr("balance + ?", amount))
	return results.Error
}