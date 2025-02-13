package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	rubix "dapp-server/rubix-interaction"

	"github.com/gin-gonic/gin"
)

type ExecuteRequest struct {
	ContractHash      string
	ExecutorDid       string
	HomeDirectory     string
	ContractDirectory string
}

type DeployRequest struct {
	WasmPath    string
	LibPath     string
	DeployerDid string
	StatePath   string
}

func APIExecuteContract(c *gin.Context) {
	var req ExecuteRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	result, err := rubix.Execute(req.ContractHash, req.ExecutorDid, req.HomeDirectory, req.ContractDirectory, "")
	if err != nil {
		fmt.Println("Failed to execute Contract err :", err)
	}
	fmt.Println("The result returned : ", result)
	resultFinal := gin.H{
		"message": "DApp executed successfully",
		"data":    result,
	}

	// Return a response
	c.JSON(http.StatusOK, resultFinal)
}

func APIDeployContract(c *gin.Context) {
	var req DeployRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	result, err := rubix.Deploy(req.WasmPath, req.LibPath, req.DeployerDid, req.StatePath)
	if err != nil {
		fmt.Println("Failed to deploy contract err :", err)
	}
	fmt.Println("The result returned : ", result)
	resultFinal := gin.H{
		"message": "Contract Deployed Successfully",
		"data":    result,
	}
	c.JSON(http.StatusOK, resultFinal)
}
