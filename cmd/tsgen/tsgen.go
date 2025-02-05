package main

import (
	"fmt"
	"log"

	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	outputFile := "api_types.ts"
	converter := typescriptify.New()

	converter.Add(models.Package{})

	converter.CreateInterface = true

	err := converter.ConvertToFile(outputFile)

	if err != nil {
		log.Fatalf("Error generating TypeScript types: %v", err)
	}

	fmt.Printf("TypeScript types successfully generated at: %s\n", outputFile)
}
