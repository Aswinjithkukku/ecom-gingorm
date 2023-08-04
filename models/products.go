package models

import (
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
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
	Slug               string    `json:"slug" gorm:"unique"`
	CreatedAt          time.Time `json:"createdAt"`
}

func (p *Products) BeforeSave(tx *gorm.DB) error {
	p.Slug = slug.Make(p.LongName)

	// Check whether the slug is unique or not.
	var count int64
	tx.Model(&Products{}).Where("slug=?", p.Slug).Count(&count)
	if count > 0 {
		return fmt.Errorf("Slug is not unique. Product with same long name exist. Provide another longName")
	}
	return nil
}
