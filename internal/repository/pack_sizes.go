package repository

import (
	"calculate_product_packs/internal/domain"
	"sync"
)

type MemoryPackSizeRepository struct {
	mu        sync.RWMutex
	packSizes []domain.PackSize
}

func NewMemoryPackSizeRepository(packSizes []domain.PackSize) domain.PackSizeRepository {
	cp := make([]domain.PackSize, len(packSizes))
	copy(cp, packSizes)
	return &MemoryPackSizeRepository{packSizes: cp}
}

func (r *MemoryPackSizeRepository) GetPackSizes() []domain.PackSize {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cp := make([]domain.PackSize, len(r.packSizes))
	copy(cp, r.packSizes)
	return cp
}

func (r *MemoryPackSizeRepository) UpdatePackSizes(sizes []domain.PackSize) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.packSizes = make([]domain.PackSize, len(sizes))
	copy(r.packSizes, sizes)
	return nil
}
