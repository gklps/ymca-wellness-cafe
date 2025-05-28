package server

import (
	"dapp-server/config"
	rubix_interaction "dapp-server/rubix-interaction"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	wasmbridge "github.com/rubixchain/rubix-wasm/go-wasm-bridge"
)

const SMART_CONTRACT_HASH = "QmZdkRPESpodVMMpYaf6bvPQ2bMckjMzQKaoBaY7C9jjdD"

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
	SCTDataReply []SCTDataReply `json:"SCTDataReply"`
}

type SCTDataReply struct {
	BlockNo            uint64 `json:"BlockNo"`
	BlockId            string `json:"BlockId"`
	SmartContractData  string `json:"SmartContractData"`
	Epoch              uint64 `json:"Epoch"`
	InitiatorSignature string `json:"InitiatorSignature"`
	ExecutorDID        string `json:"ExecutorDID"`
	InitiatorSignData  string `json:"InitiatorSignData"`
}

type AddActivityRequest struct {
	ActivityID   string `json:"activity_id"`
	RewardPoints int    `json:"reward_points"`
	AdminDID     string `json:"admin_did"`
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
	router.POST("/api/trigger-contract-2", ftContract2Handler)
	router.POST("/api/deploy-contract", APIDeployContract)
	router.POST("/api/execute-contract", APIExecuteContract)
	router.POST("/api/activity/add", APIAddActivity)
	router.POST("/api/callback/trigger", APICallBackTrigger)

	// router.GET("/request-status", getRequestStatusHandler)

	// Start the server on port 8080
	router.Run(":8080")
}

func APIAddActivity(c *gin.Context) {
	fmt.Println("APIAddActivity triggered")
	var req AddActivityRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	fmt.Println("The request body is:", req)
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println("failed to load config: %w", err)
	}
	nodePort, exists := config.GetPortByDid(cfg, req.AdminDID)
	if !exists {
		fmt.Println("failed to get node port: not found")
		return
	}
	fmt.Println("The node port is:", nodePort)
	url := fmt.Sprintf("http://localhost:%s", nodePort)
	fmt.Println("The url is :", url)
	contractMsg := fmt.Sprintf(`{"activity_id":"%s","reward_points":%d}`, req.ActivityID, req.RewardPoints)
	fmt.Println("The contract message is:", contractMsg)
	smartContractResponse, err := rubix_interaction.ExecuteSmartContract(url, SMART_CONTRACT_HASH, req.AdminDID, contractMsg)
	if err != nil {
		fmt.Println("failed to execute smart contract:", err)
		return
	}
	fmt.Println("Smart contract response:", smartContractResponse)
	rubix_interaction.SignatureResponse(url, smartContractResponse)
	if err != nil {
		fmt.Println("failed to send signature response:", err)
		return
	}
	fmt.Println("Signature response sent successfully")
	resultFinal := gin.H{
		"message": "DApp executed successfully",
		"data":    smartContractResponse,
	}

	// Return a response
	c.JSON(http.StatusOK, resultFinal)

}

