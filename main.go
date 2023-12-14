package main

import (
	"user-management-servie/api"
	"user-management-servie/database"
	"user-management-servie/server"
)

func main() {
	dbClient := database.SetupDatabase()
	defer dbClient.Close()

	proxy := api.NewProxy(dbClient)
	server.SetupServer(proxy)
}
