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
	db   *sql.DB
)

type Book struct {
	ID int32 `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Available bool `json:"available"`
	Customer string `json:"customer"`
	TimeStamp int64 `json:"time_stamp"`
}

type BookList struct {
	IDs []int32 `json:"ids"`
	Name string `json:"name"`
	Available bool `json:"available"`
	TimeStamp int64 `json:"time_stamp"`
}

type Search struct {
	Books []Book `json:"books"`
	Query string `json:"query"`
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
	var ts int64

	msg := `UPDATE books SET available = $1, customer=$2, ts = $3 WHERE id = $4`

	if !b.Available {
		ts = time.Now().Unix()
	}

	if _, err := db.Exec(msg, b.Available, b.Customer, ts, b.ID); err != nil {
		return err
	}
	return nil
}

func (b *Book) GetBookByID() error  {
	msg := `SELECT title, available, author FROM books WHERE id = $1`

	if err := db.QueryRow(msg, b.ID).Scan(&b.Title, &b.Available, &b.Author); err != nil {
		if !strings.Contains(err.Error(), "converting NULL to string") {
			return err
		}
	}

	if !b.Available {
		msg = `SELECT customer, ts FROM books WHERE id = $1`
		if err := db.QueryRow(msg, b.ID).Scan(&b.Customer, &b.TimeStamp); err != nil {
			return err
		}
	}

	return nil
}

func (b *BookList) UpdateListBooks() {
	for _, id := range b.IDs {
		msg := `UPDATE books SET available = $1, customer = $2, ts = $3 WHERE id = $4`
		if _, err := db.Exec(msg, b.Available, b.Name, time.Now().Unix(), id); err != nil {
			fmt.Println("error updating list book for id", id, err)
		}
	}
}

func (s *Search) SearchByTitle() error  {
	msg := `SELECT title, id, available, author FROM books WHERE title ILIKE $1 order by title`

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
			msg = `SELECT customer FROM books WHERE id = $1`
			if err := db.QueryRow(msg, book.ID).Scan(&book.Customer); err != nil {
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
			msg = `SELECT customer, ts FROM books WHERE id = $1`
			if err := db.QueryRow(msg, book.ID).Scan(&book.Customer, &book.TimeStamp); err != nil {
				if !strings.Contains(err.Error(), "converting NULL to string is unsupported") {
					fmt.Println("error getting customer for book ID", book.ID, err)
				}
			}
		}
		b = append(b, book)
	}

	return b, nil
}

func createTable()  {
	msg := `CREATE TABLE books (
	  id SERIAL PRIMARY KEY,
	  title TEXT NOT NULL ,
	  author TEXT,
	  available BOOLEAN DEFAULT true,
	  ts BIGINT,
	  customer TEXT
	);`


	if _, err := db.Exec(msg); err != nil {
		fmt.Println("error creating table", err)
	}
}