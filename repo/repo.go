package repo

import (
	"errors"
	"fmt"
	"jumia-task/constants"
	custom_errors "jumia-task/errors"
	"jumia-task/models"
	"strings"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) FindProductStock(country string, sku string) (*models.ProductStock, *custom_errors.CustomError) {

	var stock models.ProductStock

	result := r.db.Where("country = ? And sku = ?", country, sku).First(&stock)
	err := result.Error
	if err != nil {
		err := custom_errors.NewInternalServerError(err, constants.StockRepo, constants.RepoFindProductStock)
		return nil, err
	}

	return &stock, nil
}

func (r *Repo) ConsumeProduct(country string, sku string) *custom_errors.CustomError {

	// decrement stock by 1 if [SKU & country FOUND AND stock > 0]
	result := r.db.Exec(`UPDATE product_stocks SET stock = stock - 1 WHERE stock > 0 and sku = ? and country = ?`, sku, country)
	err := result.Error
	if err != nil {
		return custom_errors.NewInternalServerError(err, constants.StockRepo, constants.RepoConsumeProduct)
	}

	if result.RowsAffected == 0 {
		errString := fmt.Sprintf("no stock of requested product with sku and country: %s & %s respectively", sku, country)
		err := errors.New(errString)
		return custom_errors.NewBadRequestError(err, errString, constants.StockRepo, constants.RepoConsumeProduct)
	}

	return nil

}

func (r *Repo) UpdateBatch(records []*models.ProductStock) *custom_errors.CustomError {

	group := new(errgroup.Group)

	// execute all insertion queries within the same transaction
	// to rollback if an error occurs
	err := r.db.Transaction(func(tx *gorm.DB) error {
		start := 0
		end := 0
		sliceSize := 100

		// process the records in slices to leverage conccurency
		for start < len(records) {

			end = end + sliceSize
			if end > len(records) {
				end = len(records)
			}

			slice := records[start:end]

			group.Go(func() error {
				return updateSlice(tx, slice)
			})

			start = start + sliceSize
			if start > len(records) {
				break
			}
		}

		// wait for all go routines to finish
		err := group.Wait()
		return err
	})

	// return custom error
	if err != nil {
		return custom_errors.NewInternalServerError(err, constants.StockRepo, constants.RepoUpdateBatch)
	}

	return nil
}

func updateSlice(tx *gorm.DB, records []*models.ProductStock) error {

	query := "INSERT INTO product_stocks (country, sku, name, stock) VALUES"
	var values []string

	for _, element := range records {
		values = append(values, fmt.Sprintf(" ('%s', '%s', '%s', %d)", element.Country, element.SKU, element.Name, element.Stock))
	}

	query = query + strings.Join(values, ",")
	query = query + " ON DUPLICATE KEY UPDATE stock=stock+values(stock);"

	result := tx.Exec(query)
	err := result.Error
	if err != nil {
		return err
	}

	return nil
}
