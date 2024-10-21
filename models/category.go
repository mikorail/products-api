package models

type Category struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"size:100;not null"`
	Description string `json:"description"`
}
