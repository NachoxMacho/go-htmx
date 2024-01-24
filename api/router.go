package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

func basic(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Welcome to the api")
}

func bad(w http.ResponseWriter, r *http.Request) {
	panic("doh no")
}

func timeQuery(w http.ResponseWriter, r *http.Request) {
	tz := r.URL.Query().Get("tz")
	loc := time.Local
	if tz != "" {
		var err error
		loc, err = time.LoadLocation(tz)
		if err != nil {
			panic(err)
		}
	}
	body := struct {
		Time string `json:"time"`
	}{Time: time.Now().In(loc).Format(time.RFC3339)}
	json.NewEncoder(w).Encode(body)
	w.Header().Set("Content-Type", "application/json")
}

func Recovery(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() { // recover from panic
			if err := recover(); err != nil { // recover from panic
				stack := debug.Stack()
				log.Printf("%s %s: panic: %v\n%s", r.Method, r.URL, err, stack)
				// write 500 status code and "internal server error" message to response so it doesn't hang
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("internal server error"))
			}
		}()
		h.ServeHTTP(w, r)
	}
}

func Mux(pattern string, mux *http.ServeMux) {
	pattern = strings.TrimSuffix(pattern, "/")

	mux.HandleFunc(pattern+"/basic", basic)
	mux.HandleFunc(pattern+"/bad", Recovery(http.HandlerFunc(bad)))
	mux.HandleFunc(pattern+"/timeQuery", Recovery(http.HandlerFunc(timeQuery)))
}
