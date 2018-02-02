package config

import (
	"database/sql"

	"golang.org/x/oauth2"
	//Blank import to get the postgress driver
	_ "github.com/lib/pq"
)

type Env struct {
	DB          *sql.DB
	OauthConfig *oauth2.Config
}

func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
