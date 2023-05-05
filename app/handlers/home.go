package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"forum/app/config"
	"forum/app/templates"
	"forum/model"
)

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

		categories, err := app.ForumData.GetCategories()
		if err != nil {
			ServerError(app, w, r, "getting data (set of categories) from DB failed", err)
			return
		}

		// get category filters
		uQ := r.URL.Query()
		var categoryID []int
		if len(uQ[F_CATEGORIESID]) > 0 {
			for _, c := range uQ[F_CATEGORIESID] {
				id, err := strconv.Atoi(c)
				if err != nil || id <= 0 || id > len(categories) {
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
