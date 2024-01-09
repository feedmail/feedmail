package app

import (
	"net/http"
)

type Config struct {
	Addr       *string
	Router     *http.ServeMux
	User       *string
	Password   *string
	CsrfSecret *string
	CacheTag   *string
	DB         DB
	Domain     *string
	Version    *string
}
