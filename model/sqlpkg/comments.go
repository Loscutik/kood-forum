package sqlpkg

import (
	"strings"
	"time"
)

/*
inserts a new comment into DB, returns an ID for the comment
*/
func (f *ForumModel) InsertComment(postID int, content string, images []string, authorID int, dateCreate time.Time) (int, error) {
	q := `INSERT INTO comments (content, images, authorID, dateCreate, postID) VALUES (?,?,?,?,?)`
	res, err := f.DB.Exec(q, content, strings.Join(images, ","), authorID, dateCreate, postID)
	if err != nil {
		return 0, err
	}

	commentID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(commentID), nil
}
