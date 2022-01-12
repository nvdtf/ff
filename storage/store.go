package storage

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Provider interface {
	Save(tx *Transaction) error
	EnableAutoCleanup(interval time.Duration, maxAge time.Duration)
}

type Storage struct {
	db *gorm.DB
}

type Transaction struct {
	gorm.Model

	Authorizers string
	Tx          string
	Code        string
	Error       string
	ImportTags  string
	Events      string
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

func NewPostgresStorage(dsn string) (Provider, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

func (s *Storage) EnableAutoCleanup(interval time.Duration, maxAge time.Duration) {
	go func() {
		for {
			err := s.db.Unscoped().Delete(&Transaction{},
				"now()-created_at > ?",
				maxAge).Error
			if err != nil {
				fmt.Printf("Cleanup failed: %s\n", err)
			}
			time.Sleep(interval)
		}
	}()
}
