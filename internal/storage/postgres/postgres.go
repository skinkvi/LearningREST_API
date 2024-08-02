package postgres

import (
	"fmt"
	"log"
	"rest_api_app/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func (s *Storage) AutoMigrationTablePg() error {
	return s.db.AutoMigrate(&models.URL{})
}

func New() (*Storage, error) {
	const op = "storage.postgres.New" //

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=postgres dbname=urlshortdb port=5433 sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgres.SaveURL"

	urlModel := &models.URL{
		Url:   urlToSave,
		Alias: alias,
	}

	res := s.db.Create(&urlModel)

	if res.Error != nil {
		log.Fatalf("falied to save url: %s: %v", op, res.Error)
		return 0, res.Error
	}

	return int64(urlModel.ID), nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	var urlModel models.URL

	res := s.db.First(&urlModel, "alias = ?", alias)

	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("URL not found for alias: %s", alias)
		}
		log.Fatalf("falied to get url: %s: %v", op, res.Error)
		return "", res.Error
	}

	return urlModel.Url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"

	var urlModel models.URL

	res := s.db.Delete(&urlModel, "alias = ?", alias)

	if res.Error != nil {
		if res.Error == gorm.ErrInvalidData {
			return res.Error
		}
		log.Fatalf("falied delete url for alias: %s: %v", op, alias)
		return res.Error
	}

	return nil
}
