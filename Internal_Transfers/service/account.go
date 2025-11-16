package account_service

import (
	database "internal-transfers/db"
	"internal-transfers/dtos"
	"internal-transfers/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateAccount(body *dtos.CreateAccountDto) error {
	// Validate: to prevent negative initial balance
	if body.Balance.LessThan(decimal.Zero) {
		zap.L().Error("Invalid initial balance", zap.String("balance", body.Balance.String()))
		return fiber.NewError(fiber.ErrBadRequest.Code, "Initial balance cannot be negative")
	}
	return repository.CreateAccount(body.AccountId, body.Balance)
}

func GetAccount(accountId int64) (*dtos.GetAccountResponseDto, error) {
	account, err := repository.GetAccount(accountId)
	if err != nil {
		return nil, err
	}
	return &dtos.GetAccountResponseDto{
		AccountId: account.AccountId,
		Balance:   account.Balance,
	}, nil
}

func TransferFunds(body *dtos.CreateTransactionDto) error {
	// Validate: to prevent self-transfer
	if body.SourceAccountId == body.DestinationAccountId {
		zap.L().Error("Cannot transfer to same account", zap.Int64("accountId", body.SourceAccountId))
		return fiber.NewError(fiber.ErrBadRequest.Code, "Cannot transfer funds to the same account")
	}

	// Validate: to prevent negative or zero amounts
	if body.Amount.LessThanOrEqual(decimal.Zero) {
		zap.L().Error("Invalid amount", zap.String("amount", body.Amount.String()))
		return fiber.NewError(fiber.ErrBadRequest.Code, "Transfer amount must be greater than zero")
	}

	db, err := database.GetDb()
	if err != nil {
		zap.L().Error("Could not connect to db", zap.Error(err), zap.Any("body", body))
		return fiber.ErrInternalServerError
	}
	sourceAccount, err := repository.GetAccount(body.SourceAccountId)
	if err != nil {
		zap.L().Error("Error fetching source account", zap.Error(err), zap.Int64("sourceAccountId", body.SourceAccountId))
		return err
	}
	if sourceAccount.Balance.Cmp(body.Amount) == -1 { // check if enough balance in account
		zap.L().Error("insufficient balance in source acccount")
		return fiber.NewError(fiber.ErrBadRequest.Code, "Insufficient Balance")
	}

	// fetch destination account
	_, err = repository.GetAccount(body.DestinationAccountId)
	if err != nil {
		zap.L().Error("Error fetching destination account", zap.Error(err), zap.Int64("destinationAccountId", body.DestinationAccountId))
		return err
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		if err = repository.CreateTransaction(body.SourceAccountId, body.DestinationAccountId, body.Amount); err != nil {
			return err
		}
		if err = repository.DebitFunds(body.SourceAccountId, body.Amount); err != nil {
			return err
		}

		if err = repository.CreditFunds(body.DestinationAccountId, body.Amount); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		zap.L().Error("Transaction failed", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	return nil
}
