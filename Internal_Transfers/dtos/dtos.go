package dtos

import (
	"github.com/shopspring/decimal"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

type Response struct {
	Status  int64       `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CreateAccountDto struct {
	AccountId int64           `json:"account_id" validate:"required,number"`
	Balance   decimal.Decimal `json:"initial_balance" validate:"required"`
}

type GetAccountResponseDto struct {
	AccountId int64           `json:"account_id"`
	Balance   decimal.Decimal `json:"balance"`
}

type CreateTransactionDto struct {
	SourceAccountId      int64           `json:"source_account_id" validate:"required,number"`
	DestinationAccountId int64           `json:"destination_account_id" validate:"required,number"`
	Amount               decimal.Decimal `json:"amount" validate:"required"`
}