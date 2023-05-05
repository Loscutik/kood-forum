package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/app/config"
	"forum/app/templates"
	"forum/model"
)

/*
the post's page. Route: /post/p{{Id}}. Methods: GET, POST. Template: post
*/
func PostPageHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only GET or PUT methods are allowed
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			MethodNotAllowed(app, w, r, http.MethodGet, http.MethodPost)
			return
		}

		// get the post id
		const prefix = "/post/p"
		stringID := strings.TrimPrefix(r.URL.Path, prefix)
		if stringID == r.URL.Path { // if the prefix doesn't exist
			NotFound(app, w, r)
			return
		}
		postID, err := strconv.Atoi(stringID)
		if err != nil || postID < 1 {
			ClientError(app, w, r, http.StatusBadRequest, fmt.Sprintf("wrong post id in the URL post/p: %s, err: %s", stringID, err))
			return
		}

		// get a session
		ses, err := checkLoggedin(app, w, r)
		if err != nil {
			// checkLoggedin has already written error status to w
			return
		}

		if r.Method == http.MethodPost { // => creating a comment
			// only for authorisated
			if ses.LoginStatus == experied {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			if ses.LoginStatus == notloggedin {
				Forbidden(app, w, r)
				return
			}
			// continue for the loggedin status only
			// get data from the request
			err := r.ParseForm()
			if err != nil {
				ServerError(app, w, r, "parsing form error", err)
				return
			}

			content := r.PostFormValue(F_CONTENT)

			authorID, err := strconv.Atoi(r.PostFormValue(F_AUTHORID))
			if err != nil || authorID < 1 {
				ClientError(app, w, r, http.StatusBadRequest, fmt.Sprintf("A comment creating is faild: wrong athor id: %s, err: %s", r.PostFormValue(F_AUTHORID), err))
				return
			}

			dateCreate := time.Now()
			if strings.TrimSpace(content) == "" {
				ClientError(app, w, r, http.StatusBadRequest, "comment creating failed: empty data")
				return
			}

			// add the comment to the DB
			_, err = app.ForumData.InsertComment(postID, content, authorID, dateCreate)
			if err != nil {
				ServerError(app, w, r, "insert a comment to DB failed", err)
				return
			}
		}

		// get the post from DB
		post, err := app.ForumData.GetPostByID(postID)
		if err != nil {
			if errors.Is(err, model.ErrNoRecord) {
				NotFound(app, w, r)
				return
			}
			ServerError(app, w, r, "getting a post faild", err)
			return
		}

		// create a page
		output := &struct {
			Session      *session
			Post         *model.Post
			LikesStorage *likesStorage
		}{Session: ses, Post: post, LikesStorage: defaultLikesStorage}

		err = templates.ExecuteTemplate(app, w, r, "post", output)
		if err != nil {
			ServerError(app, w, r, "tamplate executing faild", err)
			return
		}
	}
}
