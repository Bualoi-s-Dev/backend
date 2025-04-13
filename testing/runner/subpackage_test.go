package testing_runner

import (
	"context"
	"testing"
	"time"

	scenarios "github.com/Bualoi-s-Dev/backend/testing/scenarios"
	utils "github.com/Bualoi-s-Dev/backend/testing/utils"
)

func TestSubpackageFeatures(t *testing.T) {
	defer cleanUpSubpackageFeature()
	server := GetTestServer()

	scenario := &scenarios.SubpackageScenario{Server: server}
	testSuite := utils.SetupGodog("subpackage.feature", scenario.InitializeScenario)
	status := testSuite.Run()
	if status != 0 {
		t.Errorf("Non-zero exit code: %d", status)
	}
}

func cleanUpSubpackageFeature() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := GetTestMongoDB()
	db.Collection("Subpackage").Drop(ctx)
	db.Collection("Package").Drop(ctx)
}
