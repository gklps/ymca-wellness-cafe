package server

import (
	"dapp-server/config"
	rubix_interaction "dapp-server/rubix-interaction"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	wasmbridge "github.com/rubixchain/rubix-wasm/go-wasm-bridge"
)

type ContractInputRequest struct {
	Port              string `json:"port"`
	SmartContractHash string `json:"smart_contract_hash"` //port should also be added here, so that the api can understand which node.
}

type RubixResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type SmartContractDataReply struct {
	RubixResponse
	SCTDataReply []SCTDataReply
}

type SCTDataReply struct {
	BlockNo           uint64
	BlockId           string
	SmartContractData string
}

func BootupServer() {
	gin.SetMode(gin.ReleaseMode) //
	log.Println("Current Gin Mode:", gin.Mode())

	// Initialize a Gin router
	router := gin.Default()
	log.Println("Current Gin Mode:", gin.Mode())

	// config := GetConfig()

	log.SetFlags(log.LstdFlags)

	// Configure CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	// nftDappCallbackHandler := config.ContractsInfo["nft"].CallBackUrl
	// ftDappCallbackHandler := config.ContractsInfo["ft"].CallBackUrl

	// Define endpoints
	// router.POST(nftDappCallbackHandler, nftDappHandler) // NFT
	router.POST("/api/call-back-trigger", ftDappHandler) // FT
	router.POST("/api/deploy-contract", APIDeployContract)
	router.POST("/api/execute-contract", APIExecuteContract)

	// router.GET("/request-status", getRequestStatusHandler)

	// Start the server on port 8080
	router.Run(":8080")
}

func contractInputHandler(w http.ResponseWriter, r *http.Request) {

	var req ContractInputRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	err3 := godotenv.Load()
	if err3 != nil {
		fmt.Println("Error loading .env file:", err3)
		return
	}
	port := req.Port
	nodeName := os.Getenv(port)
	fmt.Println(nodeName)

	resp := RubixResponse{Status: true, Message: "Callback Successful", Result: "Success"}
	json.NewEncoder(w).Encode(resp)

}

// Handler function for /callback/nft
func ftDappHandler(c *gin.Context) {
	var req ContractInputRequest
	fmt.Println("Handler trggered")
	cfg, err := config.LoadConfig(".config/config.toml")
	if err != nil {
		fmt.Errorf("failed to load config: %w", err)
	}

	err = json.NewDecoder(c.Request.Body).Decode(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// // config := GetConfig()
	smartContractHash := req.SmartContractHash
	fmt.Println("Received Smart Contract hash: ", smartContractHash)

	smartContractTokenData := rubix_interaction.GetSmartContractData(smartContractHash, cfg.Network.DeployerNodeURL) //config.NodeAddress)
	if smartContractTokenData == nil {
		fmt.Println("Unable to fetch latest smart contract data")
		return
	}

	fmt.Println("Smart Contract Token Data :", string(smartContractTokenData))

	var dataReply SmartContractDataReply

	if err := json.Unmarshal(smartContractTokenData, &dataReply); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Data reply in runDappHandler", dataReply)
	smartContractData := dataReply.SCTDataReply
	var relevantData string
	for _, reply := range smartContractData {
		fmt.Println("SmartContractData:", reply.SmartContractData)
		relevantData = reply.SmartContractData
	}
	var inputMap map[string]interface{}
	err1 := json.Unmarshal([]byte(relevantData), &inputMap)
	if err1 != nil {
		return
	}
	if len(inputMap) != 1 {
		return
	}

	var funcName string
	var inputStruct interface{}
	for key, value := range inputMap {
		funcName = key
		inputStruct = value
	}
	fmt.Println("The function name extracted =", funcName)
	fmt.Println("The inputStruct Value :", inputStruct)

	hostFnRegistry := wasmbridge.NewHostFunctionRegistry()
	wasmPath, err := getWasmContractPath(smartContractHash)
	if err != nil {
		fmt.Println("Failed to get wasm path")
	}
	// Initialize the WASM module

	wasmModule, err := wasmbridge.NewWasmModule(
		// config.ContractsInfo["ft"].ContractPath,
		wasmPath,
		hostFnRegistry,
		wasmbridge.WithRubixNodeAddress(cfg.Network.DeployerNodeURL), //config.NodeAddress),
		wasmbridge.WithQuorumType(2),
	)
	if err != nil {
		log.Printf("Failed to initialize WASM module: %v", err)
		return
	}

	executionResult, errExecuteContract := executeAndGetContractResult(wasmModule, relevantData)
	fmt.Println("----------- FT Execution Result: ", executionResult)
	if errExecuteContract != nil {
		fmt.Println("Rhe executionResult is ", executionResult)
		return
	}

	var response RubixResponse

	// Convert JSON string to struct
	if executionResult == "success" {
		response = RubixResponse{Status: true, Message: "NFT Transferred Succesfully"}
	} else {
		err = json.Unmarshal([]byte(executionResult), &response)
		if err != nil {
			log.Printf("Error parsing JSON: %v", err)
			return
		}
	}

	resultFinal := gin.H{
		"message": "DApp executed successfully",
		"data":    response,
	}

	// Return a response
	c.JSON(http.StatusOK, resultFinal)
}

func getWasmContractPath(contractHash string) (string, error) {
	currentWorkingDir, err := os.Getwd()
	fmt.Println("The current working Directory is : ", currentWorkingDir)
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}
	// Here this path should be dynamic
	contractDir := filepath.Join(currentWorkingDir, "rubix-nodes/node2/SmartContract", contractHash)

	entries, err := os.ReadDir(contractDir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".wasm") {
			return filepath.Join(contractDir, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("no wasm contract found in directory: %v", contractDir)
}

func executeAndGetContractResult(wasmModule *wasmbridge.WasmModule, contractInput string) (string, error) {
	// Call the function
	contractResult, err := wasmModule.CallFunction(contractInput)
	if err != nil {
		return "", fmt.Errorf("function call failed: %v", err)
	}

	return contractResult, nil
}
