package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	ADM_NAME  = "admin"
	ADM_EMAIL = "admin@forum.com"
	ADM_PASS  = "admin"
)

type Time time.Time

// Tables' names for likes
const (
	POSTS_LIKES    = "posts_likes"
	COMMENTS_LIKES = "comments_likes"
)
const N_LIKES = 2

const (
	LIKE = iota
	DISLIKE
)

var (
	ErrNoRecord        = errors.New("there is no record in the DB")
	ErrTooManyRecords  = errors.New("there are more than one record")
	ErrUnique          = errors.New("unique constraint failed")
	ErrUniqueUserName  = errors.New("user with the given name already exists")
	ErrUniqueUserEmail = errors.New("user with the given email already exists")
)

type User struct {
	ID            int
	Name          string
	Password      []byte
	Email         string
	DateCreate    time.Time
	Session       string
	ExpirySession sql.NullTime
}

type message struct {
	Author     *User
	Content    string
	DateCreate time.Time
	Likes      []int // index 0 keeps number of likes, index 1 keeps number of dislikes
}

type Post struct {
	ID         int
	Theme      string
	Message    message
	Categories []*Category
	Comments   []*Comment
}

type Category struct {
	ID   int
	Name string
}

type Comment struct {
	ID      int
	Message message
}

type Filter struct {
	CategoryID []int
	AuthorID, LikedByUserID int
}

/*
implements the Scanner interface.
*/
func (t *Time) Scan(src any) error {
	// tt:=time.Time(*t)
	const DateTimeSQLite = "2006-01-02 15:04:05.999999999-07:00"
	if src == nil {
		*t = Time(time.Time{})
		return nil
	}
	timeStr, ok := src.(string)
	if !ok {
		return errors.New("parametr is not a string")
	}
	tt, err := time.Parse(DateTimeSQLite, timeStr)
	if err != nil {
		return err
	}
	*t = (Time(tt))
	return nil
}

func (t *Time) String() string {
	return time.Time(*t).String()
}

func (u *User) String() string {
	if u == nil {
		return "nil"
	}
	return fmt.Sprintf("user id: %d --  name: %s --  email: %s --  password: %s --  DataCreate: %s --  session: %s\n",
		u.ID, u.Name, u.Email, u.Password, u.DateCreate.String(), u.Session)
}

func (p *Post) String() string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("id: %d -- Theme: %s\nMessage: \n%s\nCategories: \n%v\nComments: \n%v\n",
		p.ID, p.Theme, p.Message.String(), p.Categories, p.Comments)
}

func (m *message) String() string {
	if m == nil {
		return "nil"
	}
	return fmt.Sprintf("  Author: (%p)\n  %v\n  Content: %s\n  DataCreate: %s -- Likes: %#v\n",
		m.Author, m.Author, m.Content, m.DateCreate.String(), m.Likes)
}

func (c *Category) String() string {
	if c == nil {
		return "nil"
	}
	return fmt.Sprintf("categ id: %d -- name: %s\n", c.ID, c.Name)
}

func (c *Comment) String() string {
	if c == nil {
		return "nil"
	}
	return fmt.Sprintf("comment id: %d -- Comment Message: \n%s\n", c.ID, c.Message.String())
}
