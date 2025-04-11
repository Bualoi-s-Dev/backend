package testing_runner

import (
	"context"
	"testing"
	"time"

	scenarios "github.com/Bualoi-s-Dev/backend/testing/scenarios"
	utils "github.com/Bualoi-s-Dev/backend/testing/utils"
)

func TestPackageFeatures(t *testing.T) {
	defer cleanUpPackageFeature()
	server := GetTestServer()

	scenario := &scenarios.PackageScenario{Server: server}
	testSuite := utils.SetupGodog("package.feature", scenario.InitializeScenario)
	status := testSuite.Run()
	if status != 0 {
		t.Errorf("Non-zero exit code: %d", status)
		t.Fail()
	}
}

func cleanUpPackageFeature() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := GetTestMongoDB()
	db.Collection("Package").Drop(ctx)
}
