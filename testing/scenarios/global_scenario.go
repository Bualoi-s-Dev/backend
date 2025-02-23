package testing_scenarios

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func theServerIsRunning(server *httptest.Server) func() error {
	return func() error {
		if server == nil {
			return fmt.Errorf("test server is not initialized")
		}

		// Parse the server URL to extract the host and port
		parsedURL, err := url.Parse(server.URL)
		if err != nil {
			return fmt.Errorf("invalid test server URL: %v", err)
		}

		// Check if the server is listening on the port
		conn, err := net.Dial("tcp", parsedURL.Host)
		if err != nil {
			return fmt.Errorf("test server is not running: %v", err)
		}
		defer conn.Close()
		return nil
	}

}

func getLoginToken(server *httptest.Server, email string, password string) (string, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	res, err := http.Post(server.URL+"/internal/firebase/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var responseBody map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return "", err
	}

	token, ok := responseBody["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response body")
	}
	return token, nil
}
