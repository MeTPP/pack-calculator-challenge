package usecases

import (
	"calculate_product_packs/internal/domain"
	"math"
	"sort"
)

type CalculatePacksUseCase struct {
	repo domain.PackSizeRepository
}

func NewCalculatePacksUseCase(repo domain.PackSizeRepository) *CalculatePacksUseCase {
	return &CalculatePacksUseCase{repo: repo}
}

func (uc *CalculatePacksUseCase) Execute(orderSize int) ([]domain.PackResult, error) {
	if orderSize <= 0 {
		return nil, domain.ErrOrderSizePositive
	}

	packSizes := uc.repo.GetPackSizes()
	if len(packSizes) == 0 {
		return nil, domain.ErrNoPackSizes
	}

	sizes := make([]int, len(packSizes))
	for i, ps := range packSizes {
		sizes[i] = int(ps)
	}

	result := calculateOptimalPacks(orderSize, sizes)

	var packResults []domain.PackResult
	for size, count := range result {
		if count > 0 {
			packResults = append(packResults, domain.PackResult{Size: domain.PackSize(size), Count: count})
		}
	}

	sort.Slice(packResults, func(i, j int) bool {
		return packResults[i].Size > packResults[j].Size
	})

	return packResults, nil
}

// calculateOptimalPacks finds the optimal pack combination for the given order.
//
// Rules (in priority order):
//  1. Only whole packs can be sent
//  2. Minimize total items sent (must be >= orderSize)
//  3. Among solutions with equal total items, minimize number of packs
//
// Uses dynamic programming (variant of coin change problem). For large orders,
// pre-allocates largest packs to keep the DP table size manageable.
func calculateOptimalPacks(orderSize int, packSizes []int) map[int]int {
	sort.Ints(packSizes)

	minPack := packSizes[0]
	maxPack := packSizes[len(packSizes)-1]

	// For large orders, pre-subtract largest packs to keep the DP table small.
	dpLimit := minPack * maxPack
	if dpLimit < maxPack+minPack {
		dpLimit = maxPack + minPack
	}

	baseLargePacks := 0
	effOrder := orderSize
	if effOrder > dpLimit {
		baseLargePacks = (effOrder - dpLimit) / maxPack
		effOrder = orderSize - baseLargePacks*maxPack
	}

	maxTarget := effOrder + minPack - 1

	const inf = math.MaxInt32
	dp := make([]int, maxTarget+1)
	from := make([]int, maxTarget+1)
	for i := range dp {
		dp[i] = inf
	}
	dp[0] = 0

	for i := 1; i <= maxTarget; i++ {
		for _, pack := range packSizes {
			if pack > i {
				break
			}
			if dp[i-pack] < inf && dp[i-pack]+1 < dp[i] {
				dp[i] = dp[i-pack] + 1
				from[i] = pack
			}
		}
	}

	type solution struct {
		dpAmount   int
		largePacks int
		totalItems int
		totalPacks int
	}
	var best *solution

	maxExtra := effOrder / maxPack
	for extra := 0; extra <= maxExtra; extra++ {
		remainder := effOrder - extra*maxPack
		if remainder < 0 {
			break
		}
		searchEnd := remainder + minPack - 1
		if searchEnd > maxTarget {
			continue
		}

		for t := remainder; t <= searchEnd; t++ {
			if dp[t] < inf {
				totalItems := (baseLargePacks+extra)*maxPack + t
				totalPacks := baseLargePacks + extra + dp[t]

				if best == nil ||
					totalItems < best.totalItems ||
					(totalItems == best.totalItems && totalPacks < best.totalPacks) {
					best = &solution{
						dpAmount:   t,
						largePacks: baseLargePacks + extra,
						totalItems: totalItems,
						totalPacks: totalPacks,
					}
				}
				break
			}
		}
	}

	if best == nil {
		count := (orderSize + minPack - 1) / minPack
		return map[int]int{minPack: count}
	}

	result := make(map[int]int)
	if best.largePacks > 0 {
		result[maxPack] = best.largePacks
	}
	remaining := best.dpAmount
	for remaining > 0 {
		pack := from[remaining]
		result[pack]++
		remaining -= pack
	}

	return result
}
