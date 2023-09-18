package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/lizongying/go-mitm/mitm"
	"github.com/lizongying/go-mitm/web/api"
	"net/http"
)

func main() {
	filterPtr := flag.String("f", "", "-f filter")
	replacePtr := flag.Bool("r", false, "-r replace")
	proxyPtr := flag.String("p", "", "-p proxy")
	flag.Parse()

	var err error
	messageChan := make(chan *api.Message, 255)

	mux := api.NewApi(messageChan).Mux()
	handler := api.CrossDomain(mux)
	srvApi := &http.Server{
		Addr:    ":8083",
		Handler: handler,
	}
	fmt.Println("Web: http://localhost:8083")
	go func() {
		err = srvApi.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	p, err := mitm.NewProxy(*filterPtr, *proxyPtr, *replacePtr)
	if err != nil {
		panic(err)
	}
	p.SetMessageChan(messageChan)
	srv := &http.Server{
		Addr:         ":8082",
		Handler:      p,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	fmt.Println("Proxy: http://localhost:8082")
	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	select {}
}
