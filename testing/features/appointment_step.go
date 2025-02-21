package testing_features

// import (
// 	"net/http/httptest"

// 	"github.com/cucumber/godog"
// )

// type AppointmentScenario struct {
// 	Server                      *httptest.Server
// 	CustomerLoggedIn            bool
// 	SubPackageSelected          bool
// 	TimeSlotSelected            bool
// 	AppointmentCreated          bool
// 	CustomerScheduleUpdated     bool
// 	PhotographerLoggedIn        bool
// 	AppointmentScheduled        bool
// 	AppointmentCompleted        bool
// 	PhotographerScheduleUpdated bool
// 	AppointmentCanceled         bool
// }

// func (s *AppointmentScenario) InitializeScenario(ctx *godog.ScenarioContext) {
// 	ctx.Step(`^the server is running$`, theServerIsRunning(s.Server))

// 	ctx.Step(`^a customer is logged in$`)

// 	ctx.Step(`^a sub package is selected$`)

// 	ctx.Step(`^a time slot is selected$`)

// 	ctx.Step(`^the customer submits the appointment request$`)

// 	ctx.Step(`^the appointment is created$`)

// 	ctx.Step(`^the customer’s schedule is updated$`)

// 	ctx.Step(`^a photographer is logged in$`)

// 	ctx.Step(`^an appointment is scheduled$`)

// 	ctx.Step(`^the photographer marks the appointment as completed$`)

// 	ctx.Step(`^the appointment status is updated to completed$`)

// 	ctx.Step(`^the photographer’s schedule is updated$`)

// 	ctx.Step(`^the photographer cancels the appointment$`)

// 	ctx.Step(`^the appointment is canceled$`)

// 	ctx.Step(`^the customer cancels the appointment$`)
// }
