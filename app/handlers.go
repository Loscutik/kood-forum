package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/model"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// TODO forbiden.html is needed

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
	F_MESSAGEID    = "messageID"
	F_FROMURL      = "fromURL"
	F_LIKE         = "like"
)

/*
The handler of the main page. Route: /. Methods: GET. Template: home
*/
func (app *application) homePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.NotFound(w, r)
		return
	}

	// only GET method is allowed
	if r.Method != http.MethodGet {
		app.MethodNotAllowed(w, r, http.MethodGet)
		return
	}

	ses, err := app.checkLoggedin(w, r)
	if err != nil {
		// checkLoggedin has already written error status to w
		return
	}

	uQ := r.URL.Query()
	var categoryID []int
	if len(uQ["category"]) > 0 {
		for _, c := range uQ["category"] {
			id, err := strconv.Atoi(c)
			if err != nil {
				app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("wrong category id in the filter request: %s, err: %s", c, err))
				return
			}

			categoryID = append(categoryID, id)
		}
	}
	
	filter := &model.Filter{
		AuthorID:      0,
		CategoryID:    categoryID,
		LikedByUserID: 0,
	}
	if ses.IsLoggedin() {
		switch {
		case uQ.Get("author") != "":
			filter.AuthorID = ses.User.ID
		case uQ.Get("likedby") != "":
			filter.LikedByUserID = ses.User.ID
		}
	}
	posts, err := app.forumData.GetPosts(filter)
	if err != nil {
		app.ServerError(w, r, "getting data from DB failed", err)
		return
	}

	categories, err := app.forumData.GetCategories()
	if err != nil {
		app.ServerError(w, r, "getting data (set of categories) from DB failed", err)
		return
	}

	// create a page
	output := &struct {
		Session    *session
		Posts      []*model.Post
		Categories []*model.Category
	}{Session: ses, Posts: posts, Categories: categories}
	// Assembling the page from templates
	app.executeTemplate(w, r, "home", output)
}

/*
the signup page.  Route: /signup. Methods: GET, POST. Template: signup
*/
func (app *application) signupPageHandler(w http.ResponseWriter, r *http.Request) {
	ses, err := app.checkLoggedin(w, r)
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

	// TODO redo it in the way that the handler will respond to a javascript function, not execute templates
	// continue only if it's notloggedin
	switch r.Method {
	case http.MethodGet:
		// respond with the empty form
		app.executeTemplate(w, r, "signup", nil)
	case http.MethodPost:
		// try to add a user
		// get data from a form
		err := r.ParseForm()
		if err != nil {
			app.ServerError(w, r, "parsing form error", err)
			return
		}

		name := r.PostFormValue(F_NAME)
		// TODO check correct email & compare passwords in javascript
		email := r.PostFormValue(F_EMAIL)
		password := r.PostFormValue(F_PASSWORD)
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
		if err != nil {
			app.ServerError(w, r, "password crypting failed", err)
			return
		}

		// add a user  to DB
		fmt.Printf("name=%s, email= %s, hashPassword = %s\n",name, email, hashPassword)
		id, err := app.forumData.AddUser(name, email, hashPassword, time.Now())
		if err == nil { // the user is added - redirect to success page
			fmt.Printf("added, id=%d, err = %v",id,err)
			tSID, err := uuid.NewV4()
			if err != nil {
				app.ServerError(w, r, "UUID creating failed", err)
				return
			}
			expiresAt := time.Now().Add(60 * time.Second)

			// set tSID
			http.SetCookie(w, &http.Cookie{
				Name:    "tSID",
				Value:   tSID.String(),
				Expires: expiresAt,
			})
			err = app.forumData.AddUsersSession(id, tSID.String(), expiresAt)
			if err != nil {
				app.ServerError(w, r, "adding session failed", err)
				return
			}
			// go to success page
			http.Redirect(w, r, "/signup/success", http.StatusSeeOther)

		} else { // adding is failed - error mesage and respond with the filled form
			var message string
			switch err {
			case model.ErrUniqueUserName:
				message = "the name already exists"
			case model.ErrUniqueUserEmail:
				message = "the email already exists"
			default:
				app.ServerError(w, r, "adding a user failed", err)
				return
			}

			// create a page
			output := &struct {
				Session    *session
				Name, Email, Password, Message string
			}{
				Session: NotloggedinSession(),
				Name:     name,
				Email:    email,
				Password: password,
				Message:  message,
			}
			app.executeTemplate(w, r, "signup", output)
		}

	default:
		app.MethodNotAllowed(w, r, http.MethodGet+", "+http.MethodPost)
	}
}

