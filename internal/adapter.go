package internal

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

func ConnectDatabase() *sqlx.DB {
	var databaseServer = struct {
		dbUsername string
		dbPassword string
		dbPort     string
		dbName     string
		dbServer   string
	}{
		dbUsername: os.Getenv("DB_USERNAME"),
		dbPassword: os.Getenv("DB_PASSWORD"),
		dbPort:     os.Getenv("DB_PORT"),
		dbName:     os.Getenv("DB_NAME"),
		dbServer:   os.Getenv("DB_SERVER"),
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
