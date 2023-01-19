package services

import (
	custom_errors "jumia-task/errors"
	"jumia-task/models"
	"jumia-task/repo"

	"gorm.io/gorm"
)

type StockService struct {
	repo *repo.Repo
}

func NewStockService(db *gorm.DB) *StockService {
	stockRepo := repo.NewRepo(db)
	return &StockService{
		repo: stockRepo,
	}
}

func (s *StockService) FindProductBySKU(country string, sku string) (*models.ProductStock, *custom_errors.CustomError) {

	stock, err := s.repo.FindProductStock(country, sku)
	if err != nil {
		return nil, err
	}

	return stock, nil
}

func (s *StockService) ConsumeProduct(country string, sku string) *custom_errors.CustomError {

	err := s.repo.ConsumeProduct(country, sku)
	if err != nil {
		return err
	}

	return nil
}

func (s *StockService) BatchUpdate(records []*models.ProductStock) *custom_errors.CustomError {

	err := s.repo.UpdateBatch(records)
	if err != nil {
		return err
	}

	return nil
}
