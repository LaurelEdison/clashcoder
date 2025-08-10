package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	godotenv.Load(".env")
	http.HandleFunc("/", greet)
	http.ListenAndServe(":8080", nil)
}
