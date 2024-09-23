package order

import (
	"net/http"

	"github.com/cmackin9500/goecom/types"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.OrderStore
}

func NewHandler(store types.OrderStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", h.handlerCheckout).Methods(http.MethodPost)
}

func (h *Handler) handlerCheckout(w http.ResponseWriter, r *http.Request) {

}