package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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

// func Bootup() {
// 	fmt.Println("Server Started")
// 	r := mux.NewRouter()

// 	r.HandleFunc("/api/v1/contract-input", contractInputHandler).Methods("POST")
// 	err := http.ListenAndServe(":8080", r)
// 	if err != nil {
// 		fmt.Printf("Error starting server: %s\n", err)
// 	}
// }

func BootupServer() {
	// Initialize a Gin router
	router := gin.Default()
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
	// folderPath, _ := GetRubixSmartContractPath(req.SmartContractHash, "binaryCodeFile.wasm", nodeName)
	// schemaPath, _ := GetRubixSchemaPath(req.SmartContractHash, nodeName, "schemaCodeFile.json")
	// fmt.Println(folderPath)
	// _, err1 := os.Stat(folderPath)
	// fmt.Println(err1)
	// if os.IsNotExist(err1) {
	// 	fmt.Println("Smart Contract not found")
	// 	RunSmartContract(folderPath, schemaPath, port, req.SmartContractHash)
	// } else if err == nil {
	// 	fmt.Printf("Folder '%s' exists", folderPath)

	// 	RunSmartContract(folderPath, schemaPath, port, req.SmartContractHash)

	// } else {
	// 	fmt.Printf("Error while checking folder: %v\n", err)
	// }

	resp := RubixResponse{Status: true, Message: "Callback Successful", Result: "Success"}
	json.NewEncoder(w).Encode(resp)

}

// Handler function for /callback/nft
func ftDappHandler(c *gin.Context) {
	var req ContractInputRequest
	fmt.Println("Jandler trggered")
	err := json.NewDecoder(c.Request.Body).Decode(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// // config := GetConfig()
	// smartContractHash := req.SmartContractHash
	// fmt.Println("Received Smart Contract hash: ", req.SmartContractHash)

	// smartContractTokenData := rubix.GetSmartContractData(smartContractHash, "") //config.NodeAddress)
	// if smartContractTokenData == nil {
	// 	fmt.Println("Unable to fetch latest smart contract data")
	// 	return
	// }

	// fmt.Println("Smart Contract Token Data :", string(smartContractTokenData))

	// var dataReply SmartContractDataReply

	// if err := json.Unmarshal(smartContractTokenData, &dataReply); err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// fmt.Println("Data reply in runDappHandler", dataReply)
	// smartContractData := dataReply.SCTDataReply
	// var relevantData string
	// for _, reply := range smartContractData {
	// 	fmt.Println("SmartContractData:", reply.SmartContractData)
	// 	relevantData = reply.SmartContractData
	// }
	// contractInput := `{"mint_sample_ft":{"name": "rubix1", "ft_info": {
	// 	"did": "bafybmihxaehnreq4ygnq3re3soob5znuj7hxoku6aeitdukif75umdv2nu",
	// 	"ft_count": 100,
	// 	"ft_name": "test5",
	// 	"token_count": 1
	//   }}}`
	relevantData := `{"mint_sample_ft":{"name": "rubix1", "ft_info": {
		"did": "bafybmieqhv5zd7m7mmtoigqupqg2si2ri2d3fuqf43p5affuagufxgyen4",
		"ft_count": 100,
		"ft_name": "test5",
		"token_count": 1
	  }}}`
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
	// var requestId string
	// switch funcName {
	// case "mint_sample_ft":
	// 	requestId = "ft-" + smartContractHash + "-mint"
	// case "transfer_sample_ft":
	// 	requestId = "ft-" + smartContractHash + "-transfer"
	// default:
	// 	fmt.Println("This function name is not allowed")
	// 	return
	// }
	// checkResult, err := checkStringInRequests(requestId)
	// if err != nil {
	// 	fmt.Println("Error checking result:", err)
	// 	return
	// }
	// if !checkResult {
	// 	err = insertRequest(requestId, Pending) //Add constants for the status
	// 	if err != nil {
	// 		fmt.Println("Error inserting request:", err)
	// 		return
	// 	}
	// }

	hostFnRegistry := wasmbridge.NewHostFunctionRegistry()

	// Initialize the WASM module
	wasmModule, err := wasmbridge.NewWasmModule(
		// config.ContractsInfo["ft"].ContractPath,
		"C:/Users/allen/Working-repo/ymca/ymca-wellness-cafe-project/first-contract/target/wasm32-unknown-unknown/debug/first_contract.wasm",
		hostFnRegistry,
		wasmbridge.WithRubixNodeAddress("http://localhost:20003"), //config.NodeAddress),
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
		// func() {
		// 	err = updateRequestStatus(requestId, Failed)
		// 	if err != nil {
		// 		fmt.Println("Error updating request status:", err)
		// 		return
		// 	}
		// }()
	}

	// if response.Status {
	// 	err = updateRequestStatus(requestId, Success)
	// 	if err != nil {
	// 		fmt.Println("Error updating request status:", err)
	// 		return
	// 	} //handle error here
	// } else {
	// 	err = updateRequestStatus(requestId, Failed)
	// 	if err != nil {
	// 		fmt.Println("Error updating request status:", err)
	// 		return
	// 	}
	// }
	resultFinal := gin.H{
		"message": "DApp executed successfully",
		"data":    response,
	}

	// Return a response
	c.JSON(http.StatusOK, resultFinal)
}

func executeAndGetContractResult(wasmModule *wasmbridge.WasmModule, contractInput string) (string, error) {
	// Call the function
	contractResult, err := wasmModule.CallFunction(contractInput)
	if err != nil {
		return "", fmt.Errorf("function call failed: %v", err)
	}

	return contractResult, nil
}
