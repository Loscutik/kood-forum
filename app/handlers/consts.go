package handlers

import (
	"forum/model"
)

const EXP_SESSION = 1200

// form fields
const (
	F_NAME         = "name"
	F_PASSWORD     = "password"
	F_EMAIL        = "email"
	F_CONTENT      = "content"
	F_IMAGES       = "images"
	F_AUTHORID     = "authorID"
	F_THEME        = "theme"
	F_CATEGORIESID = "categoriesID"
)

const USER_IMAGES_DIR   = "./images"

const (
	MaxFileUploadSize = 20 << 20               // 20MB
	MaxUploadSize     = 10 * MaxFileUploadSize // 10 files by 20MB
)

type likesStorage struct {
	Post, Comment string
}

var defaultLikesStorage = &likesStorage{model.POSTS_LIKES, model.COMMENTS_LIKES}
