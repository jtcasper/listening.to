package orm

import (
	"database/sql"
	"fmt"
	"listening.to/types"
	"log"
)

type Orm struct {
	db *sql.DB
}

func New(s string) (o *Orm, err error) {
	var db *sql.DB
	switch s {
	case "sqlite3":
		db, err = sql.Open("sqlite3", "listening.db")
	default:
		return nil, fmt.Errorf("No orm for driver %s ", s)
	}
	if err != nil {
		return nil, err
	}
	o = &Orm{db: db}
	return
}

func (o *Orm) Destroy() {
	o.db.Close()
}

func (o *Orm) Write(v interface{}) (err error) {
	switch t := v.(type) {
	case types.Account:
		log.Printf("%+v\n", o)
		_, err = o.db.Exec("INSERT OR REPLACE INTO ACCOUNTS (ID, ACCESS_TOKEN, REFRESH_TOKEN) VALUES ($1, $2, $3) ",
			t.ID,
			t.AccessToken,
			t.RefreshToken,
		)
	default:
		return fmt.Errorf("Not implemented for type %T\n", t)
	}
	return err

}
