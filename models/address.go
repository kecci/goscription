package models

import "time"

type Address struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id" validate:"required"`
	AddressTitle    string    `json:"address_title" validate:"required"`
	AddressFull     string    `json:"address_full" validate:"required"`
	DistrictName    string    `json:"district_name"`
	SubdistrictName string    `json:"subdistrict_name"`
	ZipCode         string    `json:"zip_code"`
	Primary         bool      `json:"primary"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedBy       string    `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
}
