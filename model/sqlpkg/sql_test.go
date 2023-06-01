package sqlpkg

import (
	"fmt"
	"testing"
	"time"

	"forum/model"

	"github.com/mattn/go-sqlite3"
)

func TestCreateDB(t *testing.T) {
	db, err := CreateDB("database.db", model.ADM_NAME, model.ADM_EMAIL, model.ADM_PASS)
	// db, err := OpenDB("database.db","webuser","webuser")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	f := ForumModel{db}

	uss, err := f.GetUsers()
	if err != nil {
		t.Fatal(err)
	}

	for _, us := range uss {
		fmt.Println(us)
	}
	fmt.Println("------------")

	id, err := f.InsertUser("test11", "test1@email", []byte("pass"), time.Date(2023, time.March, 3, 12, 12, 21, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("--id=%d-----\n", id)
	fmt.Println("------------")
	uss, err = f.GetUsers()
	if err != nil {
		t.Fatal(err)
	}

	for _, us := range uss {
		fmt.Println(us)
	}
	fmt.Println("----end-----")
}

func TestDeletePost(t *testing.T) {
	db, err := OpenDB("database.db", "webuser", "webuser")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	f := ForumModel{db}

	q := `DELETE FROM posts  WHERE id=3`
	res, err := f.DB.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	postID, err := res.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("---id=%d-------\n", postID)
}

func TestInsertPost(t *testing.T) {
	db, err := OpenDB("database.db", "webuser", "webuser")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	f := ForumModel{db}

	id, err := f.InsertPost("theme1", "it's content1", []string{}, 1, time.Date(2023, time.March, 7, 12, 12, 21, 0, time.UTC), []int{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("---id=%d-------\n", id)
}

func TestLikes(t *testing.T) {
	db, err := OpenDB("database.db", "webuser", "webuser")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	f := ForumModel{db}

	fmt.Println("--insert 2,5 - true --")
	id, err := f.InsertPostLike(2, 5, true)
	if err != nil {
		t.Fatal(err)
	}
	err = f.printLikes("posts_likes")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("--update 2,5->false--")
	err = f.UpdatePostLike(id, false)
	if err != nil {
		t.Fatal(err)
	}
	err = f.printLikes("posts_likes")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("--delete 2,5--")
	err = f.DeletePostLike(id)
	if err != nil {
		t.Fatal(err)
	}
	err = f.printLikes("posts_likes")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetLikes(t *testing.T) {
	db, err := OpenDB("database.db", "webuser", "webuser")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	f := ForumModel{db}

	fmt.Println("--get likes for post 2--")
	likes, err := f.GetPostLikes(2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n", likes)

	fmt.Println("--get likes for post 3 by user 2--")
	id, like, err := f.GetUsersPostLike(2, 3)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("id= %d, like=%v\n", id, like)
	/*
		fmt.Println("--get likes for post 1 by user 2 (no)--")
		id,like, err = f.GetUsersPostLike(2,1)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("id= %d, like=%v\n",id, like)
	*/
	fmt.Println("--get likes for post 1 by user 3 (no)--")
	id, like, err = f.GetUsersPostLike(3, 1)
	fmt.Printf("err=%v, type err - %T\n", err, err)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("id= %d, like=%v\n", id, like)
}

func (f *ForumModel) printLikes(table string) error {
	q := `SELECT * FROM ` + table

	rows, err := f.DB.Query(q, table)
	if err != nil {
		return err
	}
	defer rows.Close()
	fmt.Printf("--id--\t--userID--\t--messageID--\t--like--\t\n")
	for rows.Next() {
		var id, userID, postID int
		var like bool
		err := rows.Scan(&id, &userID, &postID, &like)
		if err != nil {
			return err
		}
		fmt.Printf("  %d  \t    %d   \t    %d    \t\t  %v\n", id, userID, postID, like)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func TestGetPosts(t *testing.T) {
	db, err := OpenDB("database.db", "webuser", "webuser")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := ForumModel{db}

	fmt.Println("--get posts--")

	filter := &model.Filter{
		AuthorID:      0,
		CategoryID:    nil,
		LikedByUserID: 0,
	}
	posts, err := f.GetPosts(filter)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}

	fmt.Println("--get posts author 2--")
	filter = &model.Filter{
		AuthorID:      2,
		CategoryID:    nil,
		LikedByUserID: 0,
	}
	posts, err = f.GetPosts(filter)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}

	fmt.Println("--get posts category 2--")
	filter = &model.Filter{
		AuthorID:      0,
		CategoryID:    []int{2},
		LikedByUserID: 0,
	}
	posts, err = f.GetPosts(filter)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}

	fmt.Println("--get posts liked by user 1--")
	filter = &model.Filter{
		AuthorID:      0,
		CategoryID:    nil,
		LikedByUserID: 1,
	}
	posts, err = f.GetPosts(filter)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}
	fmt.Println("--get posts by author 1 category 2--")
	filter = &model.Filter{
		AuthorID:      1,
		CategoryID:    []int{2},
		LikedByUserID: 0,
	}
	posts, err = f.GetPosts(filter)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}

	fmt.Println("--get posts liked by user 1 category 2--")
	filter = &model.Filter{
		AuthorID:      0,
		CategoryID:    []int{2},
		LikedByUserID: 1,
	}
	posts, err = f.GetPosts(filter)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}

	fmt.Println("--get posts category 1,2--")
	filter = &model.Filter{
		AuthorID:      0,
		CategoryID:    []int{1, 2},
		LikedByUserID: 0,
	}
	posts, err = f.GetPosts(filter)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}
}

func TestGetPostsByCategory(t *testing.T) {
	db, err := OpenDB("database.db", "webuser", "webuser")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := ForumModel{db}

	fmt.Println("--get posts cat 2--")

	posts, err := f.GetPostsByCategory(2)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}

	fmt.Println("--get posts cat 0--")

	posts, err = f.GetPostsByCategory(0)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}
}

