package database

import (
	"context"
	"testing"
	"user-management-servie/ent/enttest"

	"github.com/alecthomas/assert"
)

func TestSetupDatabase(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	dbClient := SetupDatabase()
	defer dbClient.Close()

	assert.NotNil(t, dbClient, "Database client should not be nil")

	testUser, err := dbClient.User.
		Create().
		SetUsername("test").
		SetEmail("test@gmail.com").
		Save(context.Background())

	assert.NoError(t, err, "Should be able to create a user")
	assert.NotNil(t, testUser, "Created user should not be nil")

	assert.Equal(t, "test", testUser.Username, "Username should match")
	assert.Equal(t, "test@gmail.com", testUser.Email, "Email should match")
}
