package brand

import (
	"fmt"
	"net/http"

	"github.com/BintangAldian17/voucher-redemption-service/types"
	"github.com/BintangAldian17/voucher-redemption-service/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.BrandStore
}

func NewHandler(store types.BrandStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/brand", h.handleCreateBrand).Methods("POST")
}

func (h *Handler) handleCreateBrand(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateBrandPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	_, err := h.store.GetBrandByName(payload.Name)

	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("brand with name %s already exists", payload.Name))
		return
	}

	err = h.store.CreateBrand(types.Brand{
		Name:        payload.Name,
		Description: payload.Description,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, err)
}
