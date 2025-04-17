package testing_scenarios

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/cucumber/godog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackageScenario struct {
	Server     *httptest.Server
	Token      string
	Package    *dto.PackageResponse
	Subpackage *dto.SubpackageResponse
}

func (s *PackageScenario) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^the server is running$`, theServerIsRunning(s.Server))
	ctx.Given(`^a photographer has a package and sub package$`, s.thePhotographerHasPackageAndSubpackage)
	ctx.Given(`^a photographer is logged in$`, s.thePhotographerLoggedIn)

	ctx.When(`^the photographer updates the package details$`, s.thePhotographerUpdatesThePackageDetails)
	ctx.When(`^the photographer deletes the package$`, s.thePhotographerDeletesThePackage)

	ctx.Then(`^the package information is updated$`, s.thePackageInformationIsUpdated)
	ctx.Then(`^the package is removed$`, s.thePackageIsRemoved)
}

func (s *PackageScenario) thePhotographerLoggedIn() error {
	token, err := getLoginToken(s.Server, os.Getenv("TEST_PHOTOGRAPHER_EMAIL"), os.Getenv("TEST_PHOTOGRAPHER_PASSWORD"))
	if err != nil {
		return err
	}
	s.Token = token

	reqBody, _ := json.Marshal(map[string]interface{}{
		"role": "Photographer",
	})
	req, err := http.NewRequest("PATCH", s.Server.URL+"/user/profile", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (s *PackageScenario) thePhotographerHasPackageAndSubpackage() error {
	token, err := getLoginToken(s.Server, os.Getenv("TEST_PHOTOGRAPHER_EMAIL"), os.Getenv("TEST_PHOTOGRAPHER_PASSWORD"))
	if err != nil {
		return err
	}
	s.Token = token

	reqBody, _ := json.Marshal(map[string]interface{}{
		"title":  "Photography Package",
		"type":   "OTHER",
		"photos": []string{},
	})
	req, err := http.NewRequest("POST", s.Server.URL+"/package", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create package, status code: %d", res.StatusCode)
	}

	var packageResponse dto.PackageResponse
	if err := json.NewDecoder(res.Body).Decode(&packageResponse); err != nil {
		return err
	}
	s.Package = &packageResponse

	reqBody, _ = json.Marshal(map[string]interface{}{
		"title":              "dev",
		"description":        "1234556",
		"price":              123,
		"duration":           30,
		"isInf":              false,
		"repeatedDay":        []string{"SUN", "WED"},
		"availableStartTime": "15:00",
		"availableEndTime":   "16:00",
		"availableStartDay":  "2030-12-22",
		"availableEndDay":    "2031-01-22",
	})
	req, err = http.NewRequest("POST", s.Server.URL+"/subpackage/"+s.Package.ID.Hex(), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err = client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create subpackage, status code: %d", res.StatusCode)
	}

	var subpackage dto.SubpackageResponse
	if err := json.NewDecoder(res.Body).Decode(&subpackage); err != nil {
		return err
	}
	s.Subpackage = &subpackage

	// Fetch the user profile to get the Photographer ID
	req, err = http.NewRequest("GET", s.Server.URL+"/user/profile", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)

	res, err = client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch user profile, status code: %d", res.StatusCode)
	}

	// Decode response to extract Photographer ID
	var userProfile struct {
		ID primitive.ObjectID `json:"id"`
	}
	if err := json.NewDecoder(res.Body).Decode(&userProfile); err != nil {
		return err
	}

	return nil
}

func (s *PackageScenario) thePhotographerUpdatesThePackageDetails() error {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"title": "Updated Photography Package",
		"type":  "WEDDING_BLISS",
	})
	req, err := http.NewRequest("PATCH", s.Server.URL+"/package/"+s.Package.ID.Hex(), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update package, status code: %d", res.StatusCode)
	}

	var updatedPackage dto.PackageResponse
	if err := json.NewDecoder(res.Body).Decode(&updatedPackage); err != nil {
		return err
	}
	s.Package = &updatedPackage

	return nil
}

func (s *PackageScenario) thePackageInformationIsUpdated() error {
	req, err := http.NewRequest("GET", s.Server.URL+"/package/"+s.Package.ID.Hex(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch package, status code: %d", res.StatusCode)
	}

	var fetchedPackage dto.PackageResponse
	if err := json.NewDecoder(res.Body).Decode(&fetchedPackage); err != nil {
		return err
	}

	if fetchedPackage.Title != "Updated Photography Package" || fetchedPackage.Type != "WEDDING_BLISS" {
		return fmt.Errorf("package information was not updated correctly")
	}

	return nil
}

func (s *PackageScenario) thePhotographerDeletesThePackage() error {
	req, err := http.NewRequest("DELETE", s.Server.URL+"/package/"+s.Package.ID.Hex(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete package, status code: %d", res.StatusCode)
	}

	// Clear the package reference after deletion
	// s.Package = nil
	return nil
}

func (s *PackageScenario) thePackageIsRemoved() error {
	req, err := http.NewRequest("GET", s.Server.URL+"/package/"+s.Package.ID.Hex(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check if the package no longer exists
	if res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("package was not removed, status code: %d", res.StatusCode)
	}

	return nil
}
