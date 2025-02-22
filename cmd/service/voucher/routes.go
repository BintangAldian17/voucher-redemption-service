package voucher

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BintangAldian17/voucher-redemption-service/types"
	"github.com/BintangAldian17/voucher-redemption-service/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.VoucherStore
}

func NewHandler(store types.VoucherStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/voucher", h.hanldeCreateVoucher).Methods("POST")
	router.HandleFunc("/voucher/brand", h.handleGetVouchersByBrandID).Methods("GET").Queries("id", "{brand_id}")
	router.HandleFunc("/voucher", h.handleGetVoucherByID).Methods("GET").Queries("id", "{voucher_id}")
}

func (h *Handler) handleGetVoucherByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["voucher_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing voucher ID"))
		return
	}

	voucherID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid voucher ID"))
		return
	}

	voucher, err := h.store.GetVoucherByID(voucherID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJSON(w, http.StatusOK, voucher)
}

func (h *Handler) hanldeCreateVoucher(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateVoucherPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusOK, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	_, err := h.store.GetVoucherByName(payload.Name)

	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("voucher with name %s already exists", payload.Name))
		return
	}
	err = h.store.CreateVoucher(types.Voucher{
		Name:         payload.Name,
		Description:  payload.Description,
		BrandID:      payload.BrandID,
		CostInPoints: payload.CostInPoints,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, err)
}

func (h *Handler) handleGetVouchersByBrandID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strBrandID, ok := vars["brand_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing brand ID"))
		return
	}

	brandID, err := strconv.Atoi(strBrandID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid brand ID"))
		return
	}

	vouchers, err := h.store.GetVouchersByBrandID(brandID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, vouchers)
}
