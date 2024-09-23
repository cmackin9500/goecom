package cart

import (
	"fmt"

	"github.com/cmackin9500/goecom/types"
)


func getCartItemsIds(items []types.CartCheckoutItem) ([]int, error) {
	produdctIds := make([]int, len(items))
	for i, item := range(items) {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductID)
		}
		produdctIds[i] = item.ProductID
	}
	return produdctIds, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartCheckoutItem, userID int) (int, float64, error) {
	productMap :=  make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}

	// check if all products are actually in stock
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}
	// calculate the total price
	totalPrice := calculateTotalPrice(items, productMap)

	// reduce quantity of products in our db
	for _, item := range(items) {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity
		h.productStore.UpdateProduct(product)
	}

	// create the order
	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID: userID,
		TotalPrice: totalPrice,
		Status: "pending",
		Address: "some address",
	})
	if err != nil {
		return 0, 0, err
	}

	// create order items
	for _, item := range items {
		h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID: orderID,
			ProductID: item.ProductID,
			Quantity: item.Quantity,
			Price: productMap[item.ProductID].Price,
		})
	}
	
	return orderID, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []types.CartCheckoutItem, productMap map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}
	
	for _, checkoutItem := range cartItems {
		p, ok := productMap[checkoutItem.ProductID]

		if !ok {
			return fmt.Errorf("product %d is not available", checkoutItem.ProductID)
		}

		if p.Quantity < checkoutItem.Quantity {
			return fmt.Errorf("not enough %s in stock", p.Name)
		}
	}
	return nil
}

func calculateTotalPrice(cartItems []types.CartCheckoutItem, productMap map[int]types.Product) float64 {
	var totalPrice float64 = 0.0 

	for _, checkoutItem := range cartItems {
		p := productMap[checkoutItem.ProductID]
		totalPrice += (float64(checkoutItem.Quantity)*p.Price)
	}

	return totalPrice
}