package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"forum/app/config"
	"forum/app/templates"
	"forum/model"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

const EXP_SESSION = 1200

// form fields
const (
	F_NAME         = "name"
	F_PASSWORD     = "password"
	F_EMAIL        = "email"
	F_CONTENT      = "content"
	F_AUTHORID     = "authorID"
	F_THEME        = "theme"
	F_CATEGORIESID = "categoriesID"
)

type likesStorage struct {
	Post, Comment string
}

var defaultLikesStorage = &likesStorage{model.POSTS_LIKES, model.COMMENTS_LIKES}

/*
The handler of the main page. Route: /. Methods: GET. Template: home
*/
func HomePageHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const (
			AUTHOR    = "author"
			LIKEBY    = "likedby"
			DISLIKEBY = "dislikedby"
		)

		if r.URL.Path != "/" {
			NotFound(app, w, r)
			return
		}

		// only GET method is allowed
		if r.Method != http.MethodGet {
			MethodNotAllowed(app, w, r, http.MethodGet)
			return
		}

		ses, err := checkLoggedin(app, w, r)
		if err != nil {
			// checkLoggedin has already written error status to w
			return
		}

		// get category filters
		uQ := r.URL.Query()
		var categoryID []int
		if len(uQ[F_CATEGORIESID]) > 0 {
			for _, c := range uQ[F_CATEGORIESID] {
				id, err := strconv.Atoi(c)
				if err != nil {
					ClientError(app, w, r, http.StatusBadRequest, fmt.Sprintf("wrong category id in the filter request: %s, err: %s", c, err))
					return
				}

				categoryID = append(categoryID, id)
			}
		}

		filter := &model.Filter{
			AuthorID:         0,
			CategoryID:       categoryID,
			LikedByUserID:    0,
			DisLikedByUserID: 0,
		}

		// get author's filters
		if ses.IsLoggedin() {
			if uQ.Get(AUTHOR) != "" {
				filter.AuthorID = ses.User.ID
			}
			if uQ.Get(LIKEBY) != "" {
				filter.LikedByUserID = ses.User.ID
			}
			if uQ.Get(DISLIKEBY) != "" {
				filter.DisLikedByUserID = ses.User.ID
			}

		}
		posts, err := app.ForumData.GetPosts(filter)
		if err != nil {
			ServerError(app, w, r, "getting data from DB failed", err)
			return
		}

		categories, err := app.ForumData.GetCategories()
		if err != nil {
			ServerError(app, w, r, "getting data (set of categories) from DB failed", err)
			return
		}

		// create a page
		output := &struct {
			Session      *session
			Posts        []*model.Post
			Categories   []*model.Category
			Filter       *model.Filter
			LikesStorage *likesStorage
		}{Session: ses, Posts: posts, Categories: categories, Filter: filter, LikesStorage: defaultLikesStorage}
		// Assembling the page from templates
		err = templates.ExecuteTemplate(app, w, r, "home", output)
		if err != nil {
			ServerError(app, w, r, "tamplate executing faild", err)
			return
		}
	}
}

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
		if !regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`).Match([]byte(email)) {
			w.Write([]byte("error: wrong email"))
			return
		}

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
the userinfo page. Route: /userinfo/@{{Id}}. Methods: GET. Template: userinfo
*/
func UserPageHandler(app *config.Application) http.HandlerFunc {
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
		err = templates.ExecuteTemplate(app, w, r, "userinfo", output)
		if err != nil {
			ServerError(app, w, r, "tamplate executing faild", err)
			return
		}
	}
}

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

			// check email
			if email != "" {
				if !regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`).Match([]byte(email)) {
					w.Write([]byte("error: wrong email"))
					return
				}

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
			if content == "" {
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

/*
the add post page. Route: /addpost. Methods: GET. Template: addpost
*/
func AddPostPageHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only GET methode is allowed
		if r.Method != http.MethodGet {
			MethodNotAllowed(app, w, r, http.MethodGet)
			return
		}

		// get a session
		ses, err := checkLoggedin(app, w, r)
		if err != nil {
			// checkLoggedin has already written error status to w
			return
		}

		if ses.LoginStatus != loggedin {
			Forbidden(app, w, r)
			return
		}

		categories, err := app.ForumData.GetCategories()
		if err != nil {
			ServerError(app, w, r, "getting data (set of categories) from DB failed", err)
			return
		}

		// create a page
		output := &struct {
			Session    *session
			Categories []*model.Category
		}{Session: ses, Categories: categories}
		err = templates.ExecuteTemplate(app, w, r, "addpost", output)
		if err != nil {
			ServerError(app, w, r, "tamplate executing faild", err)
			return
		}
		if err != nil {
			ServerError(app, w, r, "tamplate executing faild", err)
			return
		}
	}
}

