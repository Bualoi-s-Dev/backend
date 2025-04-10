package testing_runner

import (
	"context"
	"testing"
	"time"

	scenarios "github.com/Bualoi-s-Dev/backend/testing/scenarios"
	utils "github.com/Bualoi-s-Dev/backend/testing/utils"
)

func TestAppointmentFeatures(t *testing.T) {
	defer cleanUpAppointmentFeature()
	server := GetTestServer()

	scenario := &scenarios.AppointmentScenario{Server: server}
	testSuite := utils.SetupGodog("appointment.feature", scenario.InitializeScenario)
	status := testSuite.Run()
	if status != 0 {
		t.Errorf("Non-zero exit code: %d", status)
		t.Fail()
	}
}

func cleanUpAppointmentFeature() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := GetTestMongoDB()
	db.Collection("Appointment").Drop(ctx)
	db.Collection("Subpackage").Drop(ctx)
	db.Collection("Package").Drop(ctx)
	db.Collection("BusyTime").Drop(ctx)
}
