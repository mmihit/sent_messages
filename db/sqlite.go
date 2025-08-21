package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DataBase struct {
	DB *sql.DB
}

func InitDb() (*DataBase, error) {
	db, err := sql.Open("sqlite3", "./db/my_app.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	for _, column := range Tables {
		_, err := db.Exec(column)
		if err != nil {
			return nil, err
		}
	}

	return &DataBase{
		DB: db,
	}, nil
}
