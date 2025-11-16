package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	AccountId int64           `gorm:"primaryKey;autoIncrement:false" json:"accountId"`
	Balance   decimal.Decimal `gorm:"type:decimal(10,5);default:0.00" json:"balance"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
