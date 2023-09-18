package main

import (
	"github.com/lizongying/go-mitm/web/api"
	"net/http"
)

func main() {
	mux := api.NewApi().Mux()
	handler := api.CrossDomain(mux)
	srv := &http.Server{
		Addr:    ":8083",
		Handler: handler,
	}
	_ = srv.ListenAndServe()
}
