// main_test.go

package main

import (
  "fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var ws WebServivce

func TestMain(m *testing.M) {
	ws = WebServivce{}
	ws.Initialize("root", "Cameron31!", "MalwareURLs", "allowAllFiles=true")
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS URLS (
  id int(6) unsigned NOT NULL AUTO_INCREMENT,
  url varchar(100) NOT NULL,
  malware bool NOT NULL,
  PRIMARY KEY (id)
) `

func ensureTableExists() {
	if _, err := ws.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	ws.DB.Exec("DELETE FROM URLS")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	reqrec := httptest.NewRecorder()
	ws.Router.ServeHTTP(reqrec, req)
	return reqrec
}

func actualEqualExpected(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected: %d. Actual: %d\n", expected, actual)
	}
}

func TestEmptyDB(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/urlinfo/1/test.com", nil)
	response := executeRequest(req)
	actualEqualExpected(t, http.StatusNotFound, response.Code)
}

func addUrlNoMalware() {
  statement := fmt.Sprintf("INSERT INTO URLS(url, malware) VALUES('test.com', false)")
	ws.DB.Exec(statement)
}

func TestGetUrl(t *testing.T) {
	addUrlNoMalware()
	req, _ := http.NewRequest("GET", "/urlinfo/1/test.com", nil)
	response := executeRequest(req)
	actualEqualExpected(t, http.StatusOK, response.Code)
}

func addUrlMalware() {
  statement := fmt.Sprintf("INSERT INTO URLS(url, malware) VALUES('test2.com', true)")
	ws.DB.Exec(statement)
}

func TestGetUrlMalware(t *testing.T) {
	addUrlMalware()
	req, _ := http.NewRequest("GET", "/urlinfo/1/test2.com", nil)
	response := executeRequest(req)
	actualEqualExpected(t, http.StatusOK, response.Code)
}
