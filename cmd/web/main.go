package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

type application struct {
	address    string
	staticPath string
	db         *sql.DB
	errorLog   *log.Logger
	infoLog    *log.Logger
}

func (app *application) initApp() {
	app.errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	flag.StringVar(&app.address, "address", ":8000", "HTTP network address")
	flag.StringVar(&app.staticPath, "static_path", "./ui/static/", "Path to the static folder")
	flag.Parse()
}

func initDb() (*sql.DB, error) {
	dsn := flag.String("sql_dsn", "web:pass@/snippetbox?parseTime=True", "MySQL data source")
	flag.Parse()
	db, err := sql.Open("mysql", *dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	app := application{}
	app.initApp()

	db, err := initDb()
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer db.Close()

	app.db = db
	app.infoLog.Printf("DB stats are %+v\n", db.Stats())

	app.infoLog.Printf("Starting server on %v address...\n", app.address)
	server := &http.Server{
		Addr:     app.address,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}
	err = server.ListenAndServe()
	app.errorLog.Fatal(err)
}
