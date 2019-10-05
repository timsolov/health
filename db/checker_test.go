package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewMySQLChecker(t *testing.T) {
	db, _, err := sqlmock.New()

	if err != nil {
		t.Errorf("sqlmock.New().error != nil, wants nil")
	}

	defer db.Close()

	c := NewMySQLChecker(db)
	expected := "SELECT 1"

	if c.CheckSQL != expected {
		t.Errorf("c.CheckSQL == %s, wants %s", c.CheckSQL, expected)
	}

	expected = "SELECT VERSION()"
	if c.VersionSQL != expected {
		t.Errorf("c.VersionSQL == %s, wants %s", c.VersionSQL, expected)
	}
}

func TestNewPostgreSQLCheckerChecker(t *testing.T) {
	db, _, err := sqlmock.New()

	if err != nil {
		t.Errorf("sqlmock.New().error != nil, wants nil")
	}

	defer db.Close()

	c := NewPostgreSQLChecker(db)
	expected := "SELECT 1"

	if c.CheckSQL != expected {
		t.Errorf("c.CheckSQL == %s, wants %s", c.CheckSQL, expected)
	}

	expected = "SELECT VERSION()"
	if c.VersionSQL != expected {
		t.Errorf("c.VersionSQL == %s, wants %s", c.VersionSQL, expected)
	}
}

func TestNewSqlite3CheckerChecker(t *testing.T) {
	db, _, err := sqlmock.New()

	if err != nil {
		t.Errorf("sqlmock.New().error != nil, wants nil")
	}

	defer db.Close()

	c := NewSqlite3Checker(db)
	expected := "SELECT 1"

	if c.CheckSQL != expected {
		t.Errorf("c.CheckSQL == %s, wants %s", c.CheckSQL, expected)
	}

	expected = "SELECT sqlite_version()"
	if c.VersionSQL != expected {
		t.Errorf("c.VersionSQL == %s, wants %s", c.VersionSQL, expected)
	}
}

func TestCheck_up(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("sqlmock.New().error != nil, wants nil")
	}

	defer db.Close()

	checker := NewChecker("SELECT 1", "SELECT VERSION()", db)

	rows := sqlmock.NewRows([]string{"1"}).AddRow("1")
	mock.ExpectQuery(checker.CheckSQL).WillReturnRows(rows)

	dummyVersion := "10.1.9-DummyDB"
	rows = sqlmock.NewRows([]string{"VERSION()"}).AddRow(dummyVersion)
	mock.ExpectQuery(checker.VersionSQL).WillReturnRows(rows)

	health := checker.Check()

	if health.IsDown() {
		t.Errorf("health.IsDown() == %t, wants %t", health.IsDown(), false)
	}

	version := health.GetInfo("version")

	if version != dummyVersion {
		t.Errorf("version == %s, wants %s", version, dummyVersion)
	}
}

func TestCheck_up_not_set_version(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("sqlmock.New().error != nil, wants nil")
	}

	defer db.Close()

	checker := NewChecker("SELECT 1", "", db)

	rows := sqlmock.NewRows([]string{"1"}).AddRow("1")
	mock.ExpectQuery(checker.CheckSQL).WillReturnRows(rows)

	health := checker.Check()

	if health.IsDown() {
		t.Errorf("health.IsDown() == %t, wants %t", health.IsDown(), false)
	}

	// must not be present in the response
	version := health.GetInfo("version")

	if version != nil {
		t.Errorf("version == %s, wants %v", version, nil)
	}
}

func TestCheck_down(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("sqlmock.New().error != nil, wants nil")
	}

	defer db.Close()

	checker := NewChecker("SELECT 1", "SELECT VERSION()", db)

	expectedError := "Expected error"
	mock.ExpectQuery(checker.CheckSQL).WillReturnError(errors.New(expectedError))

	health := checker.Check()

	if health.IsUp() {
		t.Errorf("health.IsUp() == %t, wants %t", health.IsUp(), false)
	}

	message := health.GetInfo("error")

	if message != expectedError {
		t.Errorf("message == %s, wants %s", message, expectedError)
	}
}

func TestCheck_down_version(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("sqlmock.New().error != nil, wants nil")
	}

	defer db.Close()

	checker := NewChecker("SELECT 1", "SELECT VERSION()", db)

	rows := sqlmock.NewRows([]string{"1"}).AddRow("1")
	mock.ExpectQuery(checker.CheckSQL).WillReturnRows(rows)

	expectedError := "Expected error"
	mock.ExpectQuery(checker.VersionSQL).WillReturnError(errors.New(expectedError))

	health := checker.Check()

	if health.IsUp() {
		t.Errorf("health.IsUp() == %t, wants %t", health.IsUp(), false)
	}

	message := health.GetInfo("error")

	if message != expectedError {
		t.Errorf("message == %s, wants %s", message, expectedError)
	}
}

func TestCheck_down_nil_db(t *testing.T) {
	checker := NewChecker("SELECT 1", "SELECT VERSION()", nil)

	expectedError := "Empty resource"

	health := checker.Check()

	if health.IsUp() {
		t.Errorf("health.IsUp() == %t, wants %t", health.IsUp(), false)
	}

	message := health.GetInfo("error")

	if message != expectedError {
		t.Errorf("message == %s, wants %s", message, expectedError)
	}
}
