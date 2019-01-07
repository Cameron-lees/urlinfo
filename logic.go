// model.go

package main

import (
  "fmt"
	"database/sql"
)

type entry struct {
	Url      string  `json:"-"`
	Malware  bool    `json:"malware"`
}

func (u *entry) getUrl(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT malware FROM URLS WHERE url=%s", u.Url)
	return db.QueryRow(statement).Scan(&u.Malware)
}
