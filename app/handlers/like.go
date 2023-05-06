package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"forum/app/application"
	"forum/model"
	"forum/model/sqlpkg"
)

type postLikeDB struct {
	dataSource            *sqlpkg.ForumModel
	id, userID, messageID int
	like                  bool
}

type commentLikeDB struct {
	dataSource            *sqlpkg.ForumModel
	id, userID, messageID int
	like                  bool
}
type liker interface {
	GetLike() error
	InsertLike(bool) error
	UpdateLike(bool) error
	DeleteLike() error
	CompareLike(bool) bool
}

/*
the liking handler. Route: /liking. Methods: POST. Template: -
*/
func LikingHandler(app *application.Application) http.HandlerFunc {
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
		w.Header().Set("Content-Type", "application.Application/json")
		fmt.Fprintf(w, `{"like": "%d", "dislike": "%d"}`, likes[model.LIKE], likes[model.DISLIKE])
	}
}

func (pl *postLikeDB) GetLike() error {
	var err error
	pl.id, pl.like, err = pl.dataSource.GetUsersPostLike(pl.userID, pl.messageID)
	return err
}

func (pl *postLikeDB) InsertLike(like bool) error {
	var err error
	pl.like = like
	pl.id, err = pl.dataSource.InsertPostLike(pl.userID, pl.messageID, pl.like)
	return err
}

func (pl *postLikeDB) UpdateLike(like bool) error {
	pl.like = like
	return pl.dataSource.UpdatePostLike(pl.id, pl.like)
}

func (pl *postLikeDB) DeleteLike() error {
	return pl.dataSource.DeletePostLike(pl.id)
}

func (pl *postLikeDB) CompareLike(like bool) bool {
	return pl.like == like
}

func (cl *commentLikeDB) GetLike() error {
	var err error
	cl.id, cl.like, err = cl.dataSource.GetUsersCommentLike(cl.userID, cl.messageID)
	return err
}

func (cl *commentLikeDB) InsertLike(like bool) error {
	var err error
	cl.like = like
	cl.id, err = cl.dataSource.InsertCommentLike(cl.userID, cl.messageID, cl.like)
	return err
}

func (cl *commentLikeDB) UpdateLike(like bool) error {
	cl.like = like
	return cl.dataSource.UpdateCommentLike(cl.id, cl.like)
}

func (cl *commentLikeDB) DeleteLike() error {
	return cl.dataSource.DeleteCommentLike(cl.id)
}

func (cl *commentLikeDB) CompareLike(like bool) bool {
	return cl.like == like
}

func setLike(liker liker, newLike bool) error {
	err := liker.GetLike()
	if err != nil {
		// if there is no like/dislike made by the user, add a new one
		if errors.Is(err, model.ErrNoRecord) {
			err := liker.InsertLike(newLike)
			if err != nil {
				return fmt.Errorf("insert data to DB failed: %s", err)
			}
		} else {
			return fmt.Errorf("getting data from DB failed: %s", err)
		}
	} else {
		if liker.CompareLike(newLike) { // if it is the same like, delete it
			err := liker.DeleteLike()
			if err != nil {
				return fmt.Errorf("deleting data from DB failed: %s", err)
			}
		} else {
			err := liker.UpdateLike(newLike)
			if err != nil {
				return fmt.Errorf("updating data in DB failed: %s", err)
			}
		}
	}
	return nil
}
