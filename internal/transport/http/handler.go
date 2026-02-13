package http

import (
	"calculate_product_packs/internal/domain"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

//go:generate mockgen -destination=mocks/mock_pack_calculator.go -package=mocks calculate_product_packs/internal/transport/http PackCalculator
type PackCalculator interface {
	Execute(orderSize int) ([]domain.PackResult, error)
}

//go:generate mockgen -destination=mocks/mock_pack_sizer.go -package=mocks calculate_product_packs/internal/transport/http PackSizer
type PackSizer interface {
	UpdatePackSizes(sizes []domain.PackSize) error
	GetPackSizes() []domain.PackSize
}

type PackCalculatorHandler struct {
	packCalculator   PackCalculator
	packSizesUseCase PackSizer
}

func NewPackCalculatorHandler(
	packCalculator PackCalculator,
	packSizesUseCase PackSizer,
) *PackCalculatorHandler {
	return &PackCalculatorHandler{
		packCalculator:   packCalculator,
		packSizesUseCase: packSizesUseCase,
	}
}

func (h *PackCalculatorHandler) CalculatePacks(w http.ResponseWriter, r *http.Request) {
	orderSize, err := strconv.Atoi(r.URL.Query().Get("orderSize"))
	if err != nil {
		http.Error(w, "Invalid order size", http.StatusBadRequest)
		return
	}

	result, err := h.packCalculator.Execute(orderSize)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrOrderSizePositive):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, domain.ErrNoPackSizes):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, result)
}

func (h *PackCalculatorHandler) UpdatePackSizes(w http.ResponseWriter, r *http.Request) {
	var sizes []domain.PackSize
	if err := json.NewDecoder(r.Body).Decode(&sizes); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.packSizesUseCase.UpdatePackSizes(sizes); err != nil {
		switch {
		case errors.Is(err, domain.ErrEmptyPackSizes),
			errors.Is(err, domain.ErrInvalidPackSize),
			errors.Is(err, domain.ErrTooManyPackSizes):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Failed to update pack sizes", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Pack sizes updated successfully")); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}

func (h *PackCalculatorHandler) GetPackSizes(w http.ResponseWriter, r *http.Request) {
	sizes := h.packSizesUseCase.GetPackSizes()
	writeJSON(w, sizes)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
