package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"forum/app/application"
	"forum/app/templates"
	"forum/model"
)

/*
the add post page. Route: /addpost. Methods: GET. Template: addpost
*/
func AddPostPageHandler(app *application.Application) http.HandlerFunc {
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
		err = templates.ExecuteTemplate(app.TemlateCashe, w, r, "addpost", output)
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
func PostCreatorHandler(app *application.Application) http.HandlerFunc {
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
		if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
			ServerError(app, w, r, "Cannot parse multipart form", err)
			return
		}

		dateCreate := time.Now()
		theme := r.PostFormValue(F_THEME)
		content := r.PostFormValue(F_CONTENT)
		imageFiles := r.MultipartForm.File[F_IMAGES]
		imagesTmpDir := fmt.Sprintf("%s/tmp_%d%d_%d", USER_IMAGES_DIR, dateCreate.Second(), dateCreate.Nanosecond(), rand.Intn(100))
		if imageFiles != nil {
			err := os.Mkdir(imagesTmpDir, 0o750)
			if err != nil && !os.IsExist(err) {
				ServerError(app, w, r, "Can't create tmp directory", err)
				return
			}
		}
		var imagesList []string
		for _, fileHeader := range imageFiles {
			newFileName, err := uploadFile(MaxFileUploadSize, fileHeader, imagesTmpDir)
			imagesList = append(imagesList, newFileName)
			if err != nil {
				ServerError(app, w, r, "Can't upload the file", err)
				return
			}

		}

		authorID, err := strconv.Atoi(r.PostFormValue(F_AUTHORID))
		if err != nil || authorID < 1 {
			ClientError(app, w, r, http.StatusBadRequest, fmt.Sprintf("wrong athor id: %s, err: %s", r.PostFormValue(F_AUTHORID), err))
			return
		}

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

		if strings.TrimSpace(theme) == "" || (strings.TrimSpace(content) == "" && len(imagesList) == 0) || len(categories) == 0 || categoriesID[0] == 0 {
			ClientError(app, w, r, http.StatusBadRequest, "post creating failed: empty data")
			return
		}

		// add post to the DB
		id, err := app.ForumData.InsertPost(theme, content, imagesList, authorID, dateCreate, categoriesID)
		if err != nil {
			ServerError(app, w, r, "insert to DB failed", err)
			return
		}

		os.Rename(imagesTmpDir, fmt.Sprintf("p%d", id))
		// redirect to the post page
		http.Redirect(w, r, "/post/p"+strconv.Itoa(id), http.StatusSeeOther)
	}
}
