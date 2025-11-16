package account_service

import (
	"testing"

	"internal-transfers/dtos"

	"github.com/shopspring/decimal"
)

func TestTransferFunds_SelfTransfer(t *testing.T) {
	body := &dtos.CreateTransactionDto{
		SourceAccountId:      123,
		DestinationAccountId: 123,
		Amount:               decimal.NewFromFloat(100.0),
	}

	err := TransferFunds(body)
	if err == nil {
		t.Error("Expected error for self-transfer, got nil")
	}
	if err.Error() != "Cannot transfer funds to the same account" {
		t.Errorf("Expected 'Cannot transfer funds to the same account' error, got: %v", err)
	}
}

func TestTransferFunds_ZeroAmount(t *testing.T) {
	body := &dtos.CreateTransactionDto{
		SourceAccountId:      123,
		DestinationAccountId: 456,
		Amount:               decimal.Zero,
	}

	err := TransferFunds(body)
	if err == nil {
		t.Error("Expected error for zero amount, got nil")
	}
	if err.Error() != "Transfer amount must be greater than zero" {
		t.Errorf("Expected 'Transfer amount must be greater than zero' error, got: %v", err)
	}
}

func TestTransferFunds_NegativeAmount(t *testing.T) {
	body := &dtos.CreateTransactionDto{
		SourceAccountId:      123,
		DestinationAccountId: 456,
		Amount:               decimal.NewFromFloat(-50.0),
	}

	err := TransferFunds(body)
	if err == nil {
		t.Error("Expected error for negative amount, got nil")
	}
	if err.Error() != "Transfer amount must be greater than zero" {
		t.Errorf("Expected 'Transfer amount must be greater than zero' error, got: %v", err)
	}
}

func TestCreateAccount_NegativeBalance(t *testing.T) {
	body := &dtos.CreateAccountDto{
		AccountId: 123,
		Balance:   decimal.NewFromFloat(-100.0),
	}

	err := CreateAccount(body)
	if err == nil {
		t.Error("Expected error for negative balance, got nil")
	}
	if err.Error() != "Initial balance cannot be negative" {
		t.Errorf("Expected 'Initial balance cannot be negative' error, got: %v", err)
	}
}
