package app

import (
	"fmt"
	"log"
	"net/http"
)

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
	Msg  string
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}

func (se StatusError) Message() string {
	return se.Msg
}

func HandleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case StatusError:
		log.Printf("StatusError HTTP %d - %s - %s", e.Status(), e, e.Message())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.Status())
		if e.Message() != "" {
			w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", e.Message())))
		} else {
			//http.Error(w, e.Error(), e.Status())
			w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", e)))
		}
	case Error:
		log.Printf("Error HTTP %d - %s", e.Status(), e)
		http.Error(w, e.Error(), e.Status())
	default:
		log.Printf(fmt.Sprintf("default Error HTTP %v", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}
