package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/noonacedia/sourcepaste/internal/models"
)

type application struct {
	address        string
	staticPath     string
	db             *sql.DB
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       models.SnippetInterface
	users          models.UserInterface
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {
	app := application{}
	app.initApp()
	defer app.db.Close()

	app.infoLog.Printf("Starting server on %v address...\n", app.address)
	server := &http.Server{
		Addr:         app.address,
		ErrorLog:     app.errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	app.errorLog.Fatal(err)
}

func (app *application) initApp() {
	app.errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	flag.StringVar(&app.address, "address", ":8000", "HTTP network address")
	flag.StringVar(&app.staticPath, "static_path", "./ui/static/", "Path to the static folder")
	flag.Parse()

	db, err := initDb()
	if err != nil {
		app.errorLog.Fatal(err)
	}
	app.db = db
	app.snippets = &models.SnippetModel{DB: app.db}
	app.users = &models.UserModel{DB: app.db}

	templateCache, err := newTemplateCache()
	if err != nil {
		app.errorLog.Fatal(err)
	}
	app.templateCache = templateCache

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	app.sessionManager = sessionManager
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