/*
the successreg page. Route: /signup/success. Methods: GET. Template: successreg
*/
func (app *application) signupSuccessPageHandler(w http.ResponseWriter, r *http.Request) {
	ses, err := app.checkLoggedin(w, r)
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
		app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("getting cookie tSID failed: %s, url: %s", err, r.URL))
		return
	}
	tSID := cook.Value
	// find the new user by tSID
	user, err := app.forumData.GetUserBySession(tSID)
	if err != nil {
		if err == model.ErrNoRecord {
			app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("a user with tSID=%s is not found", tSID))
			return
		}
		app.ServerError(w, r, "getting a user by tSID failed", err)
		return
	}
	// delete the temporary SID
	err = app.forumData.DeleteUsersSession(user.ID)
	if err != nil {
		app.ServerError(w, r, "deleting user's session failed", err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "tSID",
		Value:   "",
		Expires: time.Now(),
	})
	// create a page
	output := &struct {
		Session    *session
		Name string
	}{
		Session: NotloggedinSession(),
		Name:     user.Name,
	}
	app.executeTemplate(w, r, "successreg", output)
}

/*
the login page. Route: /login. Methods: GET, POST. Template: signin
*/
func (app *application) signinPageHandler(w http.ResponseWriter, r *http.Request) {
	ses, err := app.checkLoggedin(w, r)
	if err != nil {
		// checkLoggedin has already written error status to w
		return
	}
	if ses.IsLoggedin() {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// TODO redo it in the way that the handler will respond to a javascript function, not execute templates
	// continue if it's neither notloggedin nor expiried
	switch r.Method {
	case http.MethodGet:
		// respond with the empty form
		app.executeTemplate(w, r, "signin", nil)
	case http.MethodPost:
		// try to add a user
		err := r.ParseForm()
		if err != nil {
			app.ServerError(w, r, "parsing form error", err)
			return
		}

		name := r.PostFormValue(F_NAME)
		password := r.PostFormValue(F_PASSWORD)
		user, err := app.forumData.GetUserByName(name)
		if err != nil {
			if errors.Is(err, model.ErrNoRecord) { // the user doesn't exist
				// create a page
				output := &struct {
					Session    *session
					Name, Message string
				}{
					Session: NotloggedinSession(),
					Name:    name,
					Message: "wrong login",
				}
				app.executeTemplate(w, r, "signin", output)
				return
			}
			// any other errors:
			app.ServerError(w, r, "getting user for signin failed", err)
			return
		}
		// check user's password
		expectedHashPassword := user.Password
		if len(expectedHashPassword) == 0 {
			app.ServerError(w, r, "wrong data in the DB", fmt.Errorf("user's (%s) password is empty", name))
			return
		}

		err = bcrypt.CompareHashAndPassword(expectedHashPassword, []byte(password))
		if err == nil { // the password is true - create SID & redirect to the home page
			SID, err := uuid.NewV4()
			if err != nil {
				app.ServerError(w, r, "UUID creating failed", err)
				return
			}
			expiresAt := time.Now().Add(EXP_SESSION * time.Second)

			http.SetCookie(w, &http.Cookie{
				Name:    "SID",
				Value:   SID.String(),
				Expires: expiresAt,
			})
			err = app.forumData.AddUsersSession(user.ID, SID.String(), expiresAt)
			if err != nil {
				app.ServerError(w, r, "adding session failed", err)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)

		} else { // the password is wrong - error mesage and respond with the filled form
			// create a page
			output := &struct {
				Session    *session
				Name, Message string
			}{
				Session: NotloggedinSession(),
				Name:    name,
				Message: "wrong password",
			}
			app.executeTemplate(w, r, "signin", output)
		}

	default:
		app.MethodNotAllowed(w, r, http.MethodGet+", "+http.MethodPost)
	}
}

/*
the userinfo page. Route: /userinfo/. Methods: GET. Template: userinfo
*/
func (app *application) userPageHandler(w http.ResponseWriter, r *http.Request) {
	// only GET method is allowed
	if r.Method != http.MethodGet {
		app.MethodNotAllowed(w, r, http.MethodGet)
		return
	}

	// get a user id from URL
	const prefix = "/user/@"
	stringID := strings.TrimPrefix(r.URL.Path, prefix)
	if stringID == r.URL.Path { // if the prefix doesn't exist
		app.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(stringID)
	if err != nil || id < 1 {
		app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("wrong user id in a URL /userinfo/@: %s, err: %s", stringID, err))
		return
	}
	// get a user from DB
	user, err := app.forumData.GetUserByID(id)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			app.NotFound(w, r)
			return
		}
		app.ServerError(w, r, "getting a user faild", err)
		return
	}
	user.Password = []byte("")

	// get a session
	ses, err := app.checkLoggedin(w, r)
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
	app.executeTemplate(w, r, "userinfo", output)
}

