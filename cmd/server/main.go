package main

import (
	"email-app/api"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load .env file
	if err := godotenv.Load("configs/.env"); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8083"
	}

	http.HandleFunc("/send-email", api.SendEmailHandler)
	fmt.Println("Email application started")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
