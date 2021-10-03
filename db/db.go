package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" //importing blank as per the package recommendation.
	"os"
	"strings"
	"time"
)

var (
	db *sql.DB
)

type Book struct {
	ID         int32  `json:"id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Available  bool   `json:"available"`
	Customer   string `json:"customer"`
	CustomerID int32  `json:"customer_id""`
	TimeStamp  int64  `json:"time_stamp"`
}

type BookList struct {
	IDs        []int32 `json:"ids"`
	Customer   string  `json:"customer"`
	CustomerID int32   `json:"customer_id"`
	Available  bool    `json:"available"`
	TimeStamp  int64   `json:"time_stamp"`
}

type Account struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	EmailList bool   `json:"email_list"`
	LastEmail int64  `json:"last_email"`
	Books     []Book `json:"books"`
	BookCount int    `json:"book_count"`
}

type Search struct {
	Books    []Book    `json:"books"`
	Accounts []Account `json:"accounts"`
	Query    string    `json:"query"`
}

func Start() {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", "localhost", 5432, "postgres", os.Getenv("POSTGRES_PASS"), "library")

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	createBooksTable()
	createAccountsTable()

	fmt.Println("SQL Connected")
}

func (b *Book) NewBook() error {
	msg := `
	INSERT INTO books (title, author, ts)
	VALUES ($1, $2, $3)		
	RETURNING id
	`

	if err := db.QueryRow(msg, b.Title, b.Author, time.Now().Unix()).Scan(&b.ID); err != nil {
		return err
	}

	return nil
}

func (b *Book) DeleteBook() error {
	msg := `DELETE from books WHERE id = $1`

	if _, err := db.Exec(msg, b.ID); err != nil {
		return err
	}
	return nil
}

func (b *Book) CheckOutBook() error {
	var (
		ts  int64
		msg string
	)

	if !b.Available {
		msg = `SELECT name FROM accounts WHERE id=$1`
		if err := db.QueryRow(msg, b.CustomerID).Scan(&b.Customer); err != nil {
			if !strings.Contains(err.Error(), "converting NULL to string") {
				return err
			}
		}

		ts = time.Now().Unix()
	}

	msg = `UPDATE books SET available = $1, customer=$2, customer_id=$3, ts = $4 WHERE id = $5`

	if _, err := db.Exec(msg, b.Available, b.Customer, b.CustomerID, ts, b.ID); err != nil {
		return err
	}
	return nil
}

func (b *Book) GetBookByID() error {
	msg := `SELECT title, available, author FROM books WHERE id = $1`

	if err := db.QueryRow(msg, b.ID).Scan(&b.Title, &b.Available, &b.Author); err != nil {
		if !strings.Contains(err.Error(), "converting NULL to string") {
			return err
		}
	}

	if !b.Available {
		msg = `SELECT customer, customer_id, ts FROM books WHERE id = $1`
		if err := db.QueryRow(msg, b.ID).Scan(&b.Customer, &b.CustomerID, &b.TimeStamp); err != nil {
			return err
		}
	}

	return nil
}

// UpdateListBooks only runs to check out books, not return
func (b *BookList) UpdateListBooks() {
	msg := `SELECT name FROM accounts WHERE id=$1`
	if err := db.QueryRow(msg, b.CustomerID).Scan(&b.Customer); err != nil {
		fmt.Println("error getting name for id #:", b.CustomerID)
	}

	for _, id := range b.IDs {
		msg = `UPDATE books SET available = FALSE, customer_id=$1, customer = $2, ts = $3 WHERE id = $4`
		if _, err := db.Exec(msg, b.CustomerID, b.Customer, time.Now().Unix(), id); err != nil {
			fmt.Println("error updating list book for id", id, err)
		}
	}
}

func (s *Search) SearchByTitle() error {
	msg := `SELECT title, id, available, author FROM books WHERE title ILIKE $1 Or author ILIKE $1 OR customer ILIKE $1 order by title`

	rows, err := db.Query(msg, "%"+s.Query+"%")
	if err != nil {
		fmt.Println("error getting books", err)
		return err
	}

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.Title, &book.ID, &book.Available, &book.Author); err != nil {
			if !strings.Contains(err.Error(), "converting NULL to string is unsupported") {
				fmt.Println("error scanning get search by title", err)
			}
		}

		if !book.Available {
			msg = `SELECT customer, customer_id, ts FROM books WHERE id = $1`
			if err := db.QueryRow(msg, book.ID).Scan(&book.Customer, &book.CustomerID, &book.TimeStamp); err != nil {
				if !strings.Contains(err.Error(), "converting NULL to string is unsupported") {
					fmt.Println("error getting customer for book ID", book.ID, err)
				}
			}
		}
		s.Books = append(s.Books, book)
	}

	return nil
}

func GetBooks() ([]Book, error) {
	var b []Book

	msg := `SELECT title, id, available, author FROM books order by Title`

	rows, err := db.Query(msg)
	if err != nil {
		fmt.Println("error getting books", err)
		return b, err
	}

	for rows.Next() {
		var book Book

		if err := rows.Scan(&book.Title, &book.ID, &book.Available, &book.Author); err != nil {
			if !strings.Contains(err.Error(), "converting NULL to string is unsupported") {
				fmt.Println("error scanning get transfers", err)
			}
		}

		if !book.Available {
			msg = `SELECT customer, customer_id, ts FROM books WHERE id = $1`
			if err := db.QueryRow(msg, book.ID).Scan(&book.Customer, &book.CustomerID, &book.TimeStamp); err != nil {
				if !strings.Contains(err.Error(), "converting NULL to string is unsupported") {
					fmt.Println("error getting customer for book ID", book.ID, err)
				}
			}
		}
		b = append(b, book)
	}

	return b, nil
}

