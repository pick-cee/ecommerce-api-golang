package products

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {

	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	json.Write(w, http.StatusOK, products)
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var tempProduct createProduct

	if err := json.Read(r, &tempProduct); err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusBadRequest, err, err.Error())
		return
	}
	createdProduct, err := h.service.CreateProduct(r.Context(), tempProduct)

	if err != nil {
		if err == InvalidProductRequest {
			json.WriteError(w, http.StatusBadRequest, err, err.Error())
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, createdProduct)
}

func (h *handler) UpdateProductQuantity(w http.ResponseWriter, r *http.Request) {
	var idStr = chi.URLParam(r, "id")
	productId, err := strconv.Atoi(idStr)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	var quantity updateQuantity

	if err := json.Read(r, &quantity); err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	updatedProduct, err := h.service.UpdateProductQuantity(r.Context(), productId, quantity.Quantity)
	if err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	json.Write(w, http.StatusOK, updatedProduct)
}

func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var idStr = chi.URLParam(r, "id")
	productId, err := strconv.Atoi(idStr)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	err = h.service.DeleteProduct(r.Context(), productId)
	if err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	json.Write(w, http.StatusOK, "Product deleted successfully")
}