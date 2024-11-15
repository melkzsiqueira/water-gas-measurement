package main

import "github.com/melkzsiqueira/water-gas-measurement/configs"

func main() {
	config, _ := configs.LoadConfig(".")
	println(config.DBDriver)
}
