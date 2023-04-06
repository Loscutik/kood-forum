package sqlpkg

import "time"

/*
inserts a new comment into DB, returns an ID for the comment
*/
func (f *ForumModel) InsertComment(postID int, content string, authorID int, dateCreate time.Time) (int, error) {
	q := `INSERT INTO comments (content, authorID, dateCreate, postID) VALUES (?,?,?,?)`
	res, err := f.DB.Exec(q,  content, authorID, dateCreate, postID)
	if err != nil {
		return 0, err
	}

	commentID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(commentID), nil
}
