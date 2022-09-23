package controller

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

// LogRequest print request headers.
func LogRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		dump, err := httputil.DumpRequest(req, false)

		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		}
		log.Printf(string(dump))

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}
