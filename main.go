package main

import (
	"fmt"
	"net/http"
	"splitwise/api/handlers"
	"splitwise/api/routes"
	"splitwise/config"
)

func main() {
	db := config.ConnectDB()
	config.CreateTables(db)
	h := &handlers.Handler{DB: db}
	r := routes.SetupRoutes(h)
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
