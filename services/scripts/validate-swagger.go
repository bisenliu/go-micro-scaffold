package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// SwaggerDoc represents the basic structure of a Swagger document
type SwaggerDoc struct {
	Swagger string                 `json:"swagger"`
	Info    SwaggerInfo            `json:"info"`
	Host    string                 `json:"host"`
	Paths   map[string]interface{} `json:"paths"`
}

type SwaggerInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

func main() {
	// Read the generated swagger.json file
	data, err := ioutil.ReadFile("docs/swagger.json")
	if err != nil {
		log.Fatalf("Failed to read swagger.json: %v", err)
	}

	// Parse the JSON
	var doc SwaggerDoc
	if err := json.Unmarshal(data, &doc); err != nil {
		log.Fatalf("Failed to parse swagger.json: %v", err)
	}

	// Validate basic structure
	fmt.Println("=== Swagger Documentation Validation ===")
	fmt.Printf("Swagger Version: %s\n", doc.Swagger)
	fmt.Printf("API Title: %s\n", doc.Info.Title)
	fmt.Printf("API Description: %s\n", doc.Info.Description)
	fmt.Printf("API Version: %s\n", doc.Info.Version)
	fmt.Printf("Host: %s\n", doc.Host)
	fmt.Printf("Number of API paths: %d\n", len(doc.Paths))

	// List all available paths
	fmt.Println("\n=== Available API Paths ===")
	for path := range doc.Paths {
		fmt.Printf("- %s\n", path)
	}

	// Check for required paths
	requiredPaths := []string{
		"/auth/login",
		"/auth/logout",
		"/auth/wechat",
		"/users",
		"/users/{id}",
		"/health",
	}

	fmt.Println("\n=== Required Paths Validation ===")
	allFound := true
	for _, path := range requiredPaths {
		if _, exists := doc.Paths[path]; exists {
			fmt.Printf("✓ %s - Found\n", path)
		} else {
			fmt.Printf("✗ %s - Missing\n", path)
			allFound = false
		}
	}

	if allFound {
		fmt.Println("\n✅ All required API paths are documented!")
		os.Exit(0)
	} else {
		fmt.Println("\n❌ Some required API paths are missing!")
		os.Exit(1)
	}
}
