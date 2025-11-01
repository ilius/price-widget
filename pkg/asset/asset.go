package asset

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

func IdMap(assets []*Asset) map[string]struct{} {
	m := map[string]struct{}{}
	for _, a := range assets {
		m[a.ID] = struct{}{}
	}
	return m
}

// asset.Type:
// For gold and silver: "metal", "gold", "goldprice"
// For cryptocurrency: "", "coin", "crypto", "coingecko"

type Asset struct {
	Name   string `toml:"name"`
	ID     string `toml:"id"`
	Type   string `toml:"type"`
	Prefix string `toml:"prefix"`
	Suffix string `toml:"suffix"`

	Digits int `toml:"digits"`

	HumanizeFormat string `toml:"humanize_format"` // go-humanize format, used if digits < 0
}

func (a *Asset) numToStr(price float64) string {
	if a.Digits < 0 {
		return humanize.FormatFloat(a.HumanizeFormat, price)
	}
	return fmt.Sprintf("%.*f", a.Digits, price)
}

func (a *Asset) FormatPrice(price float64) string {
	return a.Prefix + a.numToStr(price) + a.Suffix
}