func APICallBackTrigger(c *gin.Context) {
	var req ContractInputRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	fmt.Println("The request body is:", req)
	url := fmt.Sprintf("http://localhost:%s", req.Port)
	fmt.Println("The url is :", url)

	// // config := GetConfig()
	smartContractHash := req.SmartContractHash
	fmt.Println("Received Smart Contract hash: ", smartContractHash)

	smartContractTokenData := rubix_interaction.GetSmartContractData(smartContractHash, url) //config.NodeAddress)
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
	fmt.Println("Data reply in APICallBackTrigger", dataReply)
	smartContractData := dataReply.SCTDataReply
	var relevantBlock *SCTDataReply
	// var data []map[string]interface{}

	// err = json.Unmarshal([]byte(smartContractData), &data)
	// if err != nil {
	// 	panic(err)
	// }

	// // Access the value
	// firstKey := ""
	// for key := range data[0] { // Extract the first key (e.g., "2")
	// 	firstKey = key
	// 	break
	// }
	// fmt.Println(firstKey)

	// var smartContractData []SCTDataReply

	// Unmarshal JSON into the variable
	// err = json.Unmarshal([]byte(smartContractData), &smartContractData)
	// if err != nil {
	// 	log.Fatalf("Error unmarshaling JSON: %v", err)
	// }

	// Print the result
	var blockId string
	var blockNo uint64
	for _, data := range smartContractData {
		// fmt.Printf("BlockNo: %d, BlockId: %s, SmartContractData: %s\n",
		// 	data.BlockNo, data.BlockId, data.SmartContractData)
		// fmt.Println("The smart contract data is :", data.SmartContractData)
		relevantBlock = &data // Assuming you want the last block
		fmt.Println("The relevant block is :", relevantBlock)
		blockId = data.BlockId
		blockNo = data.BlockNo
	}
	if blockNo == 0 {
		return
		fmt.Println("The block number is 13")
		relevantBlock = getNextSCTDataAfterBlockID(smartContractData, blockId)
		fmt.Println("The relevant block is :", relevantBlock)
	} else {
		fmt.Println("The block number is not 0")
		// blockId, err := getBlockIDFromJSONFile("C:/Users/allen/Working-repo/ymca/ymca-wellness-cafe-project/dappServer/test.json")
		// if err != nil {
		// 	fmt.Println("Error reading block ID from JSON file:", err)
		// 	// return
		// }
		// fmt.Println("The block ID is:", blockId)
		//Here we have an array of SCTDataReply, in this we need to extract the latest one
		// relevantBlock = getNextSCTDataAfterBlockID(smartContractData, blockId)
		// fmt.Println("The relevant block is :", relevantBlock)
		//
	}
	var parsedData struct {
		ActivityID   string `json:"activity_id"`
		RewardPoints int    `json:"reward_points"`
	}

	err = json.Unmarshal([]byte(relevantBlock.SmartContractData), &parsedData)
	if err != nil {
		fmt.Println("Error:", err)
		// return
	}
	registry := wasmbridge.NewHostFunctionRegistry()

	// Create your custom host function
	registry.Register(rubix_interaction.NewWriteToJsonFile())
	hostFunction := registry.GetHostFunctions()
	fmt.Println("Host function is :", hostFunction)
	wasmPath, err := getWasmContractPath(smartContractHash, req.Port)
	if err != nil {
		fmt.Println("Failed to get wasm path")
	}
	wasmModule, err := wasmbridge.NewWasmModule(
		wasmPath,
		registry,
		// wasmbridge.WithRubixNodeAddress("http://localhost:20002"), //config.NodeAddress),
		// wasmbridge.WithQuorumType(2),
	)
	if err != nil {
		log.Printf("Failed to initialize WASM module: %v", err)
		return
	}
	contractInput := fmt.Sprintf(`{"add_activity": {"activity_id":"%s","reward_points":%d,"block_hash":"%s"}}`, parsedData.ActivityID, parsedData.RewardPoints, relevantBlock.BlockId)
	fmt.Println("The contract input is :", contractInput)
	result, err := wasmModule.CallFunction(contractInput)
	if err != nil {
		log.Printf("Failed to call WASM function: %v", err)
		return
	}
	fmt.Println("The result is :", result)
}

// Function to read BlockId from a JSON file
func getBlockIDFromJSONFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var data struct {
		BlockID      string `json:"block_id"`
		ActivityID   string `json:"activity_id"`
		RewardPoints int    `json:"reward_points"`
	}

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return "", err
	}

	return data.BlockID, nil
}

// func getSCTDataAfterBlockID(sctDataReplies []SCTDataReply, blockID string) []SCTDataReply {
// 	var result []SCTDataReply
// 	found := false

// 	for _, data := range sctDataReplies {
// 		if found {
// 			result = append(result, data)
// 		} else if data.BlockId == blockID {
// 			found = true
// 		}
// 	}

// 	return result
// }