func (a *Account) NewAccount() error {
	msg := `
	INSERT INTO accounts (name, email, last_email)
	VALUES ($1, $2, $3)		
	RETURNING id
	`

	if err := db.QueryRow(msg, a.Name, a.Email, time.Now().Unix()).Scan(&a.ID); err != nil {
		return err
	}

	return nil
}

func GetAccounts() ([]Account, error) {
	var a []Account

	msg := `SELECT id, name, email FROM accounts order by name`

	rows, err := db.Query(msg)
	if err != nil {
		fmt.Println("error getting accounts", err)
		return a, err
	}

	for rows.Next() {
		var account Account

		if err := rows.Scan(&account.ID, &account.Name, &account.Email); err != nil {
			if !strings.Contains(err.Error(), "converting NULL to string is unsupported") {
				fmt.Println("error scanning get accounts", err)
			}
		}

		msg = `SELECT COUNT(*) FROM books WHERE customer_id = $1 AND available = FALSE`
		if err := db.QueryRow(msg, account.ID).Scan(&account.BookCount); err != nil {
			fmt.Println("error getting count of books per customer")
		}

		a = append(a, account)
	}

	return a, nil
}

func (a *Account) DeleteAccount() error {
	msg := `DELETE from accounts WHERE id = $1`

	if _, err := db.Exec(msg, a.ID); err != nil {
		return err
	}
	return nil
}

func (a *Account) GetAccountByID() error {
	msg := `SELECT name, email, email_list, last_email FROM accounts WHERE id = $1`

	if err := db.QueryRow(msg, a.ID).Scan(&a.Name, &a.Email, &a.EmailList, &a.LastEmail); err != nil {
		if !strings.Contains(err.Error(), "converting NULL to string") {
			return err
		}
	}

	msg = `SELECT title, id, author, ts FROM books WHERE customer_id = $1 AND available = $2`

	rows, err := db.Query(msg, a.ID, false)
	if err != nil {
		fmt.Println("error getting books", err)
		return err
	}

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.Title, &book.ID, &book.Author, &book.TimeStamp); err != nil {
			if !strings.Contains(err.Error(), "converting NULL to string is unsupported") {
				fmt.Println("error scanning for get Account books", err)
			}
		}

		book.CustomerID = a.ID

		a.Books = append(a.Books, book)
	}

	a.BookCount = len(a.Books)

	return nil
}

func (a *Account) UpdateAccount() error {
	msg := `UPDATE accounts SET name = $1, email = $2, email_list = $3 WHERE id = $4`

	if _, err := db.Exec(msg, a.Name, a.Email, a.EmailList, a.ID); err != nil {
		return err
	}
	return nil
}

func (s *Search) SearchAccounts() error {
	msg := `SELECT id, name, email, email_list FROM accounts WHERE name ILIKE $1 Or email ILIKE $1 ORDER BY name`

	rows, err := db.Query(msg, "%"+s.Query+"%")
	if err != nil {
		fmt.Println("error getting accounts", err)
		return err
	}

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.Name, &account.Email, &account.EmailList); err != nil {
			if !strings.Contains(err.Error(), "converting NULL to string is unsupported") {
				fmt.Println("error scanning search for accounts", err)
			}
		}

		msg = `SELECT COUNT(*) FROM books WHERE customer_id = $1`
		if err := db.QueryRow(msg, account.ID).Scan(&account.BookCount); err != nil {
			fmt.Println("error getting count of books per customer")
		}

		s.Accounts = append(s.Accounts, account)
	}

	return nil
}

func FindEmailAccounts() []Account {
	var accounts []Account

	fmt.Println("Looking for accounts to email...")

	// 604800 is one week
	// 1209600 is two weeks
	// Select an account id, if they have a book checked out for over two weeks, and if the account email_list is turned on, and if the account has not received an email in over a week.
	msg := `SELECT DISTINCT accounts.id FROM accounts JOIN books ON accounts.id = books.customer_id WHERE books.available = FALSE and email_list = TRUE and last_email + 604800 < $1 and books.ts + 1209600 < $1`
	rows, err := db.Query(msg, time.Now().Unix())
	if err != nil {
		fmt.Println("error getting accounts that need emails", err)
		return accounts
	}

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID); err != nil {
			fmt.Println("error getting id for account during email scan", err)
			continue
		}

		if err := account.GetAccountByID(); err != nil {
			fmt.Println("error getting account info for account id:", account.ID, err)
			continue
		}

		// update the account to reflect that an email has been sent
		msg = `UPDATE accounts SET last_email = $1 WHERE id = $2`
		if _, err := db.Exec(msg, time.Now().Unix(), account.ID); err != nil {
			continue
		}

		accounts = append(accounts, account)
	}

	if len(accounts) == 0 {
		fmt.Println("No accounts to email today")
	}

	return accounts
}

func createBooksTable() {
	msg := `CREATE TABLE IF NOT EXISTS books (
	  id SERIAL PRIMARY KEY,
	  title TEXT NOT NULL ,
	  author TEXT,
	  available BOOLEAN DEFAULT true,
	  ts BIGINT,
	  customer TEXT,
	  customer_id INT
	);`

	if _, err := db.Exec(msg); err != nil {
		fmt.Println("error creating table", err)
	}
}

func createAccountsTable() {
	msg := `CREATE TABLE IF NOT EXISTS accounts (
	  id SERIAL PRIMARY KEY,
	  name TEXT NOT NULL ,
	  email TEXT NOT NULL ,
	  email_list BOOLEAN DEFAULT true,
	  last_email BIGINT
	);`

	if _, err := db.Exec(msg); err != nil {
		fmt.Println("error creating table", err)
	}
}
