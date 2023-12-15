#!/bin/bash

# Base URL of your API
BASE_URL="http://localhost:8080"

# Function to make a POST request to create a user
function createUser() {
    curl -X POST "$BASE_URL/users" \
         -H "Content-Type: application/json" \
         -d '{"username": "testUser", "email": "test@example.com"}'
    echo
}

# Function to get a user by ID
function getUser() {
    local id=$1
    curl -X GET "$BASE_URL/users/$id"
    echo
}

# Function to update a user by ID
function updateUser() {
    local id=$1
    curl -X PUT "$BASE_URL/users/$id" \
         -H "Content-Type: application/json" \
         -d '{"username": "updatedUser", "email": "updated@example.com"}'
    echo
}

# Function to delete a user by ID
function deleteUser() {
    local id=$1
    curl -X DELETE "$BASE_URL/users/$id"
    echo
}

# Function to list all users
function listUsers() {
    curl -X GET "$BASE_URL/users"
    echo
}

# Simulate API calls
echo "Creating User..."
createUser

echo "Getting User with ID 1..."
getUser 1

echo "Updating User with ID 1..."
updateUser 1

echo "Listing Users..."
listUsers

echo "Deleting User with ID 1..."
deleteUser 1

# End of script
