package model

const (
	SakuStatusOpened = "Opened"
)

type SakuRaya struct {
	ID               int64  `gorm:"column:ID"`
	UserID           int64  `gorm:"column:USER_ID"`
	Name             string `gorm:"column:NAME"`
	CoreCode         string `gorm:"column:CORE_CODE"`
	SavingType       int    `gorm:"column:SAVING_TYPE"`
	TotalTransaction int    `gorm:"column:TOTAL_TRANSACTION"`
	Status           string `gorm:"column:STATUS"`
}

func (SakuRaya) TableName() string {
	return "_saku_raya"
}
