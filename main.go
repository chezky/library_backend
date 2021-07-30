package main

import (
	"fmt"
	"github.com/chezky/library/db"
	"github.com/chezky/library/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("No .env file found")
	}

	db.Start()

	r := mux.NewRouter()

	r.HandleFunc("/new", routes.NewBookHandler)
	r.HandleFunc("/get", routes.GetBooksHandler)
	r.HandleFunc("/delete", routes.DeleteBookHandler)
	r.HandleFunc("/checkout", routes.CheckOutBookHandler)
	r.HandleFunc("/get/id", routes.GetBookByIDHandler)
	r.HandleFunc("/update", routes.UpdateBooksHandler)
	r.HandleFunc("/search/title", routes.SearchByTitleHandler)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", r)

	port := fmt.Sprintf(":%d", 8080)
	fmt.Println("server running on port", port)
	log.Fatal(http.ListenAndServe(port, r))
}