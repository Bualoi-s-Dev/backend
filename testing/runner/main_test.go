package testing_runner

import (
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Bualoi-s-Dev/backend/bootstrap"
	"github.com/Bualoi-s-Dev/backend/configs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var testServer *httptest.Server
var testMongoDB *mongo.Database

func TestMain(m *testing.M) {
	rootDir, err := filepath.Abs("../../") // Adjust path to root
	if err != nil {
		log.Fatalf("Failed to get root directory: %v", err)
	}
	if err := os.Chdir(rootDir); err != nil {
		log.Fatalf("Failed to change directory: %v", err)
	}

	configs.LoadEnv()
	// Start test server
	testServer = startTestServer()

	code := m.Run()

	// Stop the test server after tests
	stopTestServer(testServer)
	os.Exit(code)
}

func startTestServer() *httptest.Server {
	log.Println("Starting test server...")

	databaseName := "Testing"
	gin.SetMode(gin.TestMode)

	client := configs.ConnectMongoDB().Database(databaseName)

	testMongoDB = client
	r, _, _ := bootstrap.SetupServer(client, true)

	testServer := httptest.NewServer(r)
	return testServer
}

func stopTestServer(server *httptest.Server) {
	log.Println("Stopping test server...")
	if server != nil {
		server.Close()
	}
}

func GetTestServer() *httptest.Server {
	return testServer
}

func GetTestMongoDB() *mongo.Database {
	return testMongoDB
}
