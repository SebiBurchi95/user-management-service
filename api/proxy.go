package api

import "net/http"

type Proxy interface {
	CreateUser(resp http.ResponseWriter, req *http.Request)
	GetUser(resp http.ResponseWriter, req *http.Request)
	UpdateUser(resp http.ResponseWriter, req *http.Request)
	DeleteUser(resp http.ResponseWriter, req *http.Request)
	ListUsers(resp http.ResponseWriter, req *http.Request)
}
