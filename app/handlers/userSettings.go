package handlers

import (
	"errors"
	"net/http"
	"net/mail"

	"forum/app/config"
	"forum/app/templates"
	"forum/model"

	"golang.org/x/crypto/bcrypt"
)

/*
the user's settings page.  Route: /settings. Methods: GET,POST. Template: settings
*/
func SettingsPageHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ses, err := checkLoggedin(app, w, r)
		if err != nil {
			// checkLoggedin has already written error status to w
			return
		}
		if ses.LoginStatus != loggedin {
			Forbidden(app, w, r)
			return
		}

		switch r.Method {
		case http.MethodPost:
			// get data from a form
			err = r.ParseForm()
			if err != nil {
				ServerError(app, w, r, "parsing form error", err)
				return
			}

			email := r.PostFormValue(F_EMAIL)
			password := r.PostFormValue(F_PASSWORD)
			if email == "" && password == "" {
				ClientError(app, w, r, http.StatusBadRequest, "nothing to change")
				return
			}

			if email != "" {
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

				err = app.ForumData.ChangeUsersEmail(ses.User.ID, email)
				if err != nil {
					if errors.Is(err, model.ErrUniqueUserEmail) {
						w.Write([]byte("error: the email already exists"))
						return
					} else {
						ServerError(app, w, r, "changing user's email failed", err)
						return
					}
				}

				w.Write([]byte("the email has been successfully changed"))
				return
			}
			if password != "" {
				hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
				if err != nil {
					ServerError(app, w, r, "password crypting failed", err)
					return
				}

				err = app.ForumData.ChangeUsersPassword(ses.User.ID, string(hashPassword))
				if err != nil {
					ServerError(app, w, r, "changing user's password failed", err)
					return
				}

				w.Write([]byte("the password has been successfully changed"))
				return
			}

		case http.MethodGet:
			// create a page
			output := &struct {
				Session *session
			}{Session: ses}
			err = templates.ExecuteTemplate(app, w, r, "settings", output)
			if err != nil {
				ServerError(app, w, r, "tamplate executing faild", err)
				return
			}
		default:
			// only GET or PUT methods are allowed
			MethodNotAllowed(app, w, r, http.MethodGet, http.MethodPost)
		}
	}
}
