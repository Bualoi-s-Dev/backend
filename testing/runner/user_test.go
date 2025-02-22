package testing_runner

import (
	"testing"

	scenarios "github.com/Bualoi-s-Dev/backend/testing/scenarios"
	utils "github.com/Bualoi-s-Dev/backend/testing/utils"
)

func TestUserFeatures(t *testing.T) {
	server := GetTestServer()

	scenario := &scenarios.UserScenario{Server: server}
	testSuite := utils.SetupGodog("user.feature", scenario.InitializeScenario)
	status := testSuite.Run()
	if status != 0 {
		t.Errorf("Non-zero exit code: %d", status)
		t.Fail()
	}
}
