package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	TransactionId        uuid.UUID       `gorm:"type:uuid;index;primaryKey;not null" json:"transactionId"`
	SourceAccountId      int64           `gorm:"not null" json:"sourceAccountId"`
	DestinationAccountId int64           `gorm:"not null" json:"destinationAccountId"`
	Amount               decimal.Decimal `gorm:"type:decimal(10,5);default:0.00" json:"amount"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
