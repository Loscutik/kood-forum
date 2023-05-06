package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"forum/app/application"
	"forum/app/templates"
	"forum/model/sqlpkg"
)

const DB_Name = "forumDB.db"

func main() {
	// app keeps all dependenses used by handlers
	app := application.New()
	err := app.SetTemplates(templates.TEMPLATES_PATH)
	if err != nil {
		app.ErrLog.Fatalln(err)
	}

	port, pristinDB, testDB, err := parseArgs()
	if err != nil {
		app.ErrLog.Fatalln(err)
	}

	// init DB pool DB_Name
	_, err = os.Stat(DB_Name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			createAndFillTestDB(app)
		} else {
			app.ErrLog.Fatalln(err)
		}
	} else {
		switch {
		case testDB:
			if os.Rename(DB_Name, DB_Name+".bak") != nil {
				app.ErrLog.Fatalln("cannot rename the DB file")
			}
			createAndFillTestDB(app)
		case  pristinDB:
			// rename DB file
			if os.Rename(DB_Name, DB_Name+".bak") != nil {
				app.ErrLog.Fatalln("cannot rename the DB file")
			}

			err := app.CreateDB(DB_Name)
			if err != nil {
				app.ErrLog.Fatalln(err)
			}
		default:
			db, err := sqlpkg.OpenDB(DB_Name, "webuser", "webuser")
			if err != nil {
				app.ErrLog.Fatalln(err)
			}
			app.ForumData = &sqlpkg.ForumModel{DB: db}
		}
	}
	defer app.ForumData.DB.Close()

	// Starting the web server
	server := &http.Server{
		Addr:     ":" + port,
		ErrorLog: app.ErrLog,
		Handler:  routers(app),
	}
	fmt.Printf("Starting server at http://www.localhost:%s\n", port)
	app.InfoLog.Printf("Starting server at port %s\n", port)
	if err := server.ListenAndServe(); err != nil {
		app.ErrLog.Fatal(err)
	}
}

// Parses the program's arguments to obtain the server port. If no arguments found, it uses the 8080 port by default
// Usage: go run .  --port=PORT_NUMBER
func parseArgs() (port string, pristinDB bool, testDB bool, err error) {
	usage := `wrong arguments
     Usage: go run ./app [OPTIONS]
     OPTIONS: 
            --port=PORT_NUMBER
            --p=PORT_NUMBER
            --pristin to drop the existing DB and create the new one from scratch
            --testdb drop the existing DB and start with the test DB`
	flag.StringVar(&port, "port", "8080", "server port")
	flag.StringVar(&port, "p", "8080", "server port (shorthand)")
	flag.BoolVar(&pristinDB, "pristin", false, "--pristin if you want drop the existing DB and create the new one from scratch")
	flag.BoolVar(&testDB, "testdb", false, "--testdb if you want drop the existing DB and start with the test DB")
	flag.Parse()
	if flag.NArg() > 0 {
		return "", false, false, fmt.Errorf(usage)
	}
	_, err = strconv.ParseUint(port, 10, 16)
	if err != nil {
		return "", false, false, fmt.Errorf("error: port must be a 16-bit unsigned number ")
	}
	return
}

func createAndFillTestDB(app *application.Application) {
	err := app.CreateDB(DB_Name)
	if err != nil {
		app.ErrLog.Fatalln(err)
	}
	err = app.FillTestDB("model/sqlpkg/testDB.sql")
	if err != nil {
		app.ErrLog.Fatalln(err)
	}
}
