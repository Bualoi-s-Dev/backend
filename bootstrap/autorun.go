package bootstrap

import (
	"context"
	"os"
	"time"
)

func AutoUpdate(ctx context.Context, serverService *ServerServices) error {
	var tickerTime time.Duration
	if os.Getenv("APP_MODE") == "development" {
		// tickerTime = 1 * time.Minute
		tickerTime = 60 * time.Second
	} else {
		tickerTime = 15 * time.Minute
	}
	ticker := time.NewTicker(tickerTime) // Runs every 15 minutes
	defer ticker.Stop()

	for {
		<-ticker.C
		go serverService.appointmentService.AutoUpdateAppointmentStatus(ctx)
	}
}
