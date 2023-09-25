package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/lizongying/go-mitm/proxy"
	"net/http"
)

func main() {
	midPortPtr := flag.Int("mid-port", 8082, "-mid-port proxyPort")
	includePtr := flag.String("include", "", "-include include")
	excludePtr := flag.String("exclude", "localhost;127.0.0.1", "-exclude exclude")
	proxyPtr := flag.String("proxy", "", "-proxy proxy")
	flag.Parse()

	p, err := proxy.NewProxy(*includePtr, *excludePtr, *proxyPtr)
	if err != nil {
		return
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", *midPortPtr),
		Handler:      p,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	_ = srv.ListenAndServe()
}