/*
the post's page. Route: /post/. Methods: GET, POST. Template: post
*/
func (app *application) postPageHandler(w http.ResponseWriter, r *http.Request) {
	// only GET or PUT methods are allowed
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		app.MethodNotAllowed(w, r, http.MethodGet+", "+http.MethodPost)
		return
	}

	// get post id
	const prefix = "/post/p"
	stringID := strings.TrimPrefix(r.URL.Path, prefix)
	if stringID == r.URL.Path { // if the prefix doesn't exist
		app.NotFound(w, r)
		return
	}
	postID, err := strconv.Atoi(stringID)
	if err != nil || postID < 1 {
		app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("wrong post id in the URL post/p: %s, err: %s", stringID, err))
		return
	}

	// get a session
	ses, err := app.checkLoggedin(w, r)
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
			app.Forbidden(w, r)
			return
		}
		// continue for the loggedin status only
		// get data from the request
		err := r.ParseForm()
		if err != nil {
			app.ServerError(w, r, "parsing form error", err)
			return
		}

		content := r.PostFormValue(F_CONTENT)

		authorID, err := strconv.Atoi(r.PostFormValue(F_AUTHORID))
		if err != nil || authorID < 1 {
			app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("A comment creating is faild: wrong athor id: %s, err: %s", r.PostFormValue(F_AUTHORID), err))
			return
		}

		dateCreate := time.Now()
		if content == "" {
			app.ClientError(w, r, http.StatusBadRequest, "comment creating failed: empty data")
			return
		}

		// add the comment to the DB
		_, err = app.forumData.InsertComment(postID, content, authorID, dateCreate)
		if err != nil {
			app.ServerError(w, r, "insert a comment to DB failed", err)
			return
		}
	}

	// get the post from DB
	post, err := app.forumData.GetPostByID(postID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			app.NotFound(w, r)
			return
		}
		app.ServerError(w, r, "getting a post faild", err)
		return
	}

	// create a page
	output := &struct {
		Session *session
		Post    *model.Post
	}{Session: ses, Post: post}

	app.executeTemplate(w, r, "post", output)
}

/*
the add post page. Route: /addpost. Methods: GET. Template: addpost
*/
func (app *application) addPostPageHandler(w http.ResponseWriter, r *http.Request) {
	// only GET methode is allowed
	if r.Method != http.MethodGet {
		app.MethodNotAllowed(w, r, http.MethodGet)
		return
	}

	// get a session
	ses, err := app.checkLoggedin(w, r)
	if err != nil {
		// checkLoggedin has already written error status to w
		return
	}

	if ses.LoginStatus != loggedin {
		app.Forbidden(w, r)
		return
	}

	// create a page
	output := &struct {
		Session *session
	}{Session: ses}

	app.executeTemplate(w, r, "addpost", output)
}

