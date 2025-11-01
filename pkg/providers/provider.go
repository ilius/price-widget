package providers

import "github.com/ilius/price-widget/pkg/asset"

type Provider interface {
	Name() string
	SupportedIDs() map[string]struct{}
	FetchPrices([]*asset.Asset) (map[string]float64, error)
}
