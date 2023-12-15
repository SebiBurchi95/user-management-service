package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"user-management-servie/ent"
	"user-management-servie/ent/enttest"

	"github.com/alecthomas/assert"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func setupTest(t *testing.T) (*ent.Client, Proxy) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	proxy := NewProxy(client)
	return client, proxy
}

func TestCreateUser(t *testing.T) {
	client, proxy := setupTest(t)
	defer client.Close()

	user := requestUser{
		Username: "test",
		Email:    "testCreate@gmail.com",
	}
	requestBody, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	respRecorder := httptest.NewRecorder()

	proxy.CreateUser(respRecorder, req)

	assert.Equal(t, http.StatusCreated, respRecorder.Code)
	var createdUser ent.User
	err := json.NewDecoder(respRecorder.Body).Decode(&createdUser)
	assert.NoError(t, err)
	assert.Equal(t, "test", createdUser.Username)
	assert.Equal(t, "testCreate@gmail.com", createdUser.Email)
}

func TestGetUser(t *testing.T) {
	client, proxy := setupTest(t)
	defer client.Close()

	mockUser, err := client.User.
		Create().
		SetUsername("test").
		SetEmail("testCreate@gmail.com").
		Save(context.Background())
	if err != nil {
		t.Fatalf("Failed to create mock user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(mockUser.ID), nil)
	respRecorder := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", proxy.GetUser)
	r.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusOK, respRecorder.Code)

	var retrievedUser userResponse
	err = json.NewDecoder(respRecorder.Body).Decode(&retrievedUser)
	assert.NoError(t, err)

	assert.Equal(t, mockUser.ID, retrievedUser.ID)
	assert.Equal(t, mockUser.Username, retrievedUser.Username)
	assert.Equal(t, mockUser.Email, retrievedUser.Email)

	assert.NotEmpty(t, retrievedUser.DogPhotoURL)
}

func TestUpdateUser(t *testing.T) {
	client, proxy := setupTest(t)
	defer client.Close()

	mockUser, err := client.User.
		Create().
		SetUsername("test").
		SetEmail("testUpdate@gmail.com").
		Save(context.Background())
	if err != nil {
		t.Fatalf("Failed to create mock user: %v", err)
	}

	updatedData := requestUser{
		Username: "update",
		Email:    "update@gmail.com",
	}
	requestBody, _ := json.Marshal(updatedData)
	req := httptest.NewRequest(http.MethodPut, "/users/"+strconv.Itoa(mockUser.ID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	respRecorder := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", proxy.UpdateUser)
	r.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusOK, respRecorder.Code)

	var updatedUser ent.User
	err = json.NewDecoder(respRecorder.Body).Decode(&updatedUser)
	assert.NoError(t, err)
	assert.Equal(t, "update", updatedUser.Username)
	assert.Equal(t, "update@gmail.com", updatedUser.Email)
	assert.Equal(t, mockUser.ID, updatedUser.ID)
}

func TestListUsers(t *testing.T) {
	client, proxy := setupTest(t)
	defer client.Close()

	var wg sync.WaitGroup

	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := client.User.
				Create().
				SetUsername(fmt.Sprintf("test%d", i)).
				SetEmail(fmt.Sprintf("test%d@gmail.com", i)).
				Save(context.Background())
			if err != nil {
				t.Errorf("Failed to create mock user %d: %v", i, err)
			}
		}(i)
	}

	wg.Wait()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	respRecorder := httptest.NewRecorder()

	proxy.ListUsers(respRecorder, req)

	assert.Equal(t, http.StatusOK, respRecorder.Code)

	var users []*ent.User
	err := json.NewDecoder(respRecorder.Body).Decode(&users)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestDeleteUser(t *testing.T) {
	client, proxy := setupTest(t)
	defer client.Close()

	mockUser, err := client.User.
		Create().
		SetUsername("testUser").
		SetEmail("testUser@gmail.com").
		Save(context.Background())
	if err != nil {
		t.Fatalf("Failed to create mock user: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%d", mockUser.ID), nil)
	respRecorder := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", proxy.DeleteUser)
	r.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusOK, respRecorder.Code)

	var result map[string]string
	err = json.NewDecoder(respRecorder.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "success", result["result"])

	_, err = client.User.Get(context.Background(), mockUser.ID)
	assert.True(t, ent.IsNotFound(err), "User should be deleted")
}
