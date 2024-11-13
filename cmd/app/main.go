package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AndreFernandesDev/boilerplate_web/internal/auth"
	"github.com/AndreFernandesDev/boilerplate_web/internal/database"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type app struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	DB             *database.Queries
	sessionManager *scs.SessionManager
}

func main() {
	godotenv.Load()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	port := os.Getenv("PORT")
	if port == "" {
		errorLog.Fatal("PORT is not found in environment")
	}

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		errorLog.Fatal("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		errorLog.Fatalf("Connection to database failed: %v", err)

	}

	auth.NewAuth()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(conn)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = false

	app := &app{
		errorLog:       errorLog,
		infoLog:        infoLog,
		DB:             database.New(conn),
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:         ":" + port,
		ErrorLog:     app.errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.infoLog.Printf("Starting server on port %s", port)
	// err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
