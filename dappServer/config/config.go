package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/BurntSushi/toml"
)

// Struct to represent each node
type Node struct {
	Name string `toml:"name"`
	Port string `toml:"port"`
	DID  string `toml:"did"`
}

// Struct to hold the configuration
type Config struct {
	Nodes map[string]Node `toml:"nodes"`
}

var (
	instance *Config
	once     sync.Once
)

// LoadConfig initializes the configuration (Singleton)
func LoadConfig(filepath string) {
	once.Do(func() {
		instance = &Config{}
		if _, err := toml.DecodeFile(filepath, instance); err != nil {
			log.Fatalf("Error loading config file: %v", err)
		}
	})
}

// GetConfig returns the global configuration instance
func GetConfig() (*Config, error) {
	if instance == nil {
		// log.Fatal("Config not loaded. Call LoadConfig() first.")
		return nil, fmt.Errorf("Config not loaded. Call LoadConfig() first")
	}
	return instance, nil
}

// GetNodeNameByPort searches for a node by its port and returns its name
func GetNodeNameByPort(config *Config, port string) (string, bool) {
	for _, node := range config.Nodes {
		if node.Port == port {
			return node.Name, true
		}
	}
	return "", false
}

func GetNodeNameByDid(config *Config, did string) (string, bool) {
	for _, node := range config.Nodes {
		if node.DID == did {
			return node.Name, true
		}
	}
	return "", false
}
