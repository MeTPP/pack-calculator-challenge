package usecases

import (
	"calculate_product_packs/internal/domain"
	"calculate_product_packs/internal/domain/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPackSizesUseCase_UpdatePackSizes(t *testing.T) {
	tests := []struct {
		name    string
		sizes   []domain.PackSize
		wantErr error
		stored  []domain.PackSize
	}{
		{
			name:   "valid sizes are sorted",
			sizes:  []domain.PackSize{500, 250, 1000},
			stored: []domain.PackSize{250, 500, 1000},
		},
		{
			name:    "empty slice",
			sizes:   []domain.PackSize{},
			wantErr: domain.ErrEmptyPackSizes,
		},
		{
			name:    "negative size",
			sizes:   []domain.PackSize{250, -1},
			wantErr: domain.ErrInvalidPackSize,
		},
		{
			name:    "zero size",
			sizes:   []domain.PackSize{0, 250},
			wantErr: domain.ErrInvalidPackSize,
		},
		{
			name:   "deduplicates",
			sizes:  []domain.PackSize{250, 500, 250},
			stored: []domain.PackSize{250, 500},
		},
		{
			name:    "exceeds max size",
			sizes:   []domain.PackSize{1_000_001},
			wantErr: domain.ErrInvalidPackSize,
		},
		{
			name:   "max size boundary",
			sizes:  []domain.PackSize{1_000_000},
			stored: []domain.PackSize{1_000_000},
		},
		{
			name:    "exceeds max pack count",
			sizes:   []domain.PackSize{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21},
			wantErr: domain.ErrTooManyPackSizes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockPackSizeRepository(ctrl)
			if tt.stored != nil {
				mockRepo.EXPECT().UpdatePackSizes(tt.stored).Return(nil)
			}

			uc := NewPackSizesUseCase(mockRepo)
			err := uc.UpdatePackSizes(tt.sizes)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPackSizesUseCase_GetPackSizes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []domain.PackSize{250, 500, 1000}
	mockRepo := mocks.NewMockPackSizeRepository(ctrl)
	mockRepo.EXPECT().GetPackSizes().Return(expected)

	uc := NewPackSizesUseCase(mockRepo)
	assert.Equal(t, expected, uc.GetPackSizes())
}
