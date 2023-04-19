package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lipandr/Snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

var dbConnCounts int

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Printf("Starting application on %s", *addr)

	// open a database connection using the DSN data
	db := connectDB()
	if db == nil {
		errorLog.Fatal("Failed to connect to MySQL database")
	}

	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("MySQL is not yet ready...")
			dbConnCounts++
		} else {
			log.Println("MySQL is ready and connected!")
			return conn
		}
		if dbConnCounts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Retrying in 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}
