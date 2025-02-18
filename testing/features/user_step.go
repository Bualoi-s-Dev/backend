package testing_features

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/Bualoi-s-Dev/backend/configs"
	"github.com/cucumber/godog"
)

type UserScenario struct {
	Server   *httptest.Server
	Username string
	Password string
	Response *http.Response
}

func (s *UserScenario) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^the server is running$`, theServerIsRunning(s.Server))

	ctx.Given(`^valid credentials are provided$`, s.validCredentialsAreProvided)
	ctx.Given(`^invalid credentials are provided$`, s.invalidCredentialsAreProvided)

	ctx.When(`^the login is submitted$`, s.loginIsSubmitted)
	ctx.When(`^the user attempts to log in$`, s.loginIsSubmitted)

	ctx.Then(`^access to the account is granted$`, s.accessToAccountIsGranted)
	ctx.Then(`^the system should reject the login and display an error message saying "([^"]*)"$`, s.errorMessageIsDisplayed)
}

func (s *UserScenario) validCredentialsAreProvided() error {
	s.Username = configs.GetEnv("TEST_USER_EMAIL")
	s.Password = configs.GetEnv("TEST_USER_PASSWORD")
	return nil
}

func (s *UserScenario) invalidCredentialsAreProvided() error {
	s.Username = configs.GetEnv("TEST_USER_EMAIL")
	s.Password = configs.GetEnv("TEST_USER_PASSWORD") + "invalid"
	return nil
}

func (s *UserScenario) loginIsSubmitted() error {
	reqBody, _ := json.Marshal(map[string]string{
		"email":    s.Username,
		"password": s.Password,
	})
	res, err := http.Post(s.Server.URL+"/internal/firebase/login", "application/json", bytes.NewBuffer(reqBody))
	s.Response = res
	return err
}

func (s *UserScenario) accessToAccountIsGranted() error {
	if s.Response.StatusCode != http.StatusOK {
		return errors.New("expected status code 200, got " + s.Response.Status)
	}
	return nil
}

func (s *UserScenario) errorMessageIsDisplayed(expectedMessage string) error {
	if s.Response.StatusCode != http.StatusUnauthorized {
		return errors.New("expected status code 401, got " + s.Response.Status)
	}
	return nil
}
