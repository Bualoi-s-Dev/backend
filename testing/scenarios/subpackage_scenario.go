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
	"github.com/Bualoi-s-Dev/backend/utils"
	"github.com/cucumber/godog"
)

type SubpackageScenario struct {
	Server     *httptest.Server
	Token      string
	Package    *dto.PackageResponse
	Subpackage *models.Subpackage
}

func (s *SubpackageScenario) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^the server is running$`, theServerIsRunning(s.Server))
	ctx.Given(`^the photographer is logged in$`, s.thePhotographerIsLoggedIn)
	ctx.Given(`^the photographer has a package$`, s.thePhotographerHasAPackage)

	ctx.When(`^the photographer creates a subpackage$`, s.thePhotographerCreatesASubpackage)
	ctx.When(`^the photographer updates a subpackage$`, s.thePhotographerUpdatesASubpackage)
	ctx.When(`^the photographer deletes a subpackage$`, s.thePhotographerDeletesASubpackage)

	ctx.Then(`^the subpackage is created and added to the package$`, s.theSubpackageIsCreated)
	ctx.Then(`^the subpackage is updated$`, s.theSubpackageIsUpdated)
	ctx.Then(`^the subpackage is deleted$`, s.theSubpackageIsDeleted)
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

	var subpackage models.Subpackage
	if err := json.NewDecoder(res.Body).Decode(&subpackage); err != nil {
		return err
	}
	s.Subpackage = &subpackage
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

func (s *SubpackageScenario) theSubpackageIsCreated() error {
	expect := models.Subpackage{
		PackageID:          s.Package.ID,
		Title:              "dev",
		Description:        "1234556",
		Price:              123,
		Duration:           23,
		IsInf:              false,
		RepeatedDay:        []models.DayName{models.Sunday, models.Wednesday},
		AvaliableStartTime: "15:11",
		AvaliableEndTime:   "16:00",
		AvaliableStartDay:  "2022-12-22",
		AvaliableEndDay:    "2023-01-22",
	}
	if err := utils.CompareStructsExcept(expect, *s.Subpackage, []string{"ID"}); err != nil {
		return err
	}
	return nil
}

func (s *SubpackageScenario) theSubpackageIsUpdated() error {
	expect := models.Subpackage{
		PackageID:          s.Package.ID,
		Title:              "dev123",
		Description:        "1234556",
		Price:              123,
		Duration:           23,
		IsInf:              true,
		RepeatedDay:        []models.DayName{models.Sunday, models.Wednesday},
		AvaliableStartTime: "15:11",
		AvaliableEndTime:   "16:00",
		AvaliableStartDay:  "2022-12-22",
		AvaliableEndDay:    "2023-01-22",
	}
	if err := utils.CompareStructsExcept(expect, *s.Subpackage, []string{"ID"}); err != nil {
		return err
	}
	return nil
}

func (s *SubpackageScenario) theSubpackageIsDeleted() error {
	req, err := http.NewRequest("GET", s.Server.URL+"/subpackage", nil)
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

	var subPackages []models.Subpackage
	if err := json.NewDecoder(res.Body).Decode(&subPackages); err != nil {
		return err
	}
	for _, subPackage := range subPackages {
		if subPackage.ID.Hex() == s.Subpackage.ID.Hex() {
			return fmt.Errorf("subpackage still exists")
		}
	}
	return nil
}
