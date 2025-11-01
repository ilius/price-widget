package metals

import (
	"time"

	"github.com/ilius/price-widget/pkg/asset"
	"github.com/ilius/price-widget/pkg/goldpriceorg"
	"github.com/ilius/price-widget/pkg/managers"
)

func New(
	assets []*asset.Asset,
	refreshInterval time.Duration,
) *managers.Manager {
	return managers.NewManager(
		goldpriceorg.New(),
		assets,
		refreshInterval,
	)
}
