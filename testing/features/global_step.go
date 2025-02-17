package testing_features

import (
	"fmt"
	"net"
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
