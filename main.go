package main

import (
	"fmt"
	"log"
	"net/http"
	"postgres-go/router"
)

func main() {
	// запуск сервера
	r := router.Router()
	fmt.Println("Server is listening at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