/*
the post creating handler. Route: /post/create. Methods: POST. Template: -
*/
func PostCreatorHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only POST method is allowed
		if r.Method != http.MethodPost {
			MethodNotAllowed(app, w, r, http.MethodPost)
			return
		}

		// get a session
		ses, err := checkLoggedin(app, w, r)
		if err != nil {
			// checkLoggedin has already written error status to w
			return
		}

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
		err = r.ParseForm()
		if err != nil {
			ServerError(app, w, r, "parsing form error", err)
			return
		}

		theme := r.PostFormValue(F_THEME)
		content := r.PostFormValue(F_CONTENT)

		authorID, err := strconv.Atoi(r.PostFormValue(F_AUTHORID))
		if err != nil || authorID < 1 {
			ClientError(app, w, r, http.StatusBadRequest, fmt.Sprintf("wrong athor id: %s, err: %s", r.PostFormValue(F_AUTHORID), err))
			return
		}

		dateCreate := time.Now()

		categories := r.PostForm[F_CATEGORIESID]
		categoriesID := make([]int, len(categories))
		for i, c := range categories {
			id, err := strconv.Atoi(c)
			if err != nil || id < 1 {
				ClientError(app, w, r, http.StatusBadRequest, fmt.Sprintf("wrong cathegory id: %s, err: %s", r.PostFormValue(F_AUTHORID), err))
				return
			}
			categoriesID[i] = id
		}

		if theme == "" || content == "" || len(categories) == 0 || categoriesID[0] == 0 {
			ClientError(app, w, r, http.StatusBadRequest, "post creating failed: empty data")
			return
		}

		// add post to the DB
		id, err := app.ForumData.InsertPost(theme, content, authorID, dateCreate, categoriesID)
		if err != nil {
			ServerError(app, w, r, "insert to DB failed", err)
			return
		}
		// redirect to the post page
		http.Redirect(w, r, "/post/p"+strconv.Itoa(id), http.StatusSeeOther)
	}
}

/*
the liking handler. Route: /liking. Methods: POST. Template: -
*/
func LikingHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only POST method is allowed
		if r.Method != http.MethodPost {
			MethodNotAllowed(app, w, r, http.MethodPost)
			return
		}

		// get a session
		ses, err := checkLoggedin(app, w, r)
		if err != nil {
			// checkLoggedin has already written an error status to w
			return
		}

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
		body, err := io.ReadAll(r.Body)
		if err != nil {
			ClientError(app, w, r, http.StatusBadRequest, fmt.Sprintf("error during reading the liking request: %s", err))
			return
		}

		var likeData struct {
			MessageType string
			MessageID   string
			Like        string
		}
		err = json.Unmarshal(body, &likeData)
		if err != nil {
			ClientError(app, w, r, http.StatusBadRequest, fmt.Sprintf("error during unmarshal the data from the liking request: %s", err))
			return
		}

		// convert data from string
		messageID, err := strconv.Atoi(likeData.MessageID)
		if err != nil || messageID < 1 {
			ClientError(app, w, r, http.StatusBadRequest, fmt.Sprintf("wrong message id: %s, err: %s", likeData.MessageID, err))
			return
		}
		newLike, err := strconv.ParseBool(likeData.Like)
		if err != nil {
			ClientError(app, w, r, http.StatusBadRequest, fmt.Sprintf("wrong value of the flag 'like': %s, err: %s", likeData.Like, err))
			return
		}

		// add or change the like into the DB
		switch likeData.MessageType {
		case model.POSTS_LIKES:
			err = setLike(&postLikeDB{dataSource: app.ForumData, userID: ses.User.ID, messageID: messageID}, newLike)
			if err != nil {
				ServerError(app, w, r, "setting a post like faild", err)
				return
			}
		case model.COMMENTS_LIKES:
			err = setLike(&commentLikeDB{dataSource: app.ForumData, userID: ses.User.ID, messageID: messageID}, newLike)
			if err != nil {
				ServerError(app, w, r, "setting a comment like faild", err)
				return
			}
		default:
			ClientError(app, w, r, http.StatusBadRequest, fmt.Sprintf("wrong type of a message: %s, err: %s", likeData.MessageType, err))
			return
		}

		// get the new number of likes/dislikes
		likes, err := app.ForumData.GetLikes(likeData.MessageType, messageID)
		if err != nil {
			ServerError(app, w, r, "getting likes faild", err)
			return
		}
		// write responce in JSON
		w.Header().Set("Content-Type", "config.Application/json")
		fmt.Fprintf(w, `{"like": "%d", "dislike": "%d"}`, likes[model.LIKE], likes[model.DISLIKE])
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