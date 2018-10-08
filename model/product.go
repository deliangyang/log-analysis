package model

type Product struct {
	Model
	Title string `'gorm:"default:''"'`
	Amount int64
}

type StatProductView struct {
	ProductId uint
	Date string
	PV uint
	UV uint
}

type ProductViewedDetail struct {
	Model
	ProductId uint
	UserId uint
}

func (Product) TableName() string {
	return "products"
}

