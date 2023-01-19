package models

type ProductStock struct {
	Country string `gorm:"primaryKey; not null; autoIncrement:false" json:"country"`
	SKU     string `gorm:"primaryKey; not null; autoIncrement:false" json:"sku"`
	Name    string `json:"name"`
	Stock   int    `gorm:"default:0" json:"stock"`
}
