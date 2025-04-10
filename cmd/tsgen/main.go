package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	outputFile := "gen/api_types.ts"
	converter := typescriptify.New()

	// Add models
	converter.
		Add(models.Package{}).
		Add(dto.PackageRequest{}).
		Add(dto.PackageResponse{}).
		AddEnum(models.ValidPackageTypes)
	converter.
		Add(dto.UserRequest{}).
		Add(dto.UserResponse{}).
		Add(dto.CheckProviderResponse{}).
		AddEnum(models.ValidUserRoles).
		AddEnum(models.ValidBankNames)
	converter.
		Add(dto.SubpackageRequest{}).
		Add(dto.SubpackageResponse{}).
		AddEnum(models.ValidDayNames)
	converter.
		Add(models.BusyTime{}).
		Add(dto.BusyTimeStrictRequest{}).
		AddEnum(models.ValidBusyTimeTypes)
	converter.
		Add(dto.AppointmentRequest{}).
		Add(dto.AppointmentStrictRequest{}).
		Add(dto.AppointmentUpdateStatusRequest{}).
		Add(dto.AppointmentResponse{}).
		Add(dto.AppointmentDetail{}).
		Add(dto.CreateAppointmentResponse{}).
		AddEnum(models.ValidAppointmentStatus)
	converter.
		Add(dto.RatingRequest{}).
		Add(dto.RatingResponse{})
	converter.
		Add(dto.PaymentResponse{}).
		Add(dto.PaymentURL{}).
		AddEnum(models.ValidPaymentStatus)

	// Change to interface
	converter.CreateInterface = true

	// Create dir
	dirErr := os.MkdirAll("gen/backup", os.ModePerm)
	if dirErr != nil {
		log.Fatalf("Error creating directory: %v", dirErr)
		return
	}

	// Backup
	converter.BackupDir = "gen/backup"

	// Convert
	err := converter.ConvertToFile(outputFile)

	if err != nil {
		log.Fatalf("Error generating TypeScript types: %v", err)
	}

	fmt.Printf("TypeScript types successfully generated at: %s\n", outputFile)
}
