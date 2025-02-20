package main

import (
	"dapp-server/config"
	"dapp-server/server"
)

const CONFIG_PATH = ".config/config.toml"

func main() {
	config.LoadConfig(CONFIG_PATH)
	server.BootupServer()
}
