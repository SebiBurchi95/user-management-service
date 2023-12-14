package database

import (
	"context"
	"log"
	"user-management-servie/ent"

	"github.com/spf13/viper"
)

func SetupDatabase() *ent.Client {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	dbType := viper.GetString("database.type")
	dataSource := viper.GetString("database.dataSource")

	client, err := ent.Open(dbType, dataSource)
	if err != nil {
		log.Fatalf("failed opening connection to %s: %v", dbType, err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}
