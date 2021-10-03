package main

import (
	"fmt"
	"github.com/chezky/library/db"
	"github.com/chezky/library/mail"
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
	mail.Start()

	r := mux.NewRouter()

	r.HandleFunc("/book/new", routes.NewBookHandler)
	r.HandleFunc("/book/get", routes.GetBooksHandler)
	r.HandleFunc("/book/delete", routes.DeleteBookHandler)
	r.HandleFunc("/book/checkout", routes.CheckOutBookHandler)
	r.HandleFunc("/book/get/id", routes.GetBookByIDHandler)
	r.HandleFunc("/book/update", routes.UpdateBooksHandler)
	r.HandleFunc("/book/search", routes.SearchBooksHandler)
	r.HandleFunc("/account/new", routes.NewAccountHandler)
	r.HandleFunc("/account/get", routes.GetAccountsHandler)
	r.HandleFunc("/account/delete", routes.DeleteAccountHandler)
	r.HandleFunc("/account/get/id", routes.GetAccountByIDHandler)
	r.HandleFunc("/account/update", routes.UpdateAccountHandler)
	r.HandleFunc("/account/search", routes.SearchAccountsHandler)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", r)

	port := fmt.Sprintf(":%d", 8080)
	fmt.Println("server running on port", port)

	log.Fatal(http.ListenAndServe(port, r))
}
