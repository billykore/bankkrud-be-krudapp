package model

import "time"

type Transaction struct {
	ID                   int64     `gorm:"column:ID"`
	UUID                 string    `gorm:"column:UUID"`
	UserID               string    `gorm:"column:USER_ID"`
	WalletIDSource       int64     `gorm:"column:WALLET_ID_SOURCE"`
	Destination          string    `gorm:"column:DESTINATION"`
	Amount               string    `gorm:"column:AMOUNT"`
	TransactionType      string    `gorm:"column:TRANSACTION_TYPE"`
	TransactionReference string    `gorm:"column:TRREFN"`
	SequenceJournal      string    `gorm:"column:SEQUENCE_JOURNAL"`
	Remarks              string    `gorm:"column:REMARKS"`
	Note                 string    `gorm:"column:NOTE"`
	Status               string    `gorm:"column:STATUS"`
	FreeFeeId            int       `gorm:"column:FREE_FEE_ID"`
	CreatedAt            time.Time `gorm:"column:CREATED_AT"`
	Fee                  string    `gorm:"column:FEE"`
	DestinationName      string    `gorm:"column:DEST_NAME"`
	InitialSourceBalance float64   `gorm:"column:INITIAL_SOURCE_BALANCE"`
	StatusCode           string    `gorm:"column:STATUS_CODE"`
	SequenceNumber       string    `gorm:"column:SEQ_NO"`
	BankCode             string    `gorm:"column:BANK_CODE"`
	SuccessTrxDate       string    `gorm:"column:SUCCESS_TRX_DATE"`
	Discount             string    `gorm:"column:DISCOUNT"`
	Cashback             string    `gorm:"column:CASHBACK"`
	RewardClaimId        string    `gorm:"column:REWARD_CLAIM_ID"`
	Currency             string    `gorm:"column:CURRENCY"`
	AdditionalData       string    `gorm:"column:ADDITIONAL_DATA"`
	DestinationSakuId    string    `gorm:"column:DESTINATION_SAKU_ID"`
	DestinationSakuType  string    `gorm:"column:DEST_SAKU_TYPE"`
}
