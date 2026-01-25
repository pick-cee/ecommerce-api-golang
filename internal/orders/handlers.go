package orders

import (
	"errors"
	"log"
	"net/http"

	"github.com/pick-cee/go-ecommerce-api/internal/json"
	"github.com/pick-cee/go-ecommerce-api/internal/users"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,	
	}
}
 
func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var tempOrder createOrderRequest

	if err := json.Read(r, &tempOrder); err != nil{
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := users.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
	}
	
	tempOrder.CustomerID = user.ID
	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		log.Println(err)

		if errors.Is(err, ProductNotFoundError) {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		if errors.Is(err, ProductNoStockError) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.Write(w, http.StatusCreated, createdOrder)
}