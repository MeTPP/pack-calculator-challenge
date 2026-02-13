package http

import (
	"bytes"
	"calculate_product_packs/internal/domain"
	"calculate_product_packs/internal/transport/http/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPackCalculatorHandler_CalculatePacks(t *testing.T) {
	tests := []struct {
		name           string
		orderSize      string
		mockSetup      func(m *mocks.MockPackCalculator)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Valid order size",
			orderSize: "500",
			mockSetup: func(m *mocks.MockPackCalculator) {
				m.EXPECT().Execute(500).Return([]domain.PackResult{{Size: 500, Count: 1}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"size":500,"count":1}]` + "\n",
		},
		{
			name:           "Invalid order size",
			orderSize:      "invalid",
			mockSetup:      func(m *mocks.MockPackCalculator) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid order size\n",
		},
		{
			name:      "Order size must be greater than zero",
			orderSize: "0",
			mockSetup: func(m *mocks.MockPackCalculator) {
				m.EXPECT().Execute(0).Return(nil, domain.ErrOrderSizePositive)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "order size must be greater than 0\n",
		},
		{
			name:      "No pack sizes available",
			orderSize: "100",
			mockSetup: func(m *mocks.MockPackCalculator) {
				m.EXPECT().Execute(100).Return(nil, domain.ErrNoPackSizes)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "no pack sizes available\n",
		},
		{
			name:      "Internal server error",
			orderSize: "1000",
			mockSetup: func(m *mocks.MockPackCalculator) {
				m.EXPECT().Execute(1000).Return(nil, errors.New("unexpected error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "unexpected error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCalculator := mocks.NewMockPackCalculator(ctrl)
			tt.mockSetup(mockCalculator)

			handler := NewPackCalculatorHandler(mockCalculator, nil)

			req := httptest.NewRequest("GET", "/api/calculate?orderSize="+tt.orderSize, nil)
			rr := httptest.NewRecorder()
			handler.CalculatePacks(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestPackCalculatorHandler_CalculatePacks_JSONResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCalculator := mocks.NewMockPackCalculator(ctrl)
	expectedResult := []domain.PackResult{
		{Size: 500, Count: 1},
		{Size: 250, Count: 1},
	}
	mockCalculator.EXPECT().Execute(750).Return(expectedResult, nil)

	handler := NewPackCalculatorHandler(mockCalculator, nil)

	req := httptest.NewRequest("GET", "/api/calculate?orderSize=750", nil)
	rr := httptest.NewRecorder()
	handler.CalculatePacks(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var result []domain.PackResult
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestPackCalculatorHandler_UpdatePackSizes(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockSetup      func(m *mocks.MockPackSizer)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "valid update",
			body: `[250, 500, 1000]`,
			mockSetup: func(m *mocks.MockPackSizer) {
				m.EXPECT().UpdatePackSizes([]domain.PackSize{250, 500, 1000}).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Pack sizes updated successfully",
		},
		{
			name:           "invalid JSON",
			body:           `not json`,
			mockSetup:      func(m *mocks.MockPackSizer) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body\n",
		},
		{
			name: "empty array",
			body: `[]`,
			mockSetup: func(m *mocks.MockPackSizer) {
				m.EXPECT().UpdatePackSizes([]domain.PackSize{}).Return(domain.ErrEmptyPackSizes)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "pack sizes cannot be empty\n",
		},
		{
			name: "negative sizes",
			body: `[250, -1]`,
			mockSetup: func(m *mocks.MockPackSizer) {
				m.EXPECT().UpdatePackSizes([]domain.PackSize{250, -1}).Return(domain.ErrInvalidPackSize)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid pack size\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSizer := mocks.NewMockPackSizer(ctrl)
			tt.mockSetup(mockSizer)

			handler := NewPackCalculatorHandler(nil, mockSizer)

			req := httptest.NewRequest("PUT", "/api/pack-sizes", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler.UpdatePackSizes(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestPackCalculatorHandler_GetPackSizes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSizer := mocks.NewMockPackSizer(ctrl)
	mockSizer.EXPECT().GetPackSizes().Return([]domain.PackSize{250, 500, 1000})

	handler := NewPackCalculatorHandler(nil, mockSizer)

	req := httptest.NewRequest("GET", "/api/pack-sizes", nil)
	rr := httptest.NewRecorder()
	handler.GetPackSizes(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var sizes []domain.PackSize
	err := json.Unmarshal(rr.Body.Bytes(), &sizes)
	assert.NoError(t, err)
	assert.Equal(t, []domain.PackSize{250, 500, 1000}, sizes)
}
