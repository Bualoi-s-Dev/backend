package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	outputFile := "gen/api_types.ts"
	converter := typescriptify.New()

	// Add models
	converter.
		Add(models.Package{}).
		AddEnum(models.ValidPackageTypes)
	converter.
		Add(models.User{}).
		AddEnum(models.ValidBankNames)

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
