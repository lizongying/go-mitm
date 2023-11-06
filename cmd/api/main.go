package main

import (
	"github.com/lizongying/go-mitm/web/api"
	"net/http"
)

func main() {
	handler := api.NewApi(nil, "", "", 8082, nil).Handler()
	handler = api.CrossDomain(handler)
	srv := &http.Server{
		Addr:    ":8083",
		Handler: handler,
	}
	_ = srv.ListenAndServe()
}
