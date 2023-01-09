package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

    "github.com/NPeykov/snippetbox/internal/models"
    _ "github.com/go-sql-driver/mysql"
)

type application struct {
    infoLog *log.Logger
    errorLog *log.Logger
    snippets *models.SnippetModel
}

func main() {
    addr := flag.String("addr", ":4000", "HTTP application port")
    dsn  := flag.String("dsn", "web:1234@/snippetbox?parseTime=true", "HTTP application port")
    flag.Parse()

    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate | log.Ltime)
    errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)

    db, err := openDB(*dsn)
    if err != nil {
        errorLog.Fatal(err)
    }
    defer db.Close()

    app := &application{
        infoLog: infoLog, 
        errorLog: errorLog,
        snippets: &models.SnippetModel{ DB: db },
    }

    srv := &http.Server{
        Addr: *addr,
        ErrorLog: errorLog,
        Handler: app.routes(),
    }

	infoLog.Printf("Starting server on port %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
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
