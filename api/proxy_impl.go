package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"user-management-servie/ent"
	"user-management-servie/util"

	"github.com/gorilla/mux"
)

type proxyImpl struct {
	client *ent.Client
}

type requestUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type userResponse struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DogPhotoURL string    `json:"dog_photo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateUser implements Proxy.
func (p *proxyImpl) CreateUser(resp http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	var reqUser requestUser
	if err := json.NewDecoder(req.Body).Decode(&reqUser); err != nil {
		http.Error(resp, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdUser, err := p.client.User.
		Create().
		SetUsername(reqUser.Username).
		SetEmail(reqUser.Email).
		Save(ctx)

	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusCreated)
	json.NewEncoder(resp).Encode(createdUser)
}

// DeleteUser implements Proxy.
func (p *proxyImpl) DeleteUser(resp http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	userID, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		http.Error(resp, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = p.client.User.DeleteOneID(userID).Exec(ctx)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(map[string]string{"result": "success"})
}

// GetUser implements Proxy.
func (p *proxyImpl) GetUser(resp http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	userID, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		http.Error(resp, "Invalid user ID", http.StatusBadRequest)
		return
	}

	u, err := p.client.User.Get(ctx, userID)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	dogPhotoURL, err := util.GetRandomDogPhotoURL()
	if err != nil {
		log.Printf("Failed to fetch dog photo: %v", err)
	}

	userResp := userResponse{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		DogPhotoURL: dogPhotoURL,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}

	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(userResp)
}

// ListUsers implements Proxy.
func (p *proxyImpl) ListUsers(resp http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	users, err := p.client.User.Query().All(ctx)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	// Channel to collect users
	userCh := make(chan *ent.User)

	// Concurrently fetch users
	for _, u := range users {
		go func(u *ent.User) {
			userCh <- u
		}(u)
	}

	// Collect users from the channel
	var userList []*ent.User
	for i := 0; i < len(users); i++ {
		userList = append(userList, <-userCh)
	}

	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(userList)
}

// UpdateUser implements Proxy.
func (p *proxyImpl) UpdateUser(resp http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	userID, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		http.Error(resp, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var reqUser requestUser
	if err := json.NewDecoder(req.Body).Decode(&reqUser); err != nil {
		http.Error(resp, "Invalid request body", http.StatusBadRequest)
		return
	}

	u, err := p.client.User.
		UpdateOneID(userID).
		SetUsername(reqUser.Username).
		SetEmail(reqUser.Email).
		Save(ctx)

	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(u)
}

func NewProxy(client *ent.Client) Proxy {
	return &proxyImpl{client: client}
}
