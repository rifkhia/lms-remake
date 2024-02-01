package internal

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func ConnectDatabase() *sqlx.DB {
	var databaseServer = struct {
		dbUsername string
		dbPassword string
		dbPort     string
		dbName     string
		dbServer   string
	}{
		dbUsername: viper.GetString("DB_USERNAME"),
		dbPassword: viper.GetString("DB_PASSWORD"),
		dbPort:     viper.GetString("DB_PORT"),
		dbName:     viper.GetString("DB_NAME"),
		dbServer:   viper.GetString("DB_SERVER"),
	}
	if databaseServer.dbUsername == "" {
		fmt.Println("error in connecting")

	}
	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		databaseServer.dbUsername,
		databaseServer.dbPassword,
		databaseServer.dbServer,
		databaseServer.dbPort,
		databaseServer.dbName))

	if err != nil {
		fmt.Printf("Failed to connect database, err:", err)
	}

	return db
}
