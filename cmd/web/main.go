package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NPeykov/snippetbox/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
    infoLog *log.Logger
    errorLog *log.Logger
    snippets *models.SnippetModel
    templateCache map[string]*template.Template
    formDecoder *form.Decoder
    sessionManager *scs.SessionManager
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

    tmplCache, err := newTemplateCache()
    if err != nil {
        errorLog.Fatal(err)
    }

    formDecoder := form.NewDecoder()

    sessionManager := scs.New()
    sessionManager.Store = mysqlstore.New(db)
    sessionManager.Lifetime = 12 * time.Hour
    sessionManager.Cookie.Secure = true

    app := &application{
        infoLog: infoLog, 
        errorLog: errorLog,
        snippets: &models.SnippetModel{ DB: db },
        templateCache: tmplCache,
        formDecoder: formDecoder,
        sessionManager: sessionManager,
    }

    srv := &http.Server{
        Addr: *addr,
        ErrorLog: errorLog,
        Handler: app.routes(),
        IdleTimeout: time.Minute,
        ReadTimeout:   5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

	infoLog.Printf("Starting server on port %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
