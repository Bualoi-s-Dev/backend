package testing_scenarios

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/utils"
	"github.com/cucumber/godog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentScenario struct {
	Server         *httptest.Server
	Token          string
	Package        *dto.PackageResponse
	Subpackage     *dto.SubpackageResponse
	Appointment    *models.Appointment
	PhotographerID primitive.ObjectID
	CustomerID     primitive.ObjectID
}

func (s *AppointmentScenario) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^the server is running$`, theServerIsRunning(s.Server))
	ctx.Given(`^a photographer has a package and sub package$`, s.thePhotographerHasPackageAndSubpackage)
	ctx.Given(`^a customer is logged in$`, s.theCustomerIsLoggedIn)

	ctx.When(`^a customer creates an appointment$`, s.theCustomerCreatesAnAppointment)

	ctx.Then(`^the appointment is created$`, s.theAppointmentIsCreated)
}

func (s *AppointmentScenario) thePhotographerHasPackageAndSubpackage() error {
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
		"duration":           23,
		"isInf":              false,
		"repeatedDay":        []string{"SUN", "WED"},
		"avaliableStartTime": "15:11",
		"avaliableEndTime":   "16:00",
		"avaliableStartDay":  "2022-12-22",
		"avaliableEndDay":    "2023-01-22",
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
	s.PhotographerID = userProfile.ID

	return nil
}

func (s *AppointmentScenario) theCustomerIsLoggedIn() error {
	//Login and get the token
	token, err := getLoginToken(s.Server, os.Getenv("TEST_USER_EMAIL"), os.Getenv("TEST_USER_PASSWORD"))
	if err != nil {
		return err
	}
	s.Token = token // Store token separately

	//Update role to Customer
	reqBody, _ := json.Marshal(map[string]interface{}{
		"role": "Customer",
	})
	req, err := http.NewRequest("PATCH", s.Server.URL+"/user/profile", bytes.NewBuffer(reqBody))
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
		return fmt.Errorf("failed to update user role, status code: %d", res.StatusCode)
	}

	// Fetch the user profile to get the Customer ID
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

	// Decode response to extract Customer ID
	var userProfile struct {
		ID primitive.ObjectID `json:"id"`
	}
	if err := json.NewDecoder(res.Body).Decode(&userProfile); err != nil {
		return err
	}
	s.CustomerID = userProfile.ID
	return nil
}

func (s *AppointmentScenario) theCustomerCreatesAnAppointment() error {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"start_time": "2030-02-21T10:30:00.000+00:00",
		"location":   "Bangkok, Thailand",
	})
	req, err := http.NewRequest("POST", s.Server.URL+"/appointment"+"/"+s.Subpackage.ID.Hex(), bytes.NewBuffer(reqBody))
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
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to create appointment, status code: %d, response: %s", res.StatusCode, string(body))
	}

	var appointment models.Appointment
	if err := json.NewDecoder(res.Body).Decode(&appointment); err != nil {
		return err
	}
	fmt.Println("Created Appointment:", appointment) // Debugging log
	s.Appointment = &appointment
	return nil
}

func (s *AppointmentScenario) theAppointmentIsCreated() error {
	expect := models.Appointment{
		CustomerID:     s.CustomerID,
		PhotographerID: s.PhotographerID,
		PackageID:      s.Package.ID,
		SubpackageID:   s.Subpackage.ID,
		Status:         "Pending",
		Location:       "Bangkok, Thailand",
	}
	fmt.Println("Expect:", expect) // Debugging log
	fmt.Println("Actual Appointment:", *s.Appointment)
	fmt.Println("Actual Package:", *s.Package)
	fmt.Println("Actual Subpackage:", *s.Subpackage)
	fmt.Println("Actual Photographer ID:", s.PhotographerID)
	fmt.Println("Actual Customer ID:", s.CustomerID)
	if err := utils.CompareStructsExcept(expect, *s.Appointment, []string{"ID", "BusyTimeID"}); err != nil {
		return err
	}
	return nil
}
