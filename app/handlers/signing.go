package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"forum/app/config"
	"forum/app/templates"
	"forum/model"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

/*
the signup page.  Route: /signup. Methods: POST. Template: signup
*/
func SignupPageHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only if it's notloggedin - needs wrapper

		// try to add a user
		// get data from a form
		err := r.ParseForm()
		if err != nil {
			ServerError(app, w, r, "parsing form error", err)
			return
		}

		name := r.FormValue(F_NAME)
		email := r.PostFormValue(F_EMAIL)
		password := r.PostFormValue(F_PASSWORD)
		if name == "" || email == "" || password == "" {
			ClientError(app, w, r, http.StatusBadRequest, "empty string in credential data")
			return
		}

		// check email
		// mail.ParseAddress accepts also local domens e.g witout .(dot)
		_, err = mail.ParseAddress(email)
		if err != nil {
			w.Write([]byte("error: wrong email"))
			return
		}
		// the regex allows only Internet emails, e.g. with dot-atom domain (https://www.rfc-editor.org/rfc/rfc5322.html#section-3.4)
		// if !regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`).Match([]byte(email)) {
		// 	w.Write([]byte("error: wrong email"))
		// 	return
		// }

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
		if err != nil {
			ServerError(app, w, r, "password crypting failed", err)
			return
		}

		// add a user  to DB
		id, err := app.ForumData.AddUser(name, email, hashPassword, time.Now())
		if err == nil { // the user is added - redirect to success page
			tSID, err := uuid.NewV4()
			if err != nil {
				ServerError(app, w, r, "UUID creating failed", err)
				return
			}
			expiresAt := time.Now().Add(60 * time.Second)

			// set tSID
			http.SetCookie(w, &http.Cookie{
				Name:    "tSID",
				Value:   tSID.String(),
				Expires: expiresAt,
			})
			err = app.ForumData.AddUsersSession(id, tSID.String(), expiresAt)
			if err != nil {
				ServerError(app, w, r, "adding session failed", err)
				return
			}

			// responde to JS, with status 204 it will link to /signup/success
			w.Header().Add("Location", "/signup/success")
			w.WriteHeader(204)

		} else { // adding is failed - error mesage and respond with the filled form
			var message string
			switch err {
			case model.ErrUniqueUserName:
				message = "error: the name already exists"
			case model.ErrUniqueUserEmail:
				message = "error: the email already exists"
			default:
				ServerError(app, w, r, "adding the user failed", err)
				return
			}

			// write responce to JavsScript function
			w.Write([]byte(message))
		}
	}
}

/*
the successreg page. Route: /signup/success. Methods: GET. Template: successreg
*/
func SignupSuccessPageHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ses, err := checkLoggedin(app, w, r)
		if err != nil {
			// checkLoggedin has already written error status to w
			return
		}
		if ses.LoginStatus == loggedin {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if ses.LoginStatus == experied {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// continue only if it's notloggedin
		// take tSID
		cook, err := r.Cookie("tSID")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			ClientError(app, w, r, http.StatusBadRequest, fmt.Sprintf("getting cookie tSID failed: %s, url: %s", err, r.URL))
			return
		}
		tSID := cook.Value
		// find the new user by tSID
		user, err := app.ForumData.GetUserBySession(tSID)
		if err != nil {
			if err == model.ErrNoRecord {
				ClientError(app, w, r, http.StatusBadRequest, fmt.Sprintf("a user with tSID=%s is not found", tSID))
				return
			}
			ServerError(app, w, r, "getting a user by tSID failed", err)
			return
		}
		// delete the temporary SID
		err = app.ForumData.DeleteUsersSession(user.ID)
		if err != nil {
			ServerError(app, w, r, "deleting user's session failed", err)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "tSID",
			Value:   "",
			Expires: time.Now(),
		})
		// create a page
		output := &struct {
			Session *session
			Name    string
		}{
			Session: NotloggedinSession(),
			Name:    user.Name,
		}
		err = templates.ExecuteTemplate(app, w, r, "successreg", output)
		if err != nil {
			ServerError(app, w, r, "tamplate executing faild", err)
			return
		}
	}
}

/*
the login page. Route: /login. Methods: POST. Template: signin
*/
func SigninPageHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only if it's notloggedin - needs wrapper
		// try to add a user
		err := r.ParseForm()
		if err != nil {
			ServerError(app, w, r, "parsing form error", err)
			return
		}

		name := r.PostFormValue(F_NAME)
		password := r.PostFormValue(F_PASSWORD)
		if name == "" || password == "" {
			ClientError(app, w, r, http.StatusBadRequest, "empty string in credential data")
			return
		}
		user, err := app.ForumData.GetUserByName(name)
		if err != nil {
			if errors.Is(err, model.ErrNoRecord) { // the user doesn't exist
				// write a message for JS
				w.Write([]byte("error: wrong login"))
				return
			}
			// any other errors:
			ServerError(app, w, r, "getting user for signin failed", err)
			return
		}
		// check user's password
		expectedHashPassword := user.Password
		if len(expectedHashPassword) == 0 {
			ServerError(app, w, r, "wrong data in the DB", fmt.Errorf("user's (%s) password is empty", name))
			return
		}

		err = bcrypt.CompareHashAndPassword(expectedHashPassword, []byte(password))
		if err == nil { // the password is true - create SID & redirect to the home page
			SID, err := uuid.NewV4()
			if err != nil {
				ServerError(app, w, r, "UUID creating failed", err)
				return
			}
			expiresAt := time.Now().Add(EXP_SESSION * time.Second)

			http.SetCookie(w, &http.Cookie{
				Name:    "SID",
				Value:   SID.String(),
				Expires: expiresAt,
			})
			err = app.ForumData.AddUsersSession(user.ID, SID.String(), expiresAt)
			if err != nil {
				ServerError(app, w, r, "adding session failed", err)
				return
			}

			// responde to JS, with status 204 it will link to the home page
			w.Header().Add("Location", "/")
			w.WriteHeader(204)

		} else { // the password is wrong - error mesage and respond with the filled form
			// write a message for JS
			w.Write([]byte("error: wrong password"))
		}
	}
}

/*
the logout handler. Route: /logout. Methods: any. Template: -
*/
func LogoutHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get a session
		ses, err := checkLoggedin(app, w, r)
		if err != nil {
			// checkLoggedin has already written error status to w
			return
		}

		if ses.IsLoggedin() {
			err = app.ForumData.DeleteUsersSession(ses.User.ID)
			if err != nil {
				ServerError(app, w, r, "deleting the expired session failed", err)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "SID",
				Value:   "",
				Expires: time.Now(),
			})
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}