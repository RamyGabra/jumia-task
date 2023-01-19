package main

import (
	"errors"
	"fmt"
	"jumia-task/controllers"
	"jumia-task/models"
	"jumia-task/routes"
	"jumia-task/services"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm/logger"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	db := initMySql()

	StockService := services.NewStockService(db)
	StockController := controllers.NewStockController(StockService)

	// stock := models.ProductStock{
	// 	Country: "EG",
	// 	Name:    "Meat Lovers Pizza",
	// 	SKU:     "12345678",
	// 	Stock:   29,
	// }

	// db.Create(stock)

	// Set up router
	r := mux.NewRouter()
	router := routes.NewRouter(r, StockController)
	router.LoadRoutes()

	// Listen on server
	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Println("server running on port 3333")
	}
}

func initMySql() *gorm.DB {

	// TODO: Cleanup

	MYSQL_USER := os.Getenv("MYSQL_USER")
	MYSQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")
	MYSQL_HOST := os.Getenv("MYSQL_HOST")
	MYSQL_PORT := os.Getenv("MYSQL_PORT")
	MYSQL_DATABASE := os.Getenv("MYSQL_DATABASE")

	if MYSQL_DATABASE == "" {
		MYSQL_DATABASE = "Stocks"
	}

	dsn := MYSQL_USER + ":" + MYSQL_PASSWORD + "@tcp(" + MYSQL_HOST + ":" + MYSQL_PORT + ")/" + MYSQL_DATABASE + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.ProductStock{})

	mysqlDB, _ := db.DB()
	mysqlDB.SetMaxIdleConns(10)
	mysqlDB.SetMaxOpenConns(10)
	mysqlDB.SetConnMaxIdleTime(time.Minute * 20)

	return db
}
