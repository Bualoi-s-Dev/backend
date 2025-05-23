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
	Subpackage *dto.SubpackageResponse
}

func (s *SubpackageScenario) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^the server is running$`, theServerIsRunning(s.Server))
	ctx.Given(`^the photographer is logged in$`, s.thePhotographerIsLoggedIn)
	ctx.Given(`^the photographer has a package$`, s.thePhotographerHasAPackage)

	ctx.When(`^the photographer creates a subpackage$`, s.thePhotographerCreatesASubpackage)
	ctx.When(`^the photographer updates a subpackage$`, s.thePhotographerUpdatesASubpackage)
	ctx.When(`^the photographer deletes a subpackage$`, s.thePhotographerDeletesASubpackage)
	ctx.When(`^the photographer creates a subpackage with wrong format$`, s.thePhotographerCreatesASubpackageWithWrongFormat)
	ctx.When(`^the photographer updates a subpackage with wrong format$`, s.thePhotographerUpdatesASubpackageWithWrongFormat)
	ctx.When(`^the photographer deletes a non-existent subpackage$`, s.thePhotographerDeletesANonExistentSubpackage)

	ctx.Then(`^the subpackage is created and added to the package$`, s.theSubpackageIsCreated)
	ctx.Then(`^the subpackage is updated$`, s.theSubpackageIsUpdated)
	ctx.Then(`^the subpackage is deleted$`, s.theSubpackageIsDeleted)
	ctx.Then(`^the subpackage is not created and not added to the package$`, s.theSubpackageIsNotCreated)
	ctx.Then(`^the subpackage is not updated$`, s.theSubpackageIsNotUpdated)
	ctx.Then(`^the subpackage is not deleted$`, s.theSubpackageIsNotDeleted)
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

		"isInf":              false,
		"repeatedDay":        []string{"SUN", "WED"},
		"availableStartTime": "15:11",
		"availableEndTime":   "16:00",
		"availableStartDay":  "2030-12-22",
		"availableEndDay":    "2031-01-22",
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
		return fmt.Errorf("failed to create subpackage, status code: %d, %s", res.StatusCode, res.Body)
	}

	var subpackage dto.SubpackageResponse
	if err := json.NewDecoder(res.Body).Decode(&subpackage); err != nil {
		return err
	}
	s.Subpackage = &subpackage
	return nil
}

func (s *SubpackageScenario) thePhotographerUpdatesASubpackage() error {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"title":             "dev123",
		"availableStartDay": "2035-12-22",
		"availableEndDay":   "2036-01-22",
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

	var subpackage dto.SubpackageResponse
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

func (s *SubpackageScenario) thePhotographerCreatesASubpackageWithWrongFormat() error {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"title":              "",
		"description":        "1234556",
		"price":              -10,     // Invalid negative price
		"duration":           "23",    // Invalid non-numeric duration
		"availableStartTime": "25:00", // Invalid time format
		"availableEndTime":   "16:00",
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

	if res.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("expected status 400 Bad Request, got: %d", res.StatusCode)
	}

	var errorResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
		return err
	}
	return nil
}

func (s *SubpackageScenario) thePhotographerUpdatesASubpackageWithWrongFormat() error {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"price": "-10",
	})

	req, err := http.NewRequest("PATCH", s.Server.URL+"/subpackage/"+s.Package.ID.Hex(), bytes.NewBuffer(reqBody))
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

	if res.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("expected status 400 Bad Request, got: %d", res.StatusCode)
	}

	return nil
}

func (s *SubpackageScenario) thePhotographerDeletesANonExistentSubpackage() error {
	req, err := http.NewRequest("DELETE", s.Server.URL+"/subpackage/nonexistentid", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected status 404 Not Found, got: %d", res.StatusCode)
	}

	return nil
}

func (s *SubpackageScenario) theSubpackageIsCreated() error {
	expect := dto.SubpackageResponse{
		PackageID:          s.Package.ID,
		Title:              "dev",
		Description:        "1234556",
		Price:              123,
		Duration:           23,
		IsInf:              false,
		RepeatedDay:        []models.DayName{models.Sunday, models.Wednesday},
		AvailableStartTime: "15:11",
		AvailableEndTime:   "16:00",
		AvailableStartDay:  "2030-12-22",
		AvailableEndDay:    "2031-01-22",
	}
	if err := utils.CompareStructsExcept(expect, *s.Subpackage, []string{"ID", "BusyTimes", "BusyTimeMap"}); err != nil {
		return err
	}
	return nil
}

func (s *SubpackageScenario) theSubpackageIsUpdated() error {
	expect := dto.SubpackageResponse{
		PackageID:          s.Package.ID,
		Title:              "dev123",
		Description:        "1234556",
		Price:              123,
		Duration:           23,
		IsInf:              false,
		RepeatedDay:        []models.DayName{models.Sunday, models.Wednesday},
		AvailableStartTime: "15:11",
		AvailableEndTime:   "16:00",
		AvailableStartDay:  "2035-12-22",
		AvailableEndDay:    "2036-01-22",
	}
	if err := utils.CompareStructsExcept(expect, *s.Subpackage, []string{"ID", "BusyTimes", "BusyTimeMap"}); err != nil {
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
		return fmt.Errorf("failed to get all subpackage, status code: %d", res.StatusCode)
	}

	// Decode the response
	var response dto.FilteredSubpackageResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}

	// Check if the deleted subpackage still exists
	for _, subPackage := range response.Subpackages {
		if subPackage.ID.Hex() == s.Subpackage.ID.Hex() {
			return fmt.Errorf("subpackage still exists with ID: %s", subPackage.ID.Hex())
		}
	}
	return nil
}

func (s *SubpackageScenario) theSubpackageIsNotCreated() error {
	req, err := http.NewRequest("GET", s.Server.URL+"/subpackage/"+s.Package.ID.Hex(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusOK {
		var subPackages []models.Subpackage
		if err := json.NewDecoder(res.Body).Decode(&subPackages); err != nil {
			return err
		}
		for _, subPackage := range subPackages {
			if subPackage.ID.Hex() == s.Package.ID.Hex() {
				return fmt.Errorf("subpackage was unexpectedly created")
			}
		}
	}

	return nil
}

func (s *SubpackageScenario) theSubpackageIsNotUpdated() error {
	req, err := http.NewRequest("GET", s.Server.URL+"/subpackage/"+s.Subpackage.ID.Hex(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusOK {
		var subpackage models.Subpackage
		if err := json.NewDecoder(res.Body).Decode(&subpackage); err != nil {
			return err
		}
		if subpackage.Price >= 0 {
			return nil
		}
		return fmt.Errorf("subpackage was unexpectedly updated")
	}

	return nil
}

func (s *SubpackageScenario) theSubpackageIsNotDeleted() error {
	req, err := http.NewRequest("GET", s.Server.URL+"/subpackage/"+s.Subpackage.ID.Hex(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusOK {
		var subpackage models.Subpackage
		if err := json.NewDecoder(res.Body).Decode(&subpackage); err != nil {
			return err
		}
		if subpackage.ID.Hex() == s.Subpackage.ID.Hex() {
			return fmt.Errorf("subpackage still exists")
		}
	}

	return nil
}
