package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Provider interface {
	Save(tx *Transaction) error
}

type Storage struct {
	db *gorm.DB
}

type Transaction struct {
	gorm.Model

	Tx    string
	Code  string
	Error string
	Tags  string
}

func NewSqliteStorage() (Provider, error) {
	db, err := gorm.Open(sqlite.Open("flow.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Transaction{})
	if err != nil {
		return nil, err
	}
	return &Storage{db}, nil
}

func (s *Storage) Save(tx *Transaction) error {
	return s.db.Create(tx).Error
}
