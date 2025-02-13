package contract

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	wasmbridge "github.com/rubixchain/rubix-wasm/go-wasm-bridge"
)



func callWasm(contractDir string, contractMsg string) (string, error) {
	// This must be changed to the path inside the rubix node
	wasmModulePath, err := getWasmContractPath(contractDir)
	if err != nil {
		return "", fmt.Errorf("failed to get wasm contract path: %w", err)
	}

	hostFnRegistry := wasmbridge.NewHostFunctionRegistry()

	wasmModule, err := wasmbridge.NewWasmModule(wasmModulePath, hostFnRegistry)
	if err != nil {
		return "", fmt.Errorf("failed to create wasm module: %w", err)
	}

	contractResult, err := wasmModule.CallFunction(contractMsg)
	if err != nil {
		return "", fmt.Errorf("failed to call Contract function: %w", err)
	}

	return contractResult, nil
}

func getWasmContractPath(contractDir string) (string, error) {
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	artifactsDir := filepath.Join(currentWorkingDir, "artifacts")

	entries, err := os.ReadDir(artifactsDir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".wasm") {
			return filepath.Join(artifactsDir, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("no wasm contract found in directory: %v", contractDir)
}
