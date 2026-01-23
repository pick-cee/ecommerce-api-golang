package orders

import (
	"log"
	"net/http"

	"github.com/pick-cee/go-ecommerce-api/internal/json"
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

	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		log.Println(err)

		if err == ProductNotFoundError {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		if err == ProductNoStockError {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.Write(w, http.StatusCreated, createdOrder)
}