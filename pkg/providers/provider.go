package providers

import "github.com/ilius/price-widget/pkg/asset"

type Provider interface {
	FetchPrices([]*asset.Asset) (map[string]float64, error)
}