/*
the post creating handler. Route: /post/create. Methods: POST. Template: -
*/
func (app *application) postCreatorHandler(w http.ResponseWriter, r *http.Request) {
	// only POST method is allowed
	if r.Method != http.MethodPost {
		app.MethodNotAllowed(w, r, http.MethodPost)
		return
	}

	// get a session
	ses, err := app.checkLoggedin(w, r)
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
		app.Forbidden(w, r)
		return
	}
	// continue for the loggedin status only
	// get data from the request
	err = r.ParseForm()
	if err != nil {
		app.ServerError(w, r, "parsing form error", err)
		return
	}

	theme := r.PostFormValue(F_THEME)
	content := r.PostFormValue(F_CONTENT)

	authorID, err := strconv.Atoi(r.PostFormValue(F_AUTHORID))
	if err != nil || authorID < 1 {
		app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("wrong athor id: %s, err: %s", r.PostFormValue(F_AUTHORID), err))
		return
	}

	dateCreate := time.Now()

	categories := r.PostForm[F_CATEGORIESID]
	categoriesID := make([]int, len(categories))
	for i, c := range categories {
		id, err := strconv.Atoi(c)
		if err != nil || id < 1 {
			app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("wrong cathegory id: %s, err: %s", r.PostFormValue(F_AUTHORID), err))
			return
		}
		categoriesID[i] = id
	}

	if theme == "" || content == "" || len(categories) == 0 || categoriesID[0] == 0 {
		app.ClientError(w, r, http.StatusBadRequest, "post creating failed: empty data")
		return
	}

	// add post to the DB
	id, err := app.forumData.InsertPost(theme, content, authorID, dateCreate, categoriesID)
	if err != nil {
		app.ServerError(w, r, "insert to DB failed", err)
		return
	}
	// redirect to the post page
	http.Redirect(w, r, "/post/p"+strconv.Itoa(id), http.StatusSeeOther)
}

/*
the liking handler. Route: /liking. Methods: POST. Template: -
*/
func (app *application) likingHandler(w http.ResponseWriter, r *http.Request) {
	// only PUT method is allowed
	if r.Method != http.MethodPost {
		app.MethodNotAllowed(w, r, http.MethodPost)
		return
	}

	// get a session
	ses, err := app.checkLoggedin(w, r)
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
		app.Forbidden(w, r)
		return
	}

	// continue for the loggedin status only
	// get data from the request
	err = r.ParseForm()
	if err != nil {
		app.ServerError(w, r, "parsing form error", err)
		return
	}

	messageType := r.PostFormValue("messageType")
	messageID, err := strconv.Atoi(r.PostFormValue(F_MESSAGEID))
	if err != nil || messageID < 1 {
		app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("wrong message id: %s, err: %s", r.PostFormValue(F_MESSAGEID), err))
		return
	}

	// TODO set answer to javascript insted of redirect
	fromURL := r.PostFormValue(F_FROMURL) // TODO form this url with #id - id of what was liked

	newLike, err := strconv.ParseBool(r.PostFormValue(F_LIKE))
	if err != nil {
		app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("wrong value of the flag 'like': %s, err: %s", r.PostFormValue(F_LIKE), err))
		return
	}

	// add or change like into the DB
	switch messageType {
	case model.POSTS_LIKES:
		err = setLike(&postLikeDB{dataSource: app.forumData, userID: ses.User.ID, messageID: messageID}, newLike)
		if err != nil {
			app.ServerError(w, r, "setting a post like faild", err)
			return
		}
	case model.COMMENTS_LIKES:
		err = setLike(&commentLikeDB{dataSource: app.forumData, userID: ses.User.ID, messageID: messageID}, newLike)
		if err != nil {
			app.ServerError(w, r, "setting a comment like faild", err)
			return
		}
	default:
		app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("wrong type of a message: %s, err: %s", messageType, err))
		return
	}

	// redirect to the fromURL
	// TODO set answer to javascript insted of redirect
	http.Redirect(w, r, fromURL, http.StatusSeeOther)
}

/*
the logout handler. Route: /logout. Methods: any. Template: -
*/
func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// get a session
	ses, err := app.checkLoggedin(w, r)
	if err != nil {
		// checkLoggedin has already written error status to w
		return
	}

	if ses.IsLoggedin() {
		err = app.forumData.DeleteUsersSession(ses.User.ID)
		if err != nil {
			app.ServerError(w, r, "deleting the expired session failed", err)
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
