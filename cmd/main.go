package main

import (
	"log"
	"os"

	// Blank-import the function package so the init() runs
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/peterhoward42/dxact-analytics"
)
// This program starts a local development web server that serves a google cloud function on localhost, 
// using functions-framework-go
//
// It's not obvious how this code knows what function to run! The answer is as follows:
// - the function framework expects the function to be in the root package.
// - In our case the injestEvent() function in ./injestevent.go.
// - But you have to provide a name binding to that function also...
// - Which in our case is a name binding to "InjestEvent", and is declared in the init() function for the root
//   package in the same file. 
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
