package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"forum/model"

	"github.com/gofrs/uuid"
)

type loginStatus byte

const (
	loggedin loginStatus = iota
	experied
	notloggedin
)

const TIME_BEFORE_AFTER_REFRESH =30

type session struct {
	LoginStatus loginStatus
	User        *model.User
}

func (s *session) isExpired() bool {
	exp := s.User.ExpirySession.Time
	return exp.Before(time.Now())
}

func (s *session) timeToExpired() time.Duration {
	exp := s.User.ExpirySession.Time
	return time.Until(exp)
}

func (s *session) IsLoggedin() bool {
	return s.LoginStatus==loggedin
}

func NotloggedinSession() *session{
	return &session{notloggedin,nil}
}

/*
returns session which contains status of login and uses's data if it's logged in.
If it is left lrss than 30 sec to expiried time, it will refresh the session
If an error occurs it will response to the client with error status and return the error
*/
func (app *application) checkLoggedin(w http.ResponseWriter, r *http.Request) (*session, error) {
	session := &session{notloggedin, nil}
	cook, err := r.Cookie("SID")
	if err != nil && err != http.ErrNoCookie {
		app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("getting cookie SID failed: %s, url: %s", err, r.URL))
		return nil, err
	}
	if err == http.ErrNoCookie || cook.Value == "" {
		return session, nil
	}

	// there is a SID
	SID := cook.Value

	user, err := app.forumData.GetUserBySession(SID)
	if err != nil {
		if err == model.ErrNoRecord {
			return session, nil
		}
		app.ServerError(w, r, "getting a user by SID failed", err)
		return nil, err
	}
	session.User = user

	if session.isExpired() {
		// delete the session & return expiried status
		session.User = nil
		err := app.forumData.DeleteUsersSession(user.ID)
		if err != nil {
			app.ServerError(w, r, "deleting the expired session failed", err)
			return nil, err
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "SID",
			Value:   "",
			Expires: time.Now(),
		})
		session.LoginStatus = experied
		return session, nil
	}

	if session.timeToExpired() < TIME_BEFORE_AFTER_REFRESH*time.Second {
		// refresh the session
		newSID, err := uuid.NewV4()
		if err != nil {
			app.ServerError(w, r, "UUID creating failed", err)
			return nil, err
		}
		expiresAt := time.Now().Add(EXP_SESSION * time.Second)

		http.SetCookie(w, &http.Cookie{
			Name:    "SID",
			Value:   newSID.String(),
			Expires: expiresAt,
		})

		err = app.forumData.AddUsersSession(user.ID, newSID.String(), expiresAt)
		if err != nil {
			app.ServerError(w, r, "adding session failed", err)
			return nil, err
		}
		session.User.Session = newSID.String()
		session.User.ExpirySession=sql.NullTime{Time: expiresAt, Valid: true}
	}

	// user was found and his time was not expired or renewed:
	session.LoginStatus = loggedin
	return session, nil
}
