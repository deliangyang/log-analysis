package model

import "time"

type Api struct {
	ID uint `gorm:"primary_key;AUTO_INCREMENT"`
	Path string `gorm:"type:varchar(60);"`
	Method string `gorm:"size:10"`
	Remarks string `gorm:"size:40"`
	CreatedAt time.Time
}

type ApiLog struct {
	ID uint  `gorm:"primary_key"`
	UserId uint
	Path string `gorm:"type:varchar(60);"`
	Method string `gorm:"size:10"`
	Ip string `gorm:"size:20"`
	CreatedAt time.Time
}
