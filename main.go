package main

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/routers"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	database.StartDB()

	var PORT = helpers.GetEnv("PORT")

	fmt.Printf("Server is running on port %s", PORT)

	routers.StartServer().Run(":" + PORT)
}
