package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/feedmail/feedmail/app"
	M "github.com/feedmail/feedmail/models"
	W "github.com/feedmail/feedmail/workers"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	app := &app.Config{}
	app.Addr = flag.String("addr", ":3000", "TCP address for the server to listen on, in the form host:port")
	app.CsrfSecret = flag.String("csrf-secret", "xyz", "csrf-secret")
	app.CacheTag = flag.String("cache-tag", "123", "cache-tag")
	app.Domain = flag.String("domain", "localhost", "domain")
	app.Version = flag.String("version", "0.0.1-dev", "version")
	app.Runtime = flag.String("runtime", "server", "runtime") // server, worker
	app.DBConnection = flag.String("db-connection", "host=localhost user=feedmail password=feedmail dbname=feedmail port=5431 sslmode=disable TimeZone=UTC", "db-connection")
	flag.Parse()

	// Database pool
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, *app.DBConnection)
	if err != nil {
		panic(err)
	}
	defer dbPool.Close()

	// River migration
	migrator := rivermigrate.New(riverpgxv5.New(dbPool), nil)
	_, err = migrator.Migrate(ctx, rivermigrate.DirectionUp, &rivermigrate.MigrateOpts{})
	if err != nil {
		panic(err)
	}

	// Runtime
	switch *app.Runtime {
	case "server":
		// Database
		db, err := gorm.Open(postgres.Open(*app.DBConnection), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		db.AutoMigrate(&M.User{}, &M.Session{}, &M.Account{}) // tables with foreign keys last
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

		log.Printf("Server listen on %s", *app.Addr)
		log.Fatal(srv.ListenAndServe())
	case "worker":
		workers := river.NewWorkers()
		river.AddWorker(workers, &W.SortWorker{})

		app.River, err = river.NewClient(riverpgxv5.New(dbPool), &river.Config{
			Queues: map[string]river.QueueConfig{
				river.QueueDefault: {MaxWorkers: 100},
			},
			Workers: workers,
		})
		if err != nil {
			panic(err)
		}

		if err := app.River.Start(ctx); err != nil {
			panic(err)
		}

		log.Println("Worker started")
	}
}
