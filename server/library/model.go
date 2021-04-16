package library

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Book struct {
	ID          int
	Author      string
	Title       string
	Description string
	ISBN        string
}

type Error struct {
	Msg  string `json:"msg"`
	Body string `json:",omitempty"`
}

type Response struct {
	Books     []Book    `json:"books,omitempty"`
	Total     int       `json:"total_count"`
	RCreation time.Time `json:"request_time"`
	Status    string    `json:"status"`
	Error     `json:"error,omitempty"`
}

type SQLDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Ping() error
}

type SQLResult interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

func writeResponse(res Response, w *http.ResponseWriter) {
	bytes, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		http.Error(*w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(*w, "%s\n", bytes)
}

func (b *Book) checkDB(db SQLDB) error {
	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}

func (b *Book) addBook(db SQLDB) error {
	err := db.QueryRow(
		"INSERT INTO books(author, title, description, isbn) VALUES($1, $2, $3, current_timestamp, current_timestamp) RETURNING id",
		b.Author, b.Title, b.Description, b.ISBN).Scan(&b.ID)
	if err != nil {
		return err
	}
	return nil
}
