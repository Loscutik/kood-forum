package sqlpkg

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"forum/model"
)

/*
inserts a new post into DB, returns an ID for the post
*/
func (f *ForumModel) InsertPost(theme, content string, authorID int, dateCreate time.Time, categoriesID []int) (int, error) {
	q := `INSERT INTO posts (theme, content, authorID, dateCreate) VALUES (?,?,?,?)`
	res, err := f.DB.Exec(q, theme, content, authorID, dateCreate)
	if err != nil {
		return 0, err
	}

	postID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	q = `INSERT INTO post_categories (categoryID, postID) VALUES (?,?)`
	for i := 1; i < len(categoriesID); i++ {
		q += `,(?,?)`
	}
	insertData := make([]any, 2*len(categoriesID))
	for i := 0; i < len(categoriesID); i++ {
		insertData[2*i] = categoriesID[i]
		insertData[2*i+1] = int(postID)
	}
	res, err = f.DB.Exec(q, insertData...)
	if err != nil {
		return 0, err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

/*
search in the DB a post by the given ID
*/
func (f *ForumModel) GetPostByID(id int) (*model.Post, error) {
	query := `SELECT p.id, p.theme, p.content, p.authorID, u.name, u.dateCreate, c.id, c.name,  p.dateCreate, 
				 count(CASE WHEN pl.like THEN TRUE END), count(CASE WHEN NOT pl.like THEN TRUE END) 
			  FROM posts p
 			  LEFT JOIN users u ON u.id=p.authorID
			  LEFT JOIN post_categories pc ON pc.postID=p.id
			  LEFT JOIN categories c ON c.id=pc.categoryID
			  LEFT JOIN posts_likes pl ON pl.messageID=p.id 
			  WHERE p.id = ?		 
			  GROUP BY c.id;
		`

	// exequting the query
	var rows *sql.Rows
	var err error
	rows, err = f.DB.Query(query, id, id)
	if err != nil {
		return nil, err
	}

	// parsing the query's result
	var post *model.Post
	var category *model.Category

	// add the first post without condition
	if rows.Next() {
		post, category, err = rowScanForPostByID(rows)
		if err != nil {
			return nil, err
		}

		post.Categories = append(post.Categories, category)
	} else {
		return nil, model.ErrNoRecord
	}

	for rows.Next() {
		// add categories only
		postTmp, categoryTmp, err := rowScanForPostByID(rows)
		if err != nil {
			return nil, err
		}

		if postTmp.ID != post.ID {
			return nil, fmt.Errorf("select failed: two different posts by one ID: %d, %d", post.ID, postTmp.ID)
		}
		post.Categories = append(post.Categories, categoryTmp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()

	query = `-- select comments.
		SELECT c.id, c.content, c.authorID, u.name, u.dateCreate, c.dateCreate, 
			count(CASE WHEN cl.like THEN TRUE END), count(CASE WHEN NOT cl.like THEN TRUE END) 
	    FROM comments c
		LEFT JOIN users u ON u.id=c.authorID
	    LEFT JOIN comments_likes cl ON cl.messageID=c.id 
		WHERE c.postID = ?		 
		GROUP BY c.id;
		`
	rows, err = f.DB.Query(query, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// get comments
	for rows.Next() {
		comment := &model.Comment{}
		comment.Message.Author = &model.User{}
		comment.Message.Likes = make([]int, model.N_LIKES)

		// parse the row with fields:
		// c.id, c.content, c.authorID, u.name, u.dateCreate, c.dateCreate,
		// count(CASE WHEN cl.like THEN TRUE END), count(CASE WHEN NOT cl.like THEN TRUE END)
		err := rows.Scan(&comment.ID,
			&comment.Message.Content,
			&comment.Message.Author.ID, &comment.Message.Author.Name, &comment.Message.Author.DateCreate,
			&comment.Message.DateCreate,
			&comment.Message.Likes[model.LIKE], &comment.Message.Likes[model.DISLIKE],
		)
		if err != nil {
			return nil, err
		}

		post.Comments = append(post.Comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	post.CommentsQuantity=len(post.Comments)

	return post, nil
}

/*
scan and prefilles an item of modelPost for getPostByID
*/
func rowScanForPostByID(rows *sql.Rows) (*model.Post, *model.Category, error) {
	post := &model.Post{}
	post.Message.Likes = make([]int, model.N_LIKES)
	post.Message.Author = &model.User{}
	category := &model.Category{}

	// parse the row with fields:
	// p.id, p.theme, p.content, p.authorID, u.name, u.dateCreate, c.id, c.name,  p.dateCreate,
	// count(CASE WHEN pl.like THEN TRUE END), count(CASE WHEN NOT pl.like THEN TRUE END)
	err := rows.Scan(&post.ID, &post.Theme,
		&post.Message.Content,
		&post.Message.Author.ID, &post.Message.Author.Name, &post.Message.Author.DateCreate,
		&category.ID, &category.Name,
		&post.Message.DateCreate,
		&post.Message.Likes[model.LIKE], &post.Message.Likes[model.DISLIKE],
	)

	return post, category, err
}

/*
returns posts created by author with the given ID, if authorID==0, returns all posts
*/
func (f *ForumModel) GetPosts(filter *model.Filter) ([]*model.Post, error) {
	v := reflect.ValueOf(*filter)
	for _, field := range reflect.VisibleFields(reflect.TypeOf(*filter)) {
		// if either of the fields !=0 add conditions to the query
		if !v.FieldByIndex(field.Index).IsZero() {
			condition := ` WHERE `
			arguments := []any{}
			if filter.AuthorID != 0 {
				condition += ` u.ID= ? AND `
				arguments = append(arguments, filter.AuthorID)
			}

			if len(filter.CategoryID) != 0 {

				condition += ` p.id IN (SELECT postID FROM post_categories pc  WHERE `
				for _, c := range filter.CategoryID {
					condition += ` pc.categoryID = ? OR `
					arguments = append(arguments, c)
				}
				condition = strings.TrimSuffix(condition, `OR `)
				condition += `GROUP BY pc.postID) AND `
			}

			if filter.LikedByUserID != 0 {
				condition += ` p.id IN (SELECT messageID FROM posts_likes pl  WHERE pl.userID = ? AND pl.like=true) AND `
				arguments = append(arguments, filter.LikedByUserID)
			}
			
			if filter.DisLikedByUserID != 0 {
				condition += ` p.id IN (SELECT messageID FROM posts_likes pl  WHERE pl.userID = ? AND pl.like=false) AND `
				arguments = append(arguments, filter.DisLikedByUserID)
			}
			
			condition = strings.TrimSuffix(condition, `AND `)
			return f.getPostsByCondition(condition, arguments)
		}
	}

	// all the fields were empty
	return f.getPostsByCondition("", nil)
}

/*
returns posts that have got the given category
*/
func (f *ForumModel) GetPostsByCategory(category int) ([]*model.Post, error) {
	return f.getPostsByCondition(` WHERE p.id IN (SELECT postID FROM post_categories pc  WHERE pc.categoryID = ?) `, []any{category})
}

/*
returns posts that have got the given category
*/
func (f *ForumModel) GetPostsLikedByUser(userID int) ([]*model.Post, error) {
	return f.getPostsByCondition(` WHERE p.id IN (SELECT messageID FROM posts_likes pl  WHERE pl.userID = ? AND pl.like=true) `, []any{userID})
}

/*
addes the condition to a query and run it. Returnes found posts
*/
func (f *ForumModel) getPostsByCondition(condition string, arguments []any) ([]*model.Post, error) {
	query := `SELECT p.id, p.theme, p.content, p.authorID, u.name, u.dateCreate, c.id, c.name,  p.dateCreate, 
				(SELECT count(id) FROM comments cm WHERE cm.postID=p.id),
				count(CASE WHEN pl.like THEN TRUE END), count(CASE WHEN NOT pl.like THEN TRUE END) 
		  FROM posts p
		  LEFT JOIN users u ON u.id=p.authorID
		  LEFT JOIN post_categories pc ON pc.postID=p.id
		  LEFT JOIN categories c ON c.id=pc.categoryID
		  LEFT JOIN posts_likes pl ON pl.messageID=p.id 
		` + condition +
		` GROUP BY p.id, c.id 
		ORDER BY p.dateCreate DESC, p.id, c.id
		`
	// exequting the query
	var rows *sql.Rows
	var err error
	rows, err = f.DB.Query(query, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// parsing the query's result
	var posts []*model.Post
	authors := make(map[int]*model.User)
	postCounter := 0 // the number of the last added post

	// add the first post without condition
	if rows.Next() {
		post, category, author, err := rowScanForPosts(rows)
		if err != nil {
			return nil, err
		}
		addNewPostStruct(&posts, post, category, author, &authors)
	}

	for rows.Next() {
		post, category, author, err := rowScanForPosts(rows)
		if err != nil {
			return nil, err
		}

		// found out do we need to add a new post or to add a category to the previouse post
		// if the next row contains the same postID not create new post, just add a category to the u post
		if post.ID == posts[postCounter].ID {
			posts[postCounter].Categories = append(posts[postCounter].Categories, category)
		} else {
			addNewPostStruct(&posts, post, category, author, &authors)
			postCounter++
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

/*
scan and prefilles an item of modelPost for getPosts
*/
func rowScanForPosts(rows *sql.Rows) (*model.Post, *model.Category, *model.User, error) {
	post := &model.Post{}
	post.Message.Likes = make([]int, model.N_LIKES)
	author := &model.User{}
	category := &model.Category{}

	// parse the row with fields:
	// p.id, p.theme, p.content,  p.authorID, u.name, u.dateCreate, c.id, c.name,  p.dateCreate, count (cm.id),
	// count(CASE WHEN pl.like THEN TRUE END), count(CASE WHEN NOT pl.like THEN TRUE END)
	err := rows.Scan(&post.ID, &post.Theme, &post.Message.Content,
		&author.ID, &author.Name, &author.DateCreate,
		&category.ID, &category.Name,
		&post.Message.DateCreate,
		&post.CommentsQuantity,
		&post.Message.Likes[model.LIKE], &post.Message.Likes[model.DISLIKE],
	)
	return post, category, author, err
}

/*
creates an item of modelPost type and addes it to the slice. Used in the getPostsByCondition function
*/
func addNewPostStruct(posts *[]*model.Post, post *model.Post, category *model.Category, author *model.User, authors *map[int]*model.User) {
	post.Categories = append(post.Categories, category)

	// find out if the author in the current row is found before, if yes, keep that previouse one
	if existingAuthor, ok := (*authors)[author.ID]; ok {
		post.Message.Author = existingAuthor
	} else {
		post.Message.Author = author
		(*authors)[author.ID] = author
	}
	*posts = append(*posts, post)
}
