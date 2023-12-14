package server

import (
	"net/http"
	"testing"
	"user-management-servie/api"
	"user-management-servie/ent/enttest"

	"gotest.tools/assert"
)

func TestServer(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	proxy := api.NewProxy(client)

	go SetupServer(proxy)

	httpClient := &http.Client{}

	testCases := []struct {
		name       string
		method     string
		endpoint   string
		statusCode int
	}{
		{"CreateUser", "POST", "/users", http.StatusCreated},
		{"GetUser", "GET", "/users/1", http.StatusOK},
		{"UpdateUser", "PUT", "/users/1", http.StatusOK},
		{"DeleteUser", "DELETE", "/users/1", http.StatusOK},
		{"ListUsers", "GET", "/users", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, "http://localhost:8080"+tc.endpoint, nil)
			resp, err := httpClient.Do(req)
			if err != nil {
				t.Errorf("Failed to make request: %v", err)
			}
			assert.Equal(t, tc.statusCode, resp.StatusCode)
		})
	}
}
