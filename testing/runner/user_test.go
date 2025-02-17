package testing_runner

import (
	"testing"

	features "github.com/Bualoi-s-Dev/backend/testing/features"
	utils "github.com/Bualoi-s-Dev/backend/testing/utils"
)

func TestUserFeatures(t *testing.T) {
	server := GetTestServer()

	scenario := &features.UserScenario{Server: server}
	testSuite := utils.SetupGodog("user.feature", scenario.InitializeScenario)
	status := testSuite.Run()
	if status != 0 {
		t.Errorf("Non-zero exit code: %d", status)
		t.Fail()
	}
}
