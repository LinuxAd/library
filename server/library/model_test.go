package library

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_checkDB(t *testing.T) {
	type fields struct {
		Title       string
		Author      string
		Description string
		ISBN        string
	}
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
			db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Errorf("got error '%s' creating mock db", err)
			}
			mock.ExpectPing()
			b := Book{
				Author:      tt.fields.Author,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				ISBN:        tt.fields.ISBN,
			}

			b.checkDB(db)
			if tt.shouldErr != (err != nil) {
				t.Errorf("error status did not meet expectations, got: '%s'", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("database situation did not meet expectations: '%s'", err)
			}
		})
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

func Test_removeBook(t *testing.T) {
	// create the book fields to create book objects for testing
	type fields struct {
		Title       string
		Author      string
		Description string
		ISBN        string
		ID          int
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
				ID:          1,
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

			// create the expected result data
			res := sqlmock.NewResult(1, 1)

			// tell mock db to expect a tx begin
			mock.ExpectBegin()
			// tell mockdb to expect this kind of query
			mock.ExpectExec("DELETE FROM books").WithArgs(
				tt.fields.ID).WillReturnResult(res)
			// expect a tx commit
			mock.ExpectCommit()

			// create the book
			b := Book{
				Author:      tt.fields.Author,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				ISBN:        tt.fields.ISBN,
				ID:          tt.fields.ID,
			}

			// call the method
			err = b.removeBook(db)

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
