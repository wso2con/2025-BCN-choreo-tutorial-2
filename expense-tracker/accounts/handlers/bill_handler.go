package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jo/choreo-tutorial/accounts/db"
	"github.com/jo/choreo-tutorial/accounts/models"
)

// BillHandler handles bill-related requests
type BillHandler struct {
	db db.Database
}

// NewBillHandler creates a new bill handler
func NewBillHandler(database db.Database) *BillHandler {
	return &BillHandler{db: database}
}

// GetBills returns all bills
// @Summary Get all bills
// @Description Returns a list of all bills with summary information
// @Tags bills
// @Produce json
// @Success 200 {array} models.BillSummary
// @Failure 500 {object} map[string]string
// @Router /bills [get]
func (h *BillHandler) GetBills(w http.ResponseWriter, r *http.Request) {
	bills, err := h.db.GetBills()
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	responseJSON(w, bills)
}

// GetBill returns a single bill
// @Summary Get a single bill
// @Description Returns a single bill with all its items
// @Tags bills
// @Produce json
// @Param id path int true "Bill ID"
// @Success 200 {object} models.Bill
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bills/{id} [get]
func (h *BillHandler) GetBill(w http.ResponseWriter, r *http.Request) {
	id, err := getBillID(r)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	bill, err := h.db.GetBill(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, errors.New("bill not found"), http.StatusNotFound)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	responseJSON(w, bill)
}

// CreateBill creates a new bill
// @Summary Create a new bill
// @Description Creates a new bill with the provided information
// @Tags bills
// @Accept json
// @Produce json
// @Param bill body models.BillInput true "Bill information"
// @Success 201 {object} map[string]int64
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bills [post]
func (h *BillHandler) CreateBill(w http.ResponseWriter, r *http.Request) {
	var billInput models.BillInput
	err := json.NewDecoder(r.Body).Decode(&billInput)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	// Validate input
	if billInput.Title == "" {
		writeError(w, errors.New("title is required"), http.StatusBadRequest)
		return
	}

	// Create bill
	id, err := h.db.CreateBill(&billInput)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	responseJSON(w, map[string]int64{"id": id})
}

// UpdateBill updates an existing bill
// @Summary Update a bill
// @Description Updates an existing bill with the provided information
// @Tags bills
// @Accept json
// @Produce json
// @Param id path int true "Bill ID"
// @Param bill body models.BillInput true "Bill information"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bills/{id} [put]
func (h *BillHandler) UpdateBill(w http.ResponseWriter, r *http.Request) {
	id, err := getBillID(r)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	var billInput models.BillInput
	err = json.NewDecoder(r.Body).Decode(&billInput)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	// Validate input
	if billInput.Title == "" {
		writeError(w, errors.New("title is required"), http.StatusBadRequest)
		return
	}

	// Check if bill exists
	_, err = h.db.GetBill(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, errors.New("bill not found"), http.StatusNotFound)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	// Update bill
	err = h.db.UpdateBill(id, &billInput)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	responseJSON(w, map[string]string{"message": "Bill updated successfully"})
}

// DeleteBill deletes a bill
// @Summary Delete a bill
// @Description Deletes a bill and all its items
// @Tags bills
// @Produce json
// @Param id path int true "Bill ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bills/{id} [delete]
func (h *BillHandler) DeleteBill(w http.ResponseWriter, r *http.Request) {
	id, err := getBillID(r)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	// Check if bill exists
	_, err = h.db.GetBill(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, errors.New("bill not found"), http.StatusNotFound)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	// Delete bill
	err = h.db.DeleteBill(id)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	responseJSON(w, map[string]string{"message": "Bill deleted successfully"})
}

// Helper functions

// getBillID extracts the bill ID from the URL
func getBillID(r *http.Request) (int64, error) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return 0, errors.New("invalid bill ID")
	}
	return id, nil
}

// writeError writes an error response
func writeError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	responseJSON(w, map[string]string{"error": err.Error()})
}

// responseJSON writes a JSON response
func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}