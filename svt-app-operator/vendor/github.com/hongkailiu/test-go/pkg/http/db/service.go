package db

import (
	"github.com/hongkailiu/test-go/pkg/http/model"
	"github.com/jinzhu/gorm"
)

// Service provides all functions with db
type Service struct {
	db *gorm.DB
}

// New constructs a db service
func New(db *gorm.DB) *Service {
	return &Service{db}
}

// GetCities get cities
func (s *Service) GetCities(limit, offset int, cities *[]model.City) error {
	return s.db.Limit(limit).Offset(offset).Find(cities).Error
}
