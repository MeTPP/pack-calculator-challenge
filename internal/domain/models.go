package domain

type PackSize int

type PackResult struct {
	Size  PackSize `json:"size"`
	Count int      `json:"count"`
}

//go:generate mockgen -destination=mocks/mock_pack_size_repository.go -package=mocks calculate_product_packs/internal/domain PackSizeRepository
type PackSizeRepository interface {
	GetPackSizes() []PackSize
	UpdatePackSizes(sizes []PackSize) error
}
