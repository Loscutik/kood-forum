package application

import (
	"fmt"
	"log"
	"os"

	"forum/app/templates"
	"forum/model"
	"forum/model/sqlpkg"

	"golang.org/x/crypto/bcrypt"
)

type Application struct {
	ErrLog       *log.Logger
	InfoLog      *log.Logger
	TemlateCashe templates.TemplateCache
	ForumData    *sqlpkg.ForumModel
}

func New() *Application {
	errLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime) // Creates logs of errors
	infoLogFile, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o664)
	if err != nil {
		errLog.Printf("Cannot open a log file. Error is %s\nStdout will be used for the info log ", err)
		infoLogFile = os.Stdout
	}
	infoLog := log.New(infoLogFile, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile)
	return &Application{ErrLog: errLog, InfoLog: infoLog}
}

func (app *Application) SetTemplates(path string) error {
	// create template's cashe - it keeps parsed temlates
	templates, err := templates.NewTemplateCache(path)
	if err != nil {
		return err
	}
	app.TemlateCashe = templates
	return nil
}

func (app *Application) CreateDB(fileName string) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(model.ADM_PASS), 8)
	if err != nil {
		return fmt.Errorf("password crypting failed: %v", err)
	}
	db, err := sqlpkg.CreateDB(fileName , model.ADM_NAME, model.ADM_EMAIL, string(hashPassword))
	if err != nil {
		return fmt.Errorf("creating DB faild: %v", err)
	}
	app.ForumData = &sqlpkg.ForumModel{DB: db}
	app.InfoLog.Printf("DB has been created")
	return nil
}

func (app *Application) FillTestDB(path string) error {
	hashPassword1, err := bcrypt.GenerateFromPassword([]byte("test1"), 8)
	if err != nil {
		return fmt.Errorf("password crypting failed: %v", err)
	}
	hashPassword2, err := bcrypt.GenerateFromPassword([]byte("test2"), 8)
	if err != nil {
		return fmt.Errorf("password crypting failed: %v", err)
	}
	app.InfoLog.Println("DB has been filled by examles of data")
	return app.ForumData.FillInDB(path, string(hashPassword1), string(hashPassword2))
}
