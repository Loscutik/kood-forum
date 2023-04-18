package main

import (
	"errors"
	"fmt"

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

func (pl *postLikeDB) GetLike() error {
	var err error
	pl.id, pl.like, err = pl.dataSource.GetUsersPostLike(pl.userID, pl.messageID)
	return err
}

func (pl *postLikeDB) InsertLike(like bool) error {
	var err error
	pl.like=like
	pl.id, err = pl.dataSource.InsertPostLike(pl.userID, pl.messageID, pl.like)
	return err
}

func (pl *postLikeDB) UpdateLike(like bool) error {
	pl.like=like
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
	cl.like=like
	cl.id, err = cl.dataSource.InsertCommentLike(cl.userID, cl.messageID, cl.like)
	return err
}

func (cl *commentLikeDB) UpdateLike(like bool) error {
	cl.like=like
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
