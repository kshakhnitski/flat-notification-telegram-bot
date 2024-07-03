package model

type Flat struct {
	ID          string  `gorm:"primaryKey"`
	Address     string  `gorm:"type:varchar(255)"`
	Description string  `gorm:"type:text"`
	PriceInByn  float64 `gorm:"type:numeric(10,2)"`
	PriceInUsd  float64 `gorm:"type:numeric(10,2)"`
	Parameters  string  `gorm:"type:varchar(255)"`
	Metro       string  `gorm:"type:varchar(255)"`
	Link        string  `gorm:"type:varchar(255)"`
	Source      string  `gorm:"type:varchar(255)"`
}
