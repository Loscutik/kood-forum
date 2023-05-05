package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"forum/app/config"
	"forum/model"
)

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
