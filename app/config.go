package app

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

type Config struct {
	Addr         *string
	Router       *http.ServeMux
	User         *string
	Password     *string
	CsrfSecret   *string
	CacheTag     *string
	DB           DB
	Domain       *string
	Version      *string
	DBConnection *string
	Runtime      *string
	River        *river.Client[pgx.Tx]
}
