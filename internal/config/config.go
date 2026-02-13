package config

import (
	"calculate_product_packs/internal/domain"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	PackSizes []domain.PackSize
	Port      string
}

func NewConfig() *Config {
	return &Config{
		PackSizes: getPackSizesFromEnv(),
		Port:      getPortFromEnv(),
	}
}

func getPackSizesFromEnv() []domain.PackSize {
	packSizesStr := os.Getenv("PACK_SIZES")
	if packSizesStr == "" {
		return []domain.PackSize{250, 500, 1000, 2000, 5000}
	}

	sizesStr := strings.Split(packSizesStr, ",")
	var sizes []domain.PackSize
	for _, s := range sizesStr {
		size, err := strconv.Atoi(strings.TrimSpace(s))
		if err == nil && size > 0 {
			sizes = append(sizes, domain.PackSize(size))
		}
	}

	if len(sizes) == 0 {
		return []domain.PackSize{250, 500, 1000, 2000, 5000}
	}

	return sizes
}

func getPortFromEnv() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}
