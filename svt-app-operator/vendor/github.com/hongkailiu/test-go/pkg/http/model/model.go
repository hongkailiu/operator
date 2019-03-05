package model

import "github.com/jinzhu/gorm"

// Order represents an order
type Order struct {
	gorm.Model
	Name   string
	CityID int
}

// City represents a city
type City struct {
	gorm.Model
	Name string
}
