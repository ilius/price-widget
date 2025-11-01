package cryptomanager

import (
	"time"

	"github.com/ilius/price-widget/pkg/asset"
	"github.com/ilius/price-widget/pkg/managers"
	"github.com/ilius/price-widget/pkg/providers/coingecko"
)

func New(
	assets []*asset.Asset,
	refreshInterval time.Duration,
) *managers.Manager {
	return managers.NewManager(
		coingecko.New(),
		assets,
		refreshInterval,
	)
}
