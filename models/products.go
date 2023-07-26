package models

import (
	"gorm.io/gorm"
)

type Discount string

const (
	Percentage Discount = "PERCENTAGE"
	Flat       Discount = "FLAT"
)

type Products struct {
	gorm.Model
	ShortName          string  `json:"shortName" gorm:"not null" validate:"required,min=3"`
	LongName           string  `json:"longName" gorm:"not null" validate:"required,min=3"`
	Cost               int     `json:"cost"`
	Price              int     `json:"price"`
	DiscountType       *string `json:"discountType"`
	DiscountPrice      *int    `json:"discountPrice"`
	Description        string  `json:"description"`
	DealerName         string  `json:"dealerName"`
	DealerPlace        string  `json:"dealerPlace"`
	ProductDestination string  `json:"productDestination"`
}
