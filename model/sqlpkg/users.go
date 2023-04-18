package sqlpkg

import (
	"database/sql"
	"errors"
	"time"

	"forum/model"
)

/*
returns list of all users in DB
*/
func (f *ForumModel) GetUsers() ([]*model.User, error) {
	q := `SELECT id, name, dateCreate FROM users`

	rows, err := f.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.DateCreate)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

/*
returns a user from DB by ID
*/
func (f *ForumModel) GetUserByID(id int) (*model.User, error) {
	q := `SELECT id, name, email, password, dateCreate, session, expirySession FROM users WHERE id=?`

	user := &model.User{}
	var ses sql.NullString
	row := f.DB.QueryRow(q, id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.DateCreate, &ses, &user.ExpirySession)
	user.Session = ses.String
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoRecord
		}
		return nil, err
	}
	return user, nil
}

/*
returns a user from DB by the name
*/
func (f *ForumModel) GetUserByName(name string) (*model.User, error) {
	q := `SELECT id, name, email, password, dateCreate, session, expirySession FROM users WHERE name=?`

	user := &model.User{}
	var ses sql.NullString
	row := f.DB.QueryRow(q, name)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.DateCreate, &ses, &user.ExpirySession)
	user.Session = ses.String
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoRecord
		}
		return nil, err
	}

	return user, nil
}

/*
returns a user from DB by the email
*/
func (f *ForumModel) GetUserByEmail(email string) (*model.User, error) {
	q := `SELECT id, name, email, password, dateCreate, session, expirySession FROM users WHERE name=?`

	user := &model.User{}
	var ses sql.NullString
	row := f.DB.QueryRow(q, email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.DateCreate, &ses, &user.ExpirySession)
	user.Session = ses.String
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoRecord
		}
		return nil, err
	}

	return user, nil
}

/*
returns a user from DB by the email
*/
func (f *ForumModel) GetUserBySession(session string) (*model.User, error) {
	q := `SELECT id, name, email, dateCreate, session, expirySession FROM users WHERE session=?`

	user := &model.User{}
	var ses sql.NullString
	row := f.DB.QueryRow(q, session)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.DateCreate, &ses, &user.ExpirySession)
	user.Session = ses.String
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoRecord
		}
		return nil, err
	}

	return user, nil
}

/*
inserts the new user into DB. It doesn't do any check of unique data. But if DB have some restricts, it will return an error
*/
func (f *ForumModel) InsertUser(name, email string, password []byte, dateCreate time.Time) (int, error) {
	q := `INSERT INTO users (name, email, password, dateCreate) VALUES (?,?,?,?)`
	res, err := f.DB.Exec(q, name, email, password, dateCreate)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

/*
adds a session to the user with the given ID
*/
func (f *ForumModel) AddUsersSession(id int, session string, expirySession time.Time) error {
	q := `UPDATE users SET session=?, expirySession=? WHERE id=?`
	res, err := f.DB.Exec(q, session, expirySession, id)
	if err != nil {
		return err
	}

	return f.checkUnique(res)
}

/*
deletes the user's session
*/
func (f *ForumModel) DeleteUsersSession(id int) error {
	q := `UPDATE users SET session=NULL, expirySession=NULL WHERE id=?`
	res, err := f.DB.Exec(q, id)
	if err != nil {
		return err
	}

	return f.checkUnique(res)
}

/*
check if a user with the given name exists,  returns nil only if there is exactly one user
*/
func (f *ForumModel) CheckUserByName(name string) error {
	err := f.checkExisting("users", "name", name)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNoRecord
	}
	return err
}

/*
check if a user with the given email exists, returns nil only if there is exactly one user
*/
func (f *ForumModel) CheckUserByEmail(email string) error {
	err := f.checkExisting("users", "email", email)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNoRecord
	}
	return err
}

/*
adds a session to the user with the given ID
*/
func (f *ForumModel) AddUser(name, email string, password []byte, dateCreate time.Time) (int, error) {
	id, err := f.InsertUser(name, email, password, dateCreate)
	if err != nil {
		errUnique := f.CheckUserByName(name)
		if errUnique == nil {
			return 0, model.ErrUniqueUserName
		}
		errUnique = f.CheckUserByEmail(email)
		if errUnique == nil {
			return 0, model.ErrUniqueUserEmail
		}
	}

	return id, nil
}
