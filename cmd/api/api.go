package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/cmackin9500/goecom/service/cart"
	"github.com/cmackin9500/goecom/service/order"
	"github.com/cmackin9500/goecom/service/product"
	"github.com/cmackin9500/goecom/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	address string
	db *sql.DB
}

func NewAPIServer(address string, db *sql.DB) *APIServer {
	return &APIServer{
		address: address,
		db: db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// USER
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)
	// PRODUCTS
	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)
	// ORDER
	orderStore := order.NewStore(s.db)
	orderHandler := order.NewHandler(orderStore)
	orderHandler.RegisterRoutes(subrouter)
	// CARTS
	cartHandler := cart.NewHandler(productStore, orderStore, userStore)
	cartHandler.RegisterRoutes(subrouter)
	
	log.Println("Listening on", s.address)
	return http.ListenAndServe(s.address, router)
}