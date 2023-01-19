package routes

import (
	"jumia-task/controllers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Router struct {
	router          *mux.Router
	StockController *controllers.StockController
}

func NewRouter(router *mux.Router, controller *controllers.StockController) *Router {

	return &Router{
		router:          router,
		StockController: controller,
	}
}

func (m *Router) LoadRoutes() {

	m.router.Handle("/stock", negroni.New(
		negroni.Wrap(http.HandlerFunc(m.StockController.GetProductBySku)),
	)).Methods("GET")

	m.router.Handle("/stock", negroni.New(
		negroni.Wrap(http.HandlerFunc(m.StockController.PutConsumeProduct)),
	)).Methods("PUT")

	m.router.Handle("/stock", negroni.New(
		negroni.Wrap(http.HandlerFunc(m.StockController.PostBatchUpdate)),
	)).Methods("POST")

	http.Handle("/", m.router)
}
