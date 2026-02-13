package usecases

import (
	"calculate_product_packs/internal/domain"
	"sort"
)

const (
	maxPackSize  = 1_000_000
	maxPackCount = 20
)

type PackSizesUseCase struct {
	repo domain.PackSizeRepository
}

func NewPackSizesUseCase(repo domain.PackSizeRepository) *PackSizesUseCase {
	return &PackSizesUseCase{repo: repo}
}

func (uc *PackSizesUseCase) UpdatePackSizes(sizes []domain.PackSize) error {
	if len(sizes) == 0 {
		return domain.ErrEmptyPackSizes
	}

	for _, size := range sizes {
		if size <= 0 || int(size) > maxPackSize {
			return domain.ErrInvalidPackSize
		}
	}

	// Deduplicate and sort.
	seen := make(map[domain.PackSize]bool, len(sizes))
	unique := make([]domain.PackSize, 0, len(sizes))
	for _, s := range sizes {
		if !seen[s] {
			seen[s] = true
			unique = append(unique, s)
		}
	}

	sort.Slice(unique, func(i, j int) bool { return unique[i] < unique[j] })

	if len(unique) > maxPackCount {
		return domain.ErrTooManyPackSizes
	}

	return uc.repo.UpdatePackSizes(unique)
}

func (uc *PackSizesUseCase) GetPackSizes() []domain.PackSize {
	return uc.repo.GetPackSizes()
}
