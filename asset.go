package main

type Asset struct {
	Name   string `toml:"name"`
	ID     string `toml:"id"`
	Digits int    `toml:"digits"`
}
