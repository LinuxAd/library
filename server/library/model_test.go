package library

import (
	"database/sql"
	"strings"
	"testing"
)

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
	b := Book{}
	m := MockModelDB{Book: b}

	got := b.addBook(&m)
	if got != nil {
		t.Errorf("got error, expected none: %v", got)
	}
	if len(m.callParams) == 0 {
		t.Errorf("no call params sent to DB server, got: %v", len(m.callParams))
	}
}
