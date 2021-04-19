package library

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	Begin() (*sql.Tx, error)
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
	tx, err := db.Begin()
	if err != nil {
		return nil
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if err := tx.QueryRow(
		"INSERT INTO books(author, title, description, isbn) VALUES($1, $2, $3, $4) RETURNING id",
		b.Author, b.Title, b.Description, b.ISBN).Scan(&b.ID); err != nil {
		return err
	}
	return nil
}

func (b *Book) removeBook(db SQLDB) error {
	if b.ID == 0 {
		return errors.New("cannot delete book with ID of 0")
	}

	tx, err := db.Begin()
	if err != nil {
		return nil
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(
		"DELETE FROM books WHERE id=$1", b.ID)
	if err != nil {
		return err
	}
	return nil
}
