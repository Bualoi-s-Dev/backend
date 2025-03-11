package testing_runner

// import (
// 	"testing"

// 	features "github.com/Bualoi-s-Dev/backend/testing/features"
// 	utils "github.com/Bualoi-s-Dev/backend/testing/utils"
// )

// func TestAppointmentFeatures(t *testing.T) {
// 	server := GetTestServer()

// 	scenario := &features.AppointmentScenario{Server: server}
// 	testSuite := utils.SetupGodog("appointment.feature", scenario.InitializeScenario)
// 	status := testSuite.Run()
// 	if status != 0 {
// 		t.Errorf("Non-zero exit code: %d", status)
// 		t.Fail()
// 	}
// }
