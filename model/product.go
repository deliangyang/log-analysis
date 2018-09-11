package model

type Product struct {
	ID           string `gorm:"AUTO_INCREMENT;primary_key"`
	LanguageCode string
	Code         string
	Name         string
}
