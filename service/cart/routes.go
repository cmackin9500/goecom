package cart

import (
	"fmt"
	"net/http"

	"github.com/cmackin9500/goecom/service/auth"
	"github.com/cmackin9500/goecom/types"
	"github.com/cmackin9500/goecom/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	orderStore types.OrderStore
	productStore types.ProductStore
	userStore types.UserStore
}

func NewHandler(productStore types.ProductStore, orderStore types.OrderStore, userStore types.UserStore) *Handler {
	return &Handler{orderStore: orderStore, productStore: productStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handlerCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handlerCheckout(w http.ResponseWriter, r *http.Request) {
	 var cart types.CartCheckoutPayload
	 if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}

	// validate the payload
	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	} 

	// get products
	productIDs, err := getCartItemsIds(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ps, err := h.productStore.GetProductsByIDs(productIDs)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	userID := auth.GetUserIDFromContext(r.Context())
	orderID, totalPrice, err := h.createOrder(ps, cart.Items, userID)

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id": orderID,
	})
}