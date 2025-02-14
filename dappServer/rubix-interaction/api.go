package rubix_interaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetSmartContractData(token string, address string) []byte {
	data := map[string]interface{}{
		"token":  token,
		"latest": true,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil
	}
	url := address + "/api/get-smart-contract-token-chain-data"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return nil
	}

	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return nil
	}
	// Process the data as needed
	fmt.Println("Response Body in get smart contract data :", string(data2))

	return data2

}

func RegisterCallBackUrl(smartContractTokenHash string, urlPort string, endPoint string, nodePort string) {
	callBackUrl := fmt.Sprintf("http://localhost:%s/%s", urlPort, endPoint)
	data := map[string]interface{}{
		"CallBackURL":        callBackUrl,
		"SmartContractToken": smartContractTokenHash,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/register-callback-url", nodePort)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	fmt.Println("Response Body in register callback url :", string(data2))
}

// // Execute handles the contract execution process
// func Execute(
// 	contractHash string, executorDid string,
// 	homeDir string, contractDir string, contractMsgFile string,
// ) (*ExecutionResult, error) {
// 	// Load config to get API URL
// 	cfg, err := config.LoadConfig(homeDir)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to load config: %w", err)
// 	}

// 	contractMsg, err := parseContractMsgFromJSON(contractMsgFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read contract message file: %w", err)
// 	}

// 	// Call execute-smart-contract API
// 	requestID, err := executeSmartContract(cfg.Network.DeployerNodeURL, contractHash, executorDid, contractMsg)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute smart contract: %w", err)
// 	}

// 	// Call signature-response API
// 	if err := signatureResponse(cfg.Network.DeployerNodeURL, requestID); err != nil {
// 		return nil, fmt.Errorf("failed to process signature response: %w", err)
// 	}

// 	contractResult, err := callWasm(contractDir, contractMsg)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to call wasm contract: %w", err)
// 	}

// 	return &ExecutionResult{
// 		ContractResult: contractResult,
// 		Success: true,
// 		Message: "Contract executed successfully",
// 	}, nil
// }

// func executeSmartContract(baseURL, contractHash, executorDid, contractMsg string) (string, error) {
// 	// Create request body
// 	requestBody := struct {
// 		Comment            string `json:"comment"`
// 		ExecutorAddr       string `json:"executorAddr"`
// 		QuorumType         int    `json:"quorumType"`
// 		SmartContractData  string `json:"smartContractData"`
// 		SmartContractToken string `json:"smartContractToken"`
// 	}{
// 		Comment:            "Contract execution",
// 		ExecutorAddr:       executorDid,
// 		QuorumType:         2,
// 		SmartContractData:  contractMsg,
// 		SmartContractToken: contractHash,
// 	}

// 	// Marshal request body
// 	bodyBytes, err := json.Marshal(requestBody)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to marshal request body: %w", err)
// 	}

// 	// Create request
// 	requestURL, err := url.JoinPath(baseURL, "/api/execute-smart-contract")
// 	if err != nil {
// 		return "", fmt.Errorf("execute: unable to form request URL")
// 	}

// 	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(bodyBytes))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create request: %w", err)
// 	}

// 	// Set headers
// 	req.Header.Set("Content-Type", "application/json")

// 	// Send request
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to send request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	// Read response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read response: %w", err)
// 	}

// 	// Parse response
// 	var apiResp SmartContractAPIResponseV2
// 	if err := json.Unmarshal(body, &apiResp); err != nil {
// 		return "", fmt.Errorf("failed to parse response: %w", err)
// 	}

// 	// Check response status
// 	if !apiResp.Status {
// 		return "", fmt.Errorf(apiResp.Message)
// 	}

// 	return apiResp.Result.Id, nil
// }
