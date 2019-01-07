// app.go

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "strconv"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type WebServivce struct {
	Router *mux.Router
	DB     *sql.DB
}

func (ws *WebServivce) Initialize(user, password, dbname string, other string) {
	connectionString := fmt.Sprintf("%s:%s@(http://127.0.0.1:8080)/%s?%s", user, password, dbname, other)

	var err error
	ws.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
  _,err = ws.DB.Exec("DROP DATABASE MalwareURLs")
   if err != nil {
       log.Fatal(err)
   }
  _,err = ws.DB.Exec("CREATE DATABASE MalwareURLs")
   if err != nil {
       log.Fatal(err)
   }
   _,err = ws.DB.Exec("USE MalwareURLs")
  if err != nil {
      log.Fatal(err)
  }
  const createTables = `
    CREATE TABLE IF NOT EXISTS URLS (
      url varchar(100) NOT NULL,
      malware bool NOT NULL,
      PRIMARY KEY (url));`
  _,err = ws.DB.Exec(createTables)
   if err != nil {
      log.Fatal(err)
   }
   _,err = ws.DB.Exec("INSERT INTO URLS(url, malware) VALUES('www.test.com', false);")
  if err != nil {
      log.Fatal(err)
  }
  _,err = ws.DB.Exec("INSERT INTO URLS(url, malware) VALUES('www.nhl.com', true);")
  if err != nil {
     log.Fatal(err)
   }
  _,err = ws.DB.Exec("INSERT INTO URLS(url, malware) VALUES('www.google.com', true);")
  if err != nil {
    log.Fatal(err)
  }
    // filePath := "/desktop/data.txt"
    // _,err = ws.DB.Exec("LOAD DATA LOCAL INFILE '" + filePath + "' INTO TABLE URLS")
    // if err != nil {
    //    log.Fatal(err)
    // }
	ws.Router = mux.NewRouter()
	ws.initializeRoutes()
}

func (ws *WebServivce) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, ws.Router))
}

func (ws *WebServivce) initializeRoutes() {
	ws.Router.HandleFunc("/urlinfo/1/{url}", ws.getUrl).Methods("GET")
}

func (ws *WebServivce) getUrl(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
	url2 := vars["url"]
  url := strconv.Quote(url2)
  // var err error
	// if err != nil {
	// 	respondWithError(w, http.StatusBadRequest, "Invalid URL")
	// 	return
	// }
	u := entry{Url: url}
	if err := u.getUrl(ws.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "URL not found in databases - proceed at caution")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
