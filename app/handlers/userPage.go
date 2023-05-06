package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/app/application"
	"forum/app/templates"
	"forum/model"
)

/*
the userinfo page. Route: /userinfo/@{{Id}}. Methods: GET. Template: userinfo
*/
func UserPageHandler(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only GET method is allowed
		if r.Method != http.MethodGet {
			MethodNotAllowed(app, w, r, http.MethodGet)
			return
		}

		// get a user id from URL
		const prefix = "/userinfo/@"
		stringID := strings.TrimPrefix(r.URL.Path, prefix)
		if stringID == r.URL.Path { // if the prefix doesn't exist
			NotFound(app, w, r)
			return
		}
		id, err := strconv.Atoi(stringID)
		if err != nil || id < 1 {
			ClientError(app, w, r, http.StatusBadRequest, fmt.Sprintf("wrong user id in a URL /userinfo/@: %s, err: %s", stringID, err))
			return
		}
		// get a user from DB
		user, err := app.ForumData.GetUserByID(id)
		if err != nil {
			if errors.Is(err, model.ErrNoRecord) {
				NotFound(app, w, r)
				return
			}
			ServerError(app, w, r, "getting a user faild", err)
			return
		}
		user.Password = []byte("")

		// get a session
		ses, err := checkLoggedin(app, w, r)
		if err != nil {
			// checkLoggedin has already written error status to w
			return
		}

		// create a page
		// if reg and name==ses.name - mypage - else -shortpage
		output := &struct {
			Session *session
			AllInfo bool
			User    *model.User
		}{ses, false, user}
		if ses.IsLoggedin() && ses.User.ID == user.ID {
			output.AllInfo = true
		}

		// Assembling the page from templates
		err = templates.ExecuteTemplate(app.TemlateCashe, w, r, "userinfo", output)
		if err != nil {
			ServerError(app, w, r, "tamplate executing faild", err)
			return
		}
	}
}