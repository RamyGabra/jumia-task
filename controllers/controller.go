package controllers

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"jumia-task/constants"
	"jumia-task/models"
	"jumia-task/services"
	"net/http"
	"strconv"
	"strings"

	custom_errors "jumia-task/errors"

	"github.com/sirupsen/logrus"
)

type StockController struct {
	StockService *services.StockService
}

type HttpErrorResponse struct {
	Description string `json:"description"`
	StatusCode  int    `json:"status_code"`
}

func NewStockController(service *services.StockService) *StockController {

	return &StockController{
		StockService: service,
	}
}

func (S *StockController) GetProductBySku(w http.ResponseWriter, r *http.Request) {

	uri := r.URL.String()
	method := r.Method

	logrus.WithFields(logrus.Fields{
		"service": "StockController",
	}).Infof("%s %s", method, uri)

	country := r.URL.Query().Get("country")
	sku := r.URL.Query().Get("sku")

	if country == "" || sku == "" {
		errString := "Error GetProductBySku - missing params"
		err := custom_errors.NewBadRequestError(errors.New(errString), errString, constants.StockController, constants.GetProductBySku)
		err.Log()
		http.Error(w, err.Error(), err.StatusCode)
		return
	}

	stock, err := S.StockService.FindProductBySKU(country, sku)
	if err != nil {
		err.Log()
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	response, marshalErorr := json.Marshal(stock)
	if marshalErorr != nil {
		err := custom_errors.NewInternalServerError(marshalErorr, constants.StockController, constants.GetProductBySku)
		err.Log()
		http.Error(w, err.Error(), err.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (S *StockController) PutConsumeProduct(w http.ResponseWriter, r *http.Request) {

	uri := r.URL.String()
	method := r.Method

	logrus.WithFields(logrus.Fields{
		"service": "StockController",
	}).Infof("%s %s", method, uri)

	country := r.URL.Query().Get("country")
	sku := r.URL.Query().Get("sku")

	if country == "" || sku == "" {
		errString := fmt.Sprintf("Error %s - missing params", constants.PutConsumerProduct)
		err := custom_errors.NewBadRequestError(errors.New(errString), errString, constants.StockController, constants.PutConsumerProduct)
		err.Log()
		http.Error(w, err.Error(), err.StatusCode)
		return
	}

	err := S.StockService.ConsumeProduct(country, sku)
	if err != nil {
		err.Log()
		http.Error(w, err.Error(), err.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func (S *StockController) PostBatchUpdate(w http.ResponseWriter, r *http.Request) {

	var records []*models.ProductStock

	uri := r.URL.String()
	method := r.Method

	logrus.WithFields(logrus.Fields{
		"service": "StockController",
	}).Infof("%s %s", method, uri)

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024) // 10MB
	file, _, err := r.FormFile("file")
	defer file.Close()

	if err != nil {
		err := custom_errors.NewInternalServerError(err, constants.StockController, constants.PostBatchUpdate)
		err.Log()
		http.Error(w, err.Error(), err.StatusCode)
	}

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			err := custom_errors.NewBadRequestError(err, "error reading csv file - invalid csv", constants.StockController, constants.PostBatchUpdate)
			err.Log()
			http.Error(w, err.Error(), err.StatusCode)
			break
		}

		// skip column headers
		if record[0] == "country" {
			continue
		}

		stockChange, err := strconv.Atoi(record[3])
		if err != nil {
			err = errors.New(fmt.Sprintf("error - invalid data in csv file - cannot convert %s to int", record[3]))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if strings.Contains(record[0], "'") || strings.Contains(record[1], "'") || strings.Contains(record[2], "'") {
			err = errors.New("error - invalid data in csv file - character ' is not allowed")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		stock := &models.ProductStock{
			Country: record[0],
			SKU:     record[1],
			Name:    record[2],
			Stock:   stockChange, // note stock here is change in stock
		}

		records = append(records, stock)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	serviceError := S.StockService.BatchUpdate(records)
	if err != nil {
		serviceError.Log()
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// fmt.Println(records)

	w.WriteHeader(200)
	w.Write(nil)
}
