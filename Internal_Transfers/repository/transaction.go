package repository

import (
	database "internal-transfers/db"
	"internal-transfers/db/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func CreateTransaction(sourceAccountId int64, destinationAccountId int64, amount decimal.Decimal) error {
	db, err := database.GetDb()
	if err != nil {
		zap.L().Error("Could not connect to db", zap.Error(err))
		return  fiber.ErrInternalServerError
	}

	transactionId := uuid.New()
	transaction := &models.Transaction {
		TransactionId: transactionId,
		SourceAccountId: sourceAccountId,
		DestinationAccountId: destinationAccountId,
		Amount: amount,
	}

	result := db.Create(transaction)
	if err := result.Error; err != nil {
		zap.L().Error("Error occurred while trying to create transaction", zap.Error(err))
		if database.IsDuplicateKeyError(err) {
			return fiber.ErrConflict
		}
		return fiber.ErrInternalServerError
	}
	zap.L().Info("Transaction created successfully", zap.Any("transaction", transaction))
	return nil
}