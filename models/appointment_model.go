package models

import "time"

type Appointment struct {
	ID             string    `json:"id,omitempty"`
	UserID         string    `json:"user_id"`
	PhotographerID string    `json:"photographer_id"`
	PackageID      string    `json:"package_id"`
	SubPackageID   string    `json:"sub_package_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Status         string    `json:"status"` // "pending", "accepted", "rejected", "completed"
	// Payment       Payment   `json:"payment,omitempty"`
}
