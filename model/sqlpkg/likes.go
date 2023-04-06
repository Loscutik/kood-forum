package sqlpkg

import (
	"database/sql"
	"errors"

	"forum/model"
)

func (f *ForumModel) AddPostLike(userID, postID int, like bool) error {
	// check if there is a like from this user
	// change the like or add the new one
	return nil
}

func (f *ForumModel) GetUsersPostLike(userID, messageID int) (int, bool, error) {
	return f.getUsersLike(model.POSTS_LIKES, userID, messageID)
}

func (f *ForumModel) InsertPostLike(userID, messageID int, like bool) (int, error) {
	return f.insertLike(model.POSTS_LIKES, userID, messageID, like)
}

func (f *ForumModel) UpdatePostLike(id int, like bool) error {
	return f.updateLike(model.POSTS_LIKES, id, like)
}

func (f *ForumModel) DeletePostLike(id int) error {
	return f.deleteLike(model.POSTS_LIKES, id)
}

func (f *ForumModel) GetUsersCommentLike(userID, messageID int) (int, bool, error) {
	return f.getUsersLike(model.COMMENTS_LIKES, userID, messageID)
}

func (f *ForumModel) InsertCommentLike(userID, messageID int, like bool) (int, error) {
	return f.insertLike(model.COMMENTS_LIKES, userID, messageID, like)
}

func (f *ForumModel) UpdateCommentLike(id int, like bool) error {
	return f.updateLike(model.COMMENTS_LIKES, id, like)
}

func (f *ForumModel) DeleteCommentLike(id int) error {
	return f.deleteLike(model.COMMENTS_LIKES, id)
}

func (f *ForumModel) getUsersLike(tableName string, userID, messageID int) (int, bool, error) {
	var id int
	var like bool
	q := `SELECT id,like FROM ` + tableName + ` WHERE userID=? AND messageID=?`
	row := f.DB.QueryRow(q, userID, messageID)

	err := row.Scan(&id, &like)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, model.ErrNoRecord
		}
		return 0, false, err
	}

	return id, like, nil
}

func (f *ForumModel) insertLike(tableName string, userID, messageID int, like bool) (int, error) {
	q := `INSERT INTO ` + tableName + ` (userID, messageID, like) VALUES (?,?,?)`
	res, err := f.DB.Exec(q, userID, messageID, like)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (f *ForumModel) updateLike(tableName string, id int, like bool) error {
	q := `UPDATE ` + tableName + ` SET like=? WHERE id=?`
	res, err := f.DB.Exec(q, like, id)
	if err != nil {
		return err
	}

	return f.checkUnique(res)
}

func (f *ForumModel) deleteLike(tableName string, id int) error {
	q := `DELETE FROM ` + tableName + ` WHERE id=?`
	res, err := f.DB.Exec(q, id)
	if err != nil {
		return err
	}

	return f.checkUnique(res)
}
