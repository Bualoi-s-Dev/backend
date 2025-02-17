package main

import (
	"github.com/Bualoi-s-Dev/backend/bootstrap"
	"github.com/Bualoi-s-Dev/backend/configs"
)

// @title PhotoMatch API
// @version 1.0
// @description This is a sample API to demonstrate Swagger with Gin.
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @host      localhost:8080

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Check the -test flag
	// testMode := flag.Bool("test", false, "Run in test mode")
	// flag.Parse()

	databaseName := "PhotoMatch"
	// if *testMode {
	// 	fmt.Println("Running in test mode...")
	// 	databaseName = "Testing"
	// }

	configs.LoadEnv()

	client := configs.ConnectMongoDB().Database(databaseName)
	// fmt.Println("Using %s Database", databaseName)
	r := bootstrap.SetupServer(client)
	r.Run()
}