func getNextSCTDataAfterBlockID(sctDataReplies []SCTDataReply, blockID string) *SCTDataReply {
	for i, data := range sctDataReplies {
		if data.BlockId == blockID && i+1 < len(sctDataReplies) {
			return &sctDataReplies[i+1] // Return the next entry
		}
	}
	return nil // Return nil if no matching block or no next entry
}

// Handler function for /callback/nft
func ftDappHandler(c *gin.Context) {
	var req ContractInputRequest
	fmt.Println("Handler trggered")
	// cfg, err := config.GetConfig()
	// if err != nil {
	// 	fmt.Println("failed to load config: %w", err)
	// }
	// config.GetNodeNameByPort(cfg, req.Port)
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s", req.Port)
	fmt.Println("The url is :", url)

	// // config := GetConfig()
	smartContractHash := req.SmartContractHash
	fmt.Println("Received Smart Contract hash: ", smartContractHash)

	smartContractTokenData := rubix_interaction.GetSmartContractData(smartContractHash, url) //config.NodeAddress)
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
	wasmPath, err := getWasmContractPath(smartContractHash, req.Port)
	if err != nil {
		fmt.Println("Failed to get wasm path")
	}
	// Initialize the WASM module

	wasmModule, err := wasmbridge.NewWasmModule(
		wasmPath,
		hostFnRegistry,
		wasmbridge.WithRubixNodeAddress(url), //config.NodeAddress),
		wasmbridge.WithQuorumType(2),
	)
	if err != nil {
		log.Printf("Failed to initialize WASM module: %v", err)
		return
	}

	executionResult, errExecuteContract := executeAndGetContractResult(wasmModule, relevantData)
	fmt.Println("----------- FT Execution Result: ", executionResult)
	if errExecuteContract != nil {
		fmt.Println("The executionResult is ", executionResult)
		return
	}

	var response RubixResponse

	// Convert JSON string to struct
	if executionResult == "success" {
		response = RubixResponse{Status: true, Message: "FT Transferred Succesfully"}
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

// Handler function for /callback/nft
func ftContract2Handler(c *gin.Context) {
	var req ContractInputRequest
	fmt.Println("Handler trggered")
	// cfg, err := config.GetConfig()
	// if err != nil {
	// 	fmt.Println("failed to load config: %w", err)
	// }

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s", req.Port)
	fmt.Println("The url is :", url)
	// // config := GetConfig()
	smartContractHash := req.SmartContractHash
	fmt.Println("Received Smart Contract hash: ", smartContractHash)

	smartContractTokenData := rubix_interaction.GetSmartContractData(smartContractHash, url) //config.NodeAddress)
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
	wasmPath, err := getWasmContractPath(smartContractHash, req.Port)
	if err != nil {
		fmt.Println("Failed to get wasm path")
	}
	// Initialize the WASM module

	wasmModule, err := wasmbridge.NewWasmModule(
		wasmPath,
		hostFnRegistry,
		wasmbridge.WithRubixNodeAddress(url), //config.NodeAddress),
		wasmbridge.WithQuorumType(2),
	)
	if err != nil {
		log.Printf("Failed to initialize WASM module: %v", err)
		return
	}

	executionResult, errExecuteContract := executeAndGetContractResult(wasmModule, relevantData)
	fmt.Println("----------- FT Execution Result: ", executionResult)
	if errExecuteContract != nil {
		fmt.Println("The executionResult is ", executionResult)
		return
	}

	var response RubixResponse

	// Convert JSON string to struct
	if executionResult == "success" {
		response = RubixResponse{Status: true, Message: "FT Transferred Succesfully"}
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

func getWasmContractPath(contractHash, port string) (string, error) {
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}
	fmt.Println("The current working Directory is:", currentWorkingDir)
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println("Failed to get config file")
	}
	nodeName, exists := config.GetNodeNameByPort(cfg, port)
	if !exists {
		fmt.Println("Failed to get node name associated with the port", port)
	}
	// Construct the path in a cleaner way
	contractDir := filepath.Join(currentWorkingDir, "rubix-nodes", nodeName, "SmartContract", contractHash)
	fmt.Println("The contract directory is:", contractDir)

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
