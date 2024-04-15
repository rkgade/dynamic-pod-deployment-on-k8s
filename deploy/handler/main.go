package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type HashData struct {
	Input string `json:"input"`
	Hash  string `json:"hash"`
}

var (
	hashList []HashData
	mutex    sync.Mutex
)

func generateHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", hash)
}

func generateHashHandler(c *gin.Context) {
	input := c.Query("input")
	if input == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is required"})
		return
	}

	hashValue := generateHash(input)

	mutex.Lock()
	hashList = append(hashList, HashData{Input: input, Hash: hashValue})
	mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"hash": hashValue})
}

func getAllHashesHandler(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	c.JSON(http.StatusOK, hashList)
}

func main() {
	router := gin.Default()

	router.POST("/generate", generateHashHandler)
	router.GET("/hashes", getAllHashesHandler)

	router.Run(":8080")
}
