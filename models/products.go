package models

import (
	"time"
)

const (
	Percentage string = "PERCENTAGE"
	Flat       string = "FLAT"
)

type Products struct {
	Id                 int       `json:"id" gorm:"primaryKey"`
	ShortName          string    `json:"shortName" gorm:"not null" validate:"required,min=3"`
	LongName           string    `json:"longName" gorm:"not null" validate:"required,min=3"`
	Cost               uint      `json:"cost" gorm:"not null"`
	Price              uint      `json:"price" gorm:"not null"`
	FinalPrice         uint      `json:"finalPrice" gorm:"not null"`
	IsDiscount         bool      `json:"isDiscount" gorm:"not null"`
	DiscountType       *string   `json:"discountType"`
	DiscountPrice      *int      `json:"discountPrice"`
	Description        string    `json:"description" gorm:"not null"`
	Stock              uint      `json:"stock" gorm:"not null"`
	DealerName         string    `json:"dealerName"`
	DealerPlace        string    `json:"dealerPlace"`
	ProductDestination string    `json:"productDestination"`
	HeroImage          string    `json:"heroImage" gorm:"not null"`
	CreatedAt          time.Time `json:"createdAt"`
}
