package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type Tmpl struct {
	Handler  string
	Fn       string
	Format   string
	Layout   string
	Partial  bool
	Code     int
	Data     any
	CacheTag *string
	Csrf     string
}

func (c *Config) InitTemplate(custom Tmpl, r *http.Request) Tmpl {
	t := custom
	if len(custom.Format) == 0 {
		t.Format = "html"
	}
	if len(custom.Layout) == 0 {
		t.Layout = "layout/app"
	}
	if custom.Code == 0 {
		t.Code = 200
	}

	t.CacheTag = c.CacheTag
	t.Csrf = fmt.Sprintf("%s", r.Context().Value("csrf"))

	// t.Partial = false
	// t.Data = nil
	// t.Code = 200
	return t
}

func (c *Config) CacheHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		splitPath := strings.Split(r.URL.Path, ".")
		log.Printf("cache %v\nCache >> %s\nsplitPath >> %s\n", r.URL.Path, c.CacheTag, splitPath)
		if len(splitPath) > 1 {
			cleanPath := RemoveIndex(splitPath, len(splitPath)-2)
			filePath := strings.Join(cleanPath, ".")
			r.URL.Path = filePath
			log.Printf("filePath >> %s\n", r.URL.Path)
			h.ServeHTTP(w, r)
		} else {
			log.Printf("Redirecting %s to root.", r.RequestURI)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func RespondOK(w http.ResponseWriter, response string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
	return nil
}

func RespondStatus(w http.ResponseWriter, status int) error {
	w.WriteHeader(status)
	return nil
}

func RespondError(args ...interface{}) StatusError {
	code := 500
	err := errors.New("Internal Server Error")
	msg := ""

	for _, arg := range args {
		switch t := arg.(type) {
		case string:
			msg = t
		case int:
			code = t
		case error:
			err = t
		}
	}
	return StatusError{Code: code, Err: err, Msg: msg}
}

func RespondText(w http.ResponseWriter, status int, payload string) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	w.Write([]byte(payload))
	return nil
}

func RespondJSON[T any](w http.ResponseWriter, status int, payload T) error {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return nil
	}
	w.Header().Set("Content-Type", "application/activity+json")
	w.WriteHeader(status)
	w.Write([]byte(response))
	return nil
}

func (c *Config) RespondErr(w http.ResponseWriter, r *http.Request, handler string, message string) error {
	t := c.InitTemplate(Tmpl{Handler: handler, Fn: "_error", Data: message, Partial: true}, r)

	w.Header().Set("Content-Type", "text/vnd.turbo-stream.html")
	w.WriteHeader(422)

	err := c.Respond(w, r, t)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) RespondData(w http.ResponseWriter, r *http.Request, handler string, data any) error {
	t := c.InitTemplate(Tmpl{Handler: handler, Fn: "_error", Data: data}, r)

	err := c.Respond(w, r, t)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) Respond(w http.ResponseWriter, r *http.Request, t Tmpl) error {
	t = c.InitTemplate(t, r)

	log.Printf("Respond: %s/%s.%s\nFormat: %v", t.Handler, t.Fn, t.Format, t)

	if t.Handler == "" || t.Fn == "" {
		return RespondError(fmt.Errorf("handler or func missing"))
	}

	var err error
	tmplFile, _ := template.ParseFiles(fmt.Sprintf("templates/%s.%s", t.Layout, t.Format), fmt.Sprintf("templates/%s/%s.%s", t.Handler, t.Fn, t.Format))
	if t.Partial {
		tmplFile, _ = template.ParseFiles(fmt.Sprintf("templates/%s/%s.%s", t.Handler, t.Fn, t.Format))
		err = tmplFile.Execute(w, t)
	} else if len(r.Header.Get("Turbo-Frame")) > 0 {
		log.Printf("Partial or Turbo-Frame: %s\n", r.Header.Get("Turbo-Frame"))
		tmplFile, _ = template.ParseFiles(fmt.Sprintf("templates/%s/%s.%s", t.Handler, t.Fn, t.Format))
		err = tmplFile.ExecuteTemplate(w, string(r.Header.Get("Turbo-Frame")), t)
	} else {
		err = tmplFile.Execute(w, t)
	}
	if err != nil {
		return RespondError(fmt.Errorf("can not parse template: %s", err.Error()))
	}
	return nil
}
