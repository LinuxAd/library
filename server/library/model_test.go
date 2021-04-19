package library

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// The MockModelDB object has a few parameters that the normal struct doesn't have - this allows us to spy on the calls made
// to the fake db
type MockModelDB struct {
	Book
	callParams []interface{}
	actions    []string
	calls      int
}

func getAction(query string) string {
	queryWords := strings.Split(query, " ")
	return queryWords[0]
}

func (m *MockModelDB) Begin() (tx *sql.Tx, err error) {
	return
}

func (m *MockModelDB) record(query string, args ...interface{}) {
	m.callParams = args
	m.actions = append(m.actions, getAction(query))
}

func (m *MockModelDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	m.record(query, args)
	m.calls++
	return nil, nil
}

func (m *MockModelDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	m.record(query, args)
	m.calls++
	return nil, nil
}

func (m *MockModelDB) QueryRow(query string, args ...interface{}) *sql.Row {
	m.record(query, args)
	m.calls++
	return nil
}

func (m *MockModelDB) Ping() error {
	m.record("", "")
	m.calls++
	return nil
}

func Test_checkDB(t *testing.T) {
	b := Book{}
	m := MockModelDB{Book: b}
	got := b.checkDB(&m)
	if got != nil {
		t.Errorf("got error should be nil: %v", got)
	}
	if m.calls > 1 {
		t.Errorf("DB called more than once")
	}
}

func Test_addBook(t *testing.T) {
	// create the book fields to create book objects for testing
	type fields struct {
		Title       string
		Author      string
		Description string
		ISBN        string
	}
	// create a list of data fields for testing
	cases := []struct {
		name      string
		fields    fields
		shouldErr bool
	}{
		{
			name: "quick test",
			fields: fields{
				Title:       "Mort",
				Author:      "Terry Pratchett",
				Description: "Mort is a fantasy novel by British writer Terry Pratchett.",
				ISBN:        "9780552144292",
			},
			shouldErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// create a new mock db
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was experienced when opening a database connection", err)
			}
			// close it at end of every test case
			defer db.Close()

			// create the expected row data
			rows := sqlmock.NewRows([]string{"id"}).
				AddRow(1)

			// tell mock db to expect a tx begin
			mock.ExpectBegin()
			// tell mockdb to expect this kind of query
			mock.ExpectQuery("INSERT INTO books").WithArgs(
				tt.fields.Author,
				tt.fields.Title,
				tt.fields.Description,
				tt.fields.ISBN).WillReturnRows(rows)
			// expect a tx commit
			mock.ExpectCommit()

			// create the book
			b := Book{
				Author:      tt.fields.Author,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				ISBN:        tt.fields.ISBN,
			}

			// call the method
			err = b.addBook(db)

			// check method return for errors
			if tt.shouldErr != (err != nil) {
				t.Errorf("error was not expected, got: %s", err)
			}
			// check if db expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %s", err)
			}

		})
	}

}
