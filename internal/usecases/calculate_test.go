package usecases

import (
	"calculate_product_packs/internal/domain"
	"calculate_product_packs/internal/domain/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCalculatePacksUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		packSizes     []domain.PackSize
		orderSize     int
		expectedPacks []domain.PackResult
		expectedError error
	}{
		{
			name:          "order smaller than smallest pack rounds up to smallest pack",
			packSizes:     []domain.PackSize{250, 500, 1000, 2000, 5000},
			orderSize:     1,
			expectedPacks: []domain.PackResult{{Size: 250, Count: 1}},
		},
		{
			name:          "order exactly matches smallest pack",
			packSizes:     []domain.PackSize{250, 500, 1000, 2000, 5000},
			orderSize:     250,
			expectedPacks: []domain.PackResult{{Size: 250, Count: 1}},
		},
		{
			name:          "order slightly over pack size rounds up to next pack",
			packSizes:     []domain.PackSize{250, 500, 1000, 2000, 5000},
			orderSize:     251,
			expectedPacks: []domain.PackResult{{Size: 500, Count: 1}},
		},
		{
			name:          "order requires combination of two packs",
			packSizes:     []domain.PackSize{250, 500, 1000, 2000, 5000},
			orderSize:     501,
			expectedPacks: []domain.PackResult{{Size: 500, Count: 1}, {Size: 250, Count: 1}},
		},
		{
			name:          "large order uses multiple pack sizes",
			packSizes:     []domain.PackSize{250, 500, 1000, 2000, 5000},
			orderSize:     12001,
			expectedPacks: []domain.PackResult{{Size: 5000, Count: 2}, {Size: 2000, Count: 1}, {Size: 250, Count: 1}},
		},
		{
			name:          "order just under round number rounds up to fewer packs",
			packSizes:     []domain.PackSize{250, 500, 1000, 2000, 5000},
			orderSize:     9999,
			expectedPacks: []domain.PackResult{{Size: 5000, Count: 2}},
		},
		{
			name:          "order exactly matches a middle pack size",
			packSizes:     []domain.PackSize{250, 500, 1000, 2000, 5000},
			orderSize:     500,
			expectedPacks: []domain.PackResult{{Size: 500, Count: 1}},
		},
		{
			name:          "very large order uses only largest pack",
			packSizes:     []domain.PackSize{250, 500, 1000, 2000, 5000},
			orderSize:     1000000,
			expectedPacks: []domain.PackResult{{Size: 5000, Count: 200}},
		},
		{
			name:          "single available pack size rounds up to fill order",
			packSizes:     []domain.PackSize{1000},
			orderSize:     2500,
			expectedPacks: []domain.PackResult{{Size: 1000, Count: 3}},
		},
		{
			name:          "returns error when no pack sizes available",
			packSizes:     []domain.PackSize{},
			orderSize:     1000,
			expectedPacks: nil,
			expectedError: domain.ErrNoPackSizes,
		},
		{
			name:      "non-standard pack sizes with large order finds optimal combination",
			packSizes: []domain.PackSize{17, 31, 47},
			orderSize: 5000,
			expectedPacks: []domain.PackResult{
				{Size: 47, Count: 105},
				{Size: 31, Count: 1},
				{Size: 17, Count: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockPackSizeRepository(ctrl)
			mockRepo.EXPECT().GetPackSizes().Return(tt.packSizes)

			useCase := NewCalculatePacksUseCase(mockRepo)
			result, err := useCase.Execute(tt.orderSize)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedPacks, result)
		})
	}
}

func TestCalculatePacksUseCase_Execute_SortsPackSizes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPackSizeRepository(ctrl)
	mockRepo.EXPECT().GetPackSizes().Return([]domain.PackSize{250, 1000, 500, 5000, 2000})

	useCase := NewCalculatePacksUseCase(mockRepo)
	result, err := useCase.Execute(12001)

	assert.NoError(t, err)
	expected := []domain.PackResult{{Size: 5000, Count: 2}, {Size: 2000, Count: 1}, {Size: 250, Count: 1}}
	assert.Equal(t, expected, result)
}

func TestCalculatePacksUseCase_Execute_EmptyPackSizes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPackSizeRepository(ctrl)
	mockRepo.EXPECT().GetPackSizes().Return([]domain.PackSize{})

	useCase := NewCalculatePacksUseCase(mockRepo)
	result, err := useCase.Execute(1000)

	assert.ErrorIs(t, err, domain.ErrNoPackSizes)
	assert.Empty(t, result)
}

func TestCalculatePacksUseCase_Execute_InvalidInputs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPackSizeRepository(ctrl)
	useCase := NewCalculatePacksUseCase(mockRepo)

	result, err := useCase.Execute(-100)
	assert.ErrorIs(t, err, domain.ErrOrderSizePositive)
	assert.Empty(t, result)

	result, err = useCase.Execute(0)
	assert.ErrorIs(t, err, domain.ErrOrderSizePositive)
	assert.Empty(t, result)
}

func BenchmarkCalculateOptimalPacks(b *testing.B) {
	sizes := []int{250, 500, 1000, 2000, 5000}
	for i := 0; i < b.N; i++ {
		calculateOptimalPacks(12001, sizes)
	}
}

func BenchmarkCalculateOptimalPacks_EdgeCase(b *testing.B) {
	sizes := []int{17, 31, 47}
	for i := 0; i < b.N; i++ {
		calculateOptimalPacks(5000, sizes)
	}
}
