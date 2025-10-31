package main

import (
	"log"
	"os"

	// Blank-import the function package so the init() runs
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/peterhoward42/dxact-analytics"
)

func main() {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	hostname := ""
	if err := funcframework.StartHostPort(hostname, port); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v\n", err)
	}
}
