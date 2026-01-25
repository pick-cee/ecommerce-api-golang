package users

import (
	"log"
	"net/http"

	"github.com/pick-cee/go-ecommerce-api/internal/json"
)

type handler struct {
	service Service
}

func NewHandler (service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var tempUser createUserDto

	err := json.Read(r, &tempUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.service.CreateUser(r.Context(), tempUser)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.Write(w, http.StatusCreated, createdUser)
}

func (h *handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var tempUser loginUserDto

	err := json.Read(r, &tempUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loginUser, err := h.service.LoginUser(r.Context(), tempUser)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.Write(w, http.StatusCreated, loginUser)
}

