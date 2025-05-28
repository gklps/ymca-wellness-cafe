package main

import (
	"dapp-server/config"
	"dapp-server/server"
)

const CONFIG_PATH = ".config/config.toml"

func main() {
	// Create a new registry
	// registry := wasmbridge.NewHostFunctionRegistry()

	// // Create your custom host function
	// registry.Register(rubix_interaction.NewWriteToJsonFile())
	// hostFunction := registry.GetHostFunctions()
	// fmt.Println("Host function is :", hostFunction)
	config.LoadConfig(CONFIG_PATH)
	config.LoadEnvConfig()
	server.BootupServer()

}
