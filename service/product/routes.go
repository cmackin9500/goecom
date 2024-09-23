package product

import (
	"fmt"
	"net/http"

	"github.com/cmackin9500/goecom/types"
	"github.com/cmackin9500/goecom/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handlerCreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/products", h.handlerGetProduct).Methods(http.MethodGet)
}

func (h *Handler) handlerCreateProduct(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterProductPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	} 
	
	// check if the product exists
	_, err := h.store.GetProductByName(payload.Name)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("product with name %s already exists", payload.Name))
		return
	}

	// if it doesn't, we create the new product
	err = h.store.CreateProduct(types.Product{		
		Name: 			payload.Name,
		Description: 	payload.Description,
		Image:			payload.Image,
		Price: 			payload.Price,
		Quantity: 		payload.Quantity,
	})
}

func (h *Handler) handlerGetProduct(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJSON(w, http.StatusOK, ps)
}