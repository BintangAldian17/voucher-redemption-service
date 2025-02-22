package transaction

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
	store         types.TransactionStore
	customerStore types.CustomerStore
	voucherStore  types.VoucherStore
}

func NewHandler(store types.TransactionStore, customerStore types.CustomerStore, voucherStore types.VoucherStore) *Handler {
	return &Handler{store: store, customerStore: customerStore, voucherStore: voucherStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/transaction/redemption", h.handleMakeRedemption).Methods("POST")
	router.HandleFunc("/transaction/redemption", h.handleGetTransactionByID).Methods("GET").Queries("transactionId", "{transactionId}")
}

func (h *Handler) handleMakeRedemption(w http.ResponseWriter, r *http.Request) {
	var payload types.RedemptionPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	customer, err := h.customerStore.GetCustomerByID(payload.CustomerID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid customer id"))
		return
	}

	totalPointsNeeded := 0
	transactionItems := make([]types.TransactionItem, len(payload.Vouchers))

	for i, item := range payload.Vouchers {
		voucher, err := h.voucherStore.GetVoucherByID(item.VoucherID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid voucher id"))
			return
		}
		totalPointsNeeded += voucher.CostInPoints * item.Quantity
		transactionItems[i] = types.TransactionItem{
			VoucherID:    item.VoucherID,
			Quantity:     item.Quantity,
			CostInPoints: voucher.CostInPoints * item.Quantity,
		}
	}

	if customer.TotalPoints < totalPointsNeeded {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("insufficient points"))
		return
	}

	tx, err := h.store.Begin()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	transactionID, err := h.store.RedeemVouchers(tx, customer, transactionItems, totalPointsNeeded)
	if err != nil {
		tx.Rollback()
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"transaction_id": transactionID,
		"message":        "Redemption successful",
	})

}

func (h *Handler) handleGetTransactionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["transactionId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing transaction ID"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid transaction ID"))
		return
	}

	transaction, err := h.store.GetTransactionWithItemsByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, transaction)
}
