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
	converter.Add(models.Package{})

	// Change to interface
	converter.CreateInterface = true

	// Create dir
	dirErr := os.MkdirAll("gen", os.ModePerm)
	if dirErr != nil {
		log.Fatalf("Error creating directory: %v", dirErr)
		return
	}

	// Convert
	err := converter.ConvertToFile(outputFile)

	if err != nil {
		log.Fatalf("Error generating TypeScript types: %v", err)
	}

	fmt.Printf("TypeScript types successfully generated at: %s\n", outputFile)
}
