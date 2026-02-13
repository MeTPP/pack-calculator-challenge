package repository

import (
	"calculate_product_packs/internal/domain"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemoryPackSizeRepository(t *testing.T) {
	packSizes := []domain.PackSize{250, 500, 1000}
	repo := NewMemoryPackSizeRepository(packSizes)
	assert.NotNil(t, repo)
}

func TestMemoryPackSizeRepository_GetPackSizes(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		packSizes := []domain.PackSize{250, 500, 1000}
		repo := NewMemoryPackSizeRepository(packSizes)
		assert.Equal(t, packSizes, repo.GetPackSizes())
	})

	t.Run("empty", func(t *testing.T) {
		repo := NewMemoryPackSizeRepository([]domain.PackSize{})
		assert.NotNil(t, repo.GetPackSizes())
		assert.Empty(t, repo.GetPackSizes())
	})
}

func TestMemoryPackSizeRepository_UpdatePackSizes(t *testing.T) {
	repo := NewMemoryPackSizeRepository([]domain.PackSize{250})

	newSizes := []domain.PackSize{100, 200, 300}
	err := repo.UpdatePackSizes(newSizes)
	require.NoError(t, err)
	assert.Equal(t, newSizes, repo.GetPackSizes())
}

func TestMemoryPackSizeRepository_GetReturnsACopy(t *testing.T) {
	repo := NewMemoryPackSizeRepository([]domain.PackSize{250, 500})

	sizes := repo.GetPackSizes()
	sizes[0] = 9999

	assert.Equal(t, []domain.PackSize{250, 500}, repo.GetPackSizes())
}

func TestMemoryPackSizeRepository_UpdateStoresACopy(t *testing.T) {
	repo := NewMemoryPackSizeRepository(nil)

	original := []domain.PackSize{100, 200}
	_ = repo.UpdatePackSizes(original)
	original[0] = 9999

	assert.Equal(t, []domain.PackSize{100, 200}, repo.GetPackSizes())
}

func TestMemoryPackSizeRepository_ConcurrentAccess(t *testing.T) {
	repo := NewMemoryPackSizeRepository([]domain.PackSize{250, 500, 1000})

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = repo.GetPackSizes()
		}()
		go func(v int) {
			defer wg.Done()
			_ = repo.UpdatePackSizes([]domain.PackSize{domain.PackSize(v)})
		}(i)
	}
	wg.Wait()
}
