package testing_scenarios

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/cucumber/godog"
)

type PackageScenario struct {
	Server *httptest.Server
	Token  string
}

func (s *PackageScenario) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^the server is running$`, theServerIsRunning(s.Server))
	ctx.Given(`^the photographer is logged in$`, s.thePhotographerLoggedIn)
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
