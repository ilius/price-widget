package main

import "fmt"

type Asset struct {
	Name   string `toml:"name"`
	ID     string `toml:"id"`
	Digits int    `toml:"digits"`
}

func (a *Asset) FormatPrice(price float64) string {
	return fmt.Sprintf("%.*f", a.Digits, price)
}
