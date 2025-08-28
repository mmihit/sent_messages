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

	err = InsertAdminAccount(db)
	if err != nil {
		return nil, err
	}

	return &DataBase{
		DB: db,
	}, nil
}

func InsertAdminAccount(Db *sql.DB) error {
	_, err := Db.Exec(`INSERT OR IGNORE INTO cliniques VALUES ('1','mohammed mihit', 'MED', 'med86004@gmail.com','', 'Medmohammed310@20',1,'','');`)
	if err != nil {
		return err
	}
	return nil
}
