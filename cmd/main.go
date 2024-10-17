package main

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	var dbConnString string

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	flag.StringVar(&dbConnString, "db", "", "Database connection string")
	flag.Parse()

}
