package server

import (
	"log"
	"net/http"
	"user-management-servie/api"

	"github.com/gorilla/mux"
)

func SetupServer(proxy api.Proxy) {
	router := mux.NewRouter()

	router.HandleFunc("/users", proxy.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", proxy.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", proxy.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", proxy.DeleteUser).Methods("DELETE")
	router.HandleFunc("/users", proxy.ListUsers).Methods("GET")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
