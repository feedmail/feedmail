package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/feedmail/feedmail/app"
	M "github.com/feedmail/feedmail/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	app := &app.Config{}
	app.Addr = flag.String("addr", ":3000", "TCP address for the server to listen on, in the form host:port")
	app.CsrfSecret = flag.String("csrf-secret", "xyz", "csrf-secret")
	app.CacheTag = flag.String("cache-tag", "123", "cache-tag")
	app.Domain = flag.String("domain", "localhost", "domain")
	app.Version = flag.String("version", "0.0.1-dev", "version")
	flag.Parse()

	// Database
	db, err := gorm.Open(sqlite.Open("db/dev.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&M.Account{}, &M.User{}, &M.Session{})
	app.DB.Client = db

	// Router
	app.Router = http.NewServeMux()
	fs := app.CacheHandler(http.FileServer(http.Dir("./static")))
	app.Router.Handle("/static/", http.StripPrefix("/static/", fs))
	InitRoutes(app)

	// Server
	srv := &http.Server{
		Addr:    *app.Addr,
		Handler: app.Router,
		//TLSConfig:    cfg,
		//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Printf("Listen on %s", *app.Addr)
	log.Fatal(srv.ListenAndServe())
}
