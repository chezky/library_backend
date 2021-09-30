package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chezky/library/db"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func NewBookHandler(w http.ResponseWriter, r *http.Request)  {
	var b db.Book

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request")
		return
	}

	if err := json.Unmarshal(body, &b); err != nil {
		fmt.Println("error unmarshalling new book", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := b.NewBook(); err != nil{
		fmt.Println("error inserting book into db", b.Title, err)
	}

	fmt.Println(fmt.Sprintf("Title is: %s, author is %s", b.Title, b.Author))
	fmt.Fprint(w, fmt.Sprintf("%d", b.ID))
}

func DeleteBookHandler(w http.ResponseWriter, r *http.Request)  {
	var b db.Book

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &b); err != nil {
		fmt.Println("error unmarshalling new book", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := b.DeleteBook(); err != nil{
		fmt.Println("error deleting book", b.ID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "success")
}

func CheckOutBookHandler(w http.ResponseWriter, r *http.Request)  {
	var b db.Book

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &b); err != nil {
		fmt.Println("error unmarshalling new book", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("customerID is:", b.CustomerID)

	if err := b.CheckOutBook(); err != nil{
		fmt.Println("error checking out book", b.ID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "success")
}

func GetBookByIDHandler(w http.ResponseWriter, r *http.Request)  {
	var b db.Book

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &b); err != nil {
		fmt.Println("error unmarshalling get book by ID", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := b.GetBookByID(); err != nil{
		fmt.Println("error getting book", b.ID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(b)
	if err != nil {
		fmt.Println("error marshaling in book", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(res))
}

//UpdateBooksHandler either returns or checks out a list of book. Requires a name to be given when the books are being checked out.
func UpdateBooksHandler(w http.ResponseWriter, r *http.Request)  {
	var b db.BookList

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &b); err != nil {
		fmt.Println("error unmarshalling get book by ID", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b.UpdateListBooks()

	fmt.Fprint(w, "success")
}

//SearchBooksHandler either returns or checks out a list of book. Requires a name to be given when the books are being checked out.
func SearchBooksHandler(w http.ResponseWriter, r *http.Request)  {
	var s db.Search

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &s); err != nil {
		fmt.Println("error unmarshalling search books by title", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.SearchByTitle(); err != nil {
		fmt.Println("error searching by title for query", s.Query, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(s.Books)
	if err != nil {
		fmt.Println("error marshaling in books for search by title", s.Query, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(res))
}

func GetBooksHandler(w http.ResponseWriter, r *http.Request)  {
	b, err := db.GetBooks()
	if err != nil {
		fmt.Println("error getting books", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(b)
	if err != nil {
		fmt.Println("error marshaling in books", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(res))
}

func NewAccountHandler(w http.ResponseWriter, r *http.Request)  {
	var a db.Account

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request")
		return
	}

	if err := json.Unmarshal(body, &a); err != nil {
		fmt.Println("error unmarshalling new account", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.NewAccount(); err != nil{
		fmt.Println("error inserting new account into db", a.Name, err)
	}

	fmt.Println(fmt.Sprintf("New account created: * %s * was added to the library.", a.Name))
	fmt.Fprint(w, fmt.Sprintf("%d", a.ID))
}

func GetAccountsHandler(w http.ResponseWriter, r *http.Request)  {
	b, err := db.GetAccounts()
	if err != nil {
		fmt.Println("error getting accounts", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(b)
	if err != nil {
		fmt.Println("error marshaling in accounts", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(res))
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request)  {
	var a db.Account

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &a); err != nil {
		fmt.Println("error unmarshalling for delete account", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.DeleteAccount(); err != nil{
		fmt.Println("error deleting account id number:", a.ID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("successfully deleted account number:", a.ID)

	fmt.Fprint(w, "success")
}

func GetAccountByIDHandler(w http.ResponseWriter, r *http.Request)  {
	var a db.Account

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &a); err != nil {
		fmt.Println("error unmarshalling get book by ID", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.GetAccountByID(); err != nil{
		fmt.Println("error getting account number:", a.ID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(a)
	if err != nil {
		fmt.Println("error marshaling in account", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(res))
}

func UpdateAccountHandler(w http.ResponseWriter, r *http.Request)  {
	var a db.Account

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request")
		return
	}

	if err := json.Unmarshal(body, &a); err != nil {
		fmt.Println("error unmarshalling update account", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.UpdateAccount(); err != nil{
		fmt.Println("error inserting updated account into db:", a.Name, err)
	}

	fmt.Println(fmt.Sprintf("Account updated: * %s * was updated in the library.", a.Name))
	fmt.Fprint(w, "success")
}

func SearchAccountsHandler(w http.ResponseWriter, r *http.Request)  {
	var s db.Search

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &s); err != nil {
		fmt.Println("error unmarshalling search accounts", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.SearchAccounts(); err != nil {
		fmt.Println("error searching for accounts for query", s.Query, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(s.Accounts)
	if err != nil {
		fmt.Println("error marshaling in accounts for search", s.Query, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(res))
}