func TestGetPostsLikedByUser(t *testing.T) {
	db, err := OpenDB("database.db", "webuser", "webuser")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := ForumModel{db}

	fmt.Println("--get liked by user 1 --")

	posts, err := f.GetPostsLikedByUser(1)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}

	fmt.Println("--get posts liked by user 2 --")

	posts, err = f.GetPostsLikedByUser(2)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}

	fmt.Println("--get liked by user 0 --")

	posts, err = f.GetPostsLikedByUser(0)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}
}

func TestGetPostByDI(t *testing.T) {
	db, err := OpenDB("database.db", "webuser", "webuser")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := ForumModel{db}

	fmt.Println("--get post 1--")

	post, err := f.GetPostByID(1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", post.String())

	fmt.Println("--get post 3--")

	post, err = f.GetPostByID(3)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", post.String())
}

func TestInsertComment(t *testing.T) {
	db, err := OpenDB("database.db", "webuser", "webuser")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	f := ForumModel{db}

	id, err := f.InsertComment(2, "comment 2 tto post 2", []string{}, 1, time.Date(2023, time.March, 8, 12, 12, 21, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("---id=%d-------\n", id)
}

func TestAddUserSession(t *testing.T) {
	db, err := OpenDB("database.db", "webuser", "webuser")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := ForumModel{db}

	fmt.Println("--- add a session to the user 1 ---")
	fmt.Println(f.AddUsersSession(1, "ses1", time.Now()))
	fmt.Println("--- add a session to the user 10(not existing) ---")
	fmt.Println(f.AddUsersSession(10, "ses1", time.Now()))
}

func TestInsertUser(t *testing.T) {
	db, err := OpenDB("database.db", "webuser", "webuser")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := ForumModel{db}

	fmt.Println("--- add a user with the existing name ---")
	n, err := f.InsertUser("ussertets1", "emailt1", []byte("pas"), time.Now())
	fmt.Printf("n= %d, err= %v\n", n, err)
	n, err = f.AddUser("ussertetsA", "emailtA", []byte("pas"), time.Now())

	fmt.Printf("n= %d, err= %v\n", n, err)
}

func TestGetUserBySession(t *testing.T) {
	db, err := OpenDB("database.db", "webuser", "webuser")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := ForumModel{db}

	fmt.Println("------")
	u, err := f.GetUserBySession("d8ce41bc-a504-4c4d-9285-c560a4bcaa7b")
	fmt.Printf("user= %v, \nerr= %v\n", u, err)

	fmt.Println("---no that session---")
	u, err = f.GetUserBySession("d8ce41bc-a504-4c4d-9285-c560a4b")
	fmt.Printf("user= %v, \nerr= %v\n", u, err)
}

func TestGetUserByID(t *testing.T) {
	db, err := OpenDB("database.db", "webuser", "webuser")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := ForumModel{db}

	fmt.Println("---id 2---")
	u, err := f.GetUserByID(2)
	fmt.Printf("user= %v, \nerr= %v\n", u, err)

	fmt.Println("---id 4---")
	u, err = f.GetUserByID(4)
	fmt.Printf("user= %v, \nerr= %v\n", u, err)
}

func TestAuthenDB(t *testing.T) {
	db, err := OpenDB("database.db", "webuser", "webuser") // open as not admin
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var sqlconn *sqlite3.SQLiteConn
	err = sqlconn.AuthUserAdd("webuser1", "webuser", false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("----end-----")
}
