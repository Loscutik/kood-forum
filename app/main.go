package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"forum/app/config"
	"forum/app/templates"
	"forum/model"
	"forum/model/sqlpkg"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Creates logs of what happened
	errLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime) // Creates logs of errors
	infoLogFile, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o664)
	if err != nil {
		errLog.Printf("Cannot open a log file. Error is %s\nStdout will be used for the info log ", err)
		infoLogFile = os.Stdout
	}
	infoLog := log.New(infoLogFile, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile)

	// create template's cashe - it keeps parsed temlates
	templates, err := templates.NewTemplateCache(templates.TEMPLATES_PATH)
	if err != nil {
		errLog.Fatal(err)
	}

	// init DB pool "forumDB.db"
	var db *sql.DB
	_, err = os.Stat("forumDB.db")
	if errors.Is(err, os.ErrNotExist) {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(model.ADM_PASS), 8)
		if err != nil {
			errLog.Fatal("password crypting failed: ", err)
			return
		}
		db, err = sqlpkg.CreateDB("forumDB.db", model.ADM_NAME, model.ADM_EMAIL, string(hashPassword))
		if err != nil {
			errLog.Fatal(err)
		}
		infoLog.Printf("DB has created in")
	} else {
		db, err = sqlpkg.OpenDB("forumDB.db", "webuser", "webuser")
		if err != nil {
			errLog.Fatal(err)
		}
	}
	defer db.Close()

	// app keeps all dependenses used by handlers
	app := &config.Application{
		ErrLog:       errLog,
		InfoLog:      infoLog,
		TemlateCashe: templates,
		ForumData:    &sqlpkg.ForumModel{DB: db},
	}

	port, err := parseArgs()
	if err != nil {
		errLog.Fatal(err)
	}
	// Starting the web server
	server := &http.Server{
		Addr:     ":" + *port,
		ErrorLog: app.ErrLog,
		Handler:  routers(app),
	}
	fmt.Printf("Starting server at http://localhost:%s\n", *port)
	infoLog.Printf("Starting server at port %s\n", *port)
	if err := server.ListenAndServe(); err != nil {
		errLog.Fatal(err)
	}
}

// Parses the program's arguments to obtain the server port. If no arguments found, it uses the 8080 port by default
// Usage: go run .  --port=PORT_NUMBER
func parseArgs() (*string, error) {
	port := flag.String("port", "8080", "server port")
	flag.Parse()
	if flag.NArg() > 0 {
		return nil, fmt.Errorf("wrong arguments\nUsage: go run .  --port=PORT_NUMBER")
	}
	_, err := strconv.ParseUint(*port, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("error: port must be a 16-bit unsigned number ")
	}
	return port, nil
}
