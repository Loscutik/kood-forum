package config

import (
	"html/template"
	"log"

	"forum/model/sqlpkg"
)

type Application struct {
	ErrLog       *log.Logger
	InfoLog      *log.Logger
	TemlateCashe map[string]*template.Template
	ForumData    *sqlpkg.ForumModel
}
