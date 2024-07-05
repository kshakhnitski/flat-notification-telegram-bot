package model

type Flat struct {
	ID          string   `gorm:"primaryKey"`
	Address     string   `gorm:"type:varchar(255)"`
	Description string   `gorm:"type:text"`
	PriceInByn  float64  `gorm:"type:numeric(10,2)"`
	PriceInUsd  float64  `gorm:"type:numeric(10,2)"`
	Parameters  string   `gorm:"type:varchar(255)"`
	Rooms       int      `gorm:"type:integer"`
	Area        *float64 `gorm:"type:numeric(10,2);null"`
	Floor       *int     `gorm:"type:integer;null"`
	TotalFloors *int     `gorm:"type:integer;null"`
	Metro       string   `gorm:"type:varchar(255)"`
	Link        string   `gorm:"type:varchar(255)"`
	Source      string   `gorm:"type:varchar(255)"`
}
