package orm

import (
	"database/sql"
	"fmt"
	"listening.to/types"
)

type Orm struct {
	db *sql.DB
}

type Queryable interface {
	Table() string
}

type Rows struct {
	rows *sql.Rows
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
		_, err = o.db.Exec("INSERT OR REPLACE INTO "+t.Table()+" (ID, ACCESS_TOKEN, REFRESH_TOKEN, EXPIRY) VALUES ($1, $2, $3, $4) ",
			t.ID,
			t.Token.AccessToken,
			t.Token.RefreshToken,
			t.Token.Expiry,
		)
	default:
		return fmt.Errorf("Not implemented for type %T\n", t)
	}
	return err
}

func (o *Orm) Query(q Queryable) (*Rows, error) {
	rows, err := o.db.Query("SELECT * FROM " + q.Table())
	if err != nil {
		return nil, err
	}
	r := &Rows{rows}
	return r, nil
}

func (r *Rows) GetAccounts() []types.Account {
	var accs []types.Account
	defer r.rows.Close()
	for r.rows.Next() {
		var acc types.Account
		r.rows.Scan(&acc.ID, &acc.Token.AccessToken, &acc.Token.RefreshToken, &acc.Token.Expiry)
		accs = append(accs, acc)
	}
	return accs
}
