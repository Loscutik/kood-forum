package sqlpkg

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"forum/model"

	_ "github.com/mattn/go-sqlite3"
)

type ForumModel struct {
	DB *sql.DB
}

func OpenDB(name, user, pass string) (*sql.DB, error) {
	// init pull (not connection)
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_auth&_auth_user=%s&_auth_pass=%s&_foreign_keys=on", name, user, pass))
	if err != nil {
		return nil, err
	}

	// check connection (create and check)
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func handleErrAndCloseDB(db *sql.DB, operation string, err error) error {
	errClose := db.Close()
	if errClose != nil {
		return fmt.Errorf("%s failed: %v, unable to close DB: %v", operation, err, errClose)
	}
	return fmt.Errorf("DB was closed cause %s failed: %v", operation, err)
}

func CreateDB(name, admName, admEmail, admPass string) (*sql.DB, error) {
	// init pull (not connection)
	db, err := OpenDB(name, "admin", "adminpass")
	if err != nil {
		return nil, err
	}

	createQuery, err := os.ReadFile("model/sqlpkg/creatingDB.sql")
	if err != nil {
		return nil, fmt.Errorf("reading creatingDB.sql faild: %v", err)
	}

	db.Exec("PRAGMA foreign_keys = ON;")
	// use a  transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, handleErrAndCloseDB(db, "transaction begin", err) // close DB and return error
	}

	// try exec transaction
	_, errExec := tx.Exec(string(createQuery), admName, admEmail, admPass, time.Now(), "cats", "dogs", "pets", "savage")
	if errExec != nil {
		errRoll := tx.Rollback()
		if errRoll != nil {
			return nil, fmt.Errorf("table creating failed: %v, unable to rollback: %v", errExec, errRoll)
		}
		return nil, handleErrAndCloseDB(db, "table creating", errExec)
	}

	// if the transaction was a success
	err = tx.Commit()
	if err != nil {
		return nil, handleErrAndCloseDB(db, "transaction commit", err)
	}

	err = db.Close()
	if err != nil {
		return nil, err
	}

	// open the DB with no admin user and check the connection
	db, err = OpenDB(name, "webuser", "webuser")
	if err != nil {
		return nil, err
	}

	return db, nil
}

/*
fills in the DB with data from the given file
*/
func (f *ForumModel) FillInDB(fileName string, params ...any) error {
	query, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("reading creatingDB.sql faild: %v", err)
	}
	// use a  transaction
	tx, err := f.DB.Begin()
	if err != nil {
		return fmt.Errorf("table filling failed: transaction begin faild: %v", err) // close DB and return error
	}

	// try exec transaction
	_, errExec := tx.Exec(string(query),params...)
	if errExec != nil {
		errRoll := tx.Rollback()
		if errRoll != nil {
			return fmt.Errorf("table filling failed: %v, unable to rollback the transaction: %v", errExec, errRoll)
		}
		return fmt.Errorf("table filling failed: transaction execute faild: %v", errExec)
	}

	// if the transaction was a success
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("table filling failed: transaction commit faild: %v", err)
	}

	return nil
}

/*
checks if the value exists in the table's field and returns the number of rows where the value was found
*/
func (f *ForumModel) checkExisting(table, field, value string) error {
	q := `SELECT ` + field + ` FROM ` + table + ` WHERE ` + field + ` = ?`
	row := f.DB.QueryRow(q, value)
	var tmp string
	return row.Scan(&tmp)
}

/*
checks the res and returns error=nil if only 1 row had been affected,
in the other cases returns  ErrNoRecord (for 0 rows), or ErrTooManyRecords (for more than 1)
*/
func (f *ForumModel) checkUnique(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 1 {
		return nil
	}
	if n == 0 {
		return model.ErrNoRecord
	}
	if n > 1 {
		return model.ErrTooManyRecords
	}
	return errors.New("negative number of rows")
}
