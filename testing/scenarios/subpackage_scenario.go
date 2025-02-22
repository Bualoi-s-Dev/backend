package testing_scenarios

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/cucumber/godog"
)

type SubpackageScenario struct {
	Server     *httptest.Server
	Token      string
	Package    *dto.PackageResponse
	Subpackage *models.Subpackage

	Response *http.Response
}

func (s *SubpackageScenario) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^the server is running$`, theServerIsRunning(s.Server))
	ctx.Given(`^the photographer is logged in$`, s.thePhotographerIsLoggedIn)
	ctx.Given(`^the photographer has a package$`, s.thePhotographerHasAPackage)

	ctx.When(`^the photographer creates a subpackage$`, s.thePhotographerCreatesASubpackage)
	ctx.When(`^the photographer updates a subpackage$`, s.thePhotographerUpdatesASubpackage)
	ctx.When(`^the photographer deletes a subpackage$`, s.thePhotographerDeletesASubpackage)

	ctx.Then(`^the subpackage is (created and added to the package|updated|deleted)$`, s.theSubpackageResponseIsOK)
}

func (s *SubpackageScenario) thePhotographerIsLoggedIn() error {
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

func (s *SubpackageScenario) thePhotographerHasAPackage() error {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"title":  "Dev Dol Package",
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
	s.Response = res

	var packageResponse dto.PackageResponse
	if err := json.NewDecoder(res.Body).Decode(&packageResponse); err != nil {
		return err
	}
	s.Package = &packageResponse

	return nil
}

func (s *SubpackageScenario) thePhotographerCreatesASubpackage() error {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"title":       "dev",
		"description": "1234556",
		"price":       123,
		"duration":    23,

		"isInf":              true,
		"repeatedDay":        []string{"SUN", "WED"},
		"avaliableStartTime": "15:11",
		"avaliableEndTime":   "16:00",
		"avaliableStartDay":  "2022-12-22",
		"avaliableEndDay":    "2023-01-22",
	})
	req, err := http.NewRequest("POST", s.Server.URL+"/subpackage/"+s.Package.ID.Hex(), bytes.NewBuffer(reqBody))
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
		return fmt.Errorf("failed to create subpackage, status code: %d", res.StatusCode)
	}

	var subpackage models.Subpackage
	if err := json.NewDecoder(res.Body).Decode(&subpackage); err != nil {
		return err
	}
	s.Subpackage = &subpackage
	return nil
}

func (s *SubpackageScenario) thePhotographerUpdatesASubpackage() error {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"title": "dev123",
	})
	req, err := http.NewRequest("PATCH", s.Server.URL+"/subpackage/"+s.Subpackage.ID.Hex(), bytes.NewBuffer(reqBody))
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
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update subpackage, status code: %d", res.StatusCode)
	}
	return nil
}

func (s *SubpackageScenario) thePhotographerDeletesASubpackage() error {
	req, err := http.NewRequest("DELETE", s.Server.URL+"/subpackage/"+s.Subpackage.ID.Hex(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete subpackage, status code: %d", res.StatusCode)
	}
	return nil
}

func (s *SubpackageScenario) theSubpackageResponseIsOK() error {
	if s.Response.StatusCode == http.StatusOK || s.Response.StatusCode == http.StatusCreated {
		return nil
	}
	return fmt.Errorf("expected status code 200 or 201, got %d", s.Response.StatusCode)
}
