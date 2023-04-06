package sqlpkg

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"forum/model"

	"github.com/mattn/go-sqlite3"
)

type ForumModel struct {
	DB *sql.DB
}

func OpenDB(name string) (*sql.DB, error) {
	// init pull (not connection)
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_auth&_auth_user=webuser&_auth_pass=webuser", name))
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

// TODO change name of categories
func CreateDB(name, admName, admEmail, admPass string) (*sql.DB, error) {
	// init pull (not connection)
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_auth&_auth_user=admin&_auth_pass=adminpass", name))
	if err != nil {
		return nil, err
	}

	// create a not-admin user
	var sqlconn *sqlite3.SQLiteConn
	err = sqlconn.AuthUserAdd("webuser", "webuser", false)
	if err != nil {
		return nil, err
	}

	q := `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			name TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			dateCreate TIMESTAMP NOT NULL,
			session TEXT,
			expirySession TIMESTAMP
		);
		
		CREATE TABLE 'posts_likes' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			userID INT NOT NULL REFERENCES users(id),
			messageID INT NOT NULL REFERENCES posts(id),
			like BOOL NOT NULL,
			UNIQUE (userID, messageID)
		);
		
		CREATE TABLE 'comments_likes' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			userID INT NOT NULL REFERENCES users(id),
			messageID INT NOT NULL REFERENCES comments(id),
			like BOOL NOT NULL,
			UNIQUE (userID, messageID)
		);

		CREATE TABLE 'posts' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			theme TEXT NOT NULL DEFAULT ('(No theme)'),
			content TEXT NOT NULL, 
			authorID INT NOT NULL REFERENCES users(id),
			dateCreate TIMESTAMP NOT NULL
		);

		CREATE TABLE 'comments' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			content TEXT NOT NULL, 
			authorID INT NOT NULL REFERENCES users(id),
			dateCreate TIMESTAMP NOT NULL,
			postID INT NOT NULL REFERENCES posts(id)
		);
		
		CREATE TABLE 'categories' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			name TEXT NOT NULL 
		);
		
		CREATE TABLE 'post_categories' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			categoryID INT NOT NULL REFERENCES categories(id), 
			postID INT NOT NULL REFERENCES posts(id),
			UNIQUE (categoryID, postID)
		);

		CREATE INDEX userssession ON users (session);

		INSERT INTO users (name,email,password, dateCreate) VALUES (?,?,?,?);
		INSERT INTO categories (name) VALUES (?), (?), (?);
		
		--!!!!! tests
		INSERT INTO users (name,email,password, dateCreate) VALUES ("test1","test1@forum","tst1","2023-03-20 09:41:04.656479916+00:00");
		INSERT INTO users (name,email,password, dateCreate) VALUES ("test2","test2@forum","tst1","2023-03-20 09:52:04.656479916+00:00");
		INSERT INTO posts (theme,content,authorID, dateCreate) VALUES ("cats", "cats are cute", 1, "2023-03-20 15:41:04.656479916+00:00");
		INSERT INTO posts (theme,content,authorID, dateCreate) VALUES ("dogs", "dogs are funny", 2, "2023-03-21 14:41:04.656479916+00:00");
		INSERT INTO posts (theme,content,authorID, dateCreate) VALUES ("My cat", "She is the best", 3, "2023-03-22 10:41:04.656479916+00:00");
		INSERT INTO posts (theme,content,authorID, dateCreate) VALUES ("My dog"," He is the best", 2, "2023-03-20 11:41:04.656479916+00:00");
		INSERT INTO posts (theme,content,authorID, dateCreate) VALUES ("My parrot", "My parrot is a pirate", 1, "2023-03-22 11:41:04.656479916+00:00");
		
		INSERT INTO comments (content,authorID, dateCreate,postID) VALUES ("No, mine", 1, "2023-03-22 11:41:04.656479916+00:00",3);
		
		INSERT INTO posts_likes (userID, messageID, like) VALUES (2, 4, 0);
		INSERT INTO posts_likes (userID, messageID, like) VALUES (2, 3, 0);
		INSERT INTO posts_likes (userID, messageID, like) VALUES (2, 2, 1);
		INSERT INTO posts_likes (userID, messageID, like) VALUES (1, 2, 1);
		INSERT INTO posts_likes (userID, messageID, like) VALUES (1, 3, 1);
		INSERT INTO posts_likes (userID, messageID, like) VALUES (1, 1, 1);
		
		INSERT INTO comments_likes (userID, messageID, like) VALUES (1, 1, 1);
		INSERT INTO comments_likes (userID, messageID, like) VALUES (3, 1, 0);
		
		INSERT INTO post_categories (categoryID, postID) VALUES (1,1);
		INSERT INTO post_categories (categoryID, postID) VALUES (3,1);
		INSERT INTO post_categories (categoryID, postID) VALUES (2,2);
		INSERT INTO post_categories (categoryID, postID) VALUES (3,2);
		INSERT INTO post_categories (categoryID, postID) VALUES (1,3);
		INSERT INTO post_categories (categoryID, postID) VALUES (3,3);
		INSERT INTO post_categories (categoryID, postID) VALUES (2,4);
		INSERT INTO post_categories (categoryID, postID) VALUES (3,4);
		INSERT INTO post_categories (categoryID, postID) VALUES (3,5);
		
	`
	// use a  transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, handleErrAndCloseDB(db, "transaction begin", err) // close DB and return error
	}

	// try exec transaction
	_, errExec := tx.Exec(q, admName, admEmail, admPass, time.Now(), "cats", "dogs", "pets")
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

	// open the DB with no admin user
	db, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?_auth&_auth_user=webuser&_auth_pass=webuser", name))
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

/*
checks if the value exists in the table's field and returns the number of rows where the value was found
*/
func (f *ForumModel) checkExisting(table, field, value string) (int, error) {
	q := `SELECT ` + field + ` FROM ` + table + ` WHERE ` + field + ` = ?`
	res, err := f.DB.Exec(q, value)
	if err != nil {
		return 0, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

/*
check the res and returns error=nil if only 1 row had been affected, 
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