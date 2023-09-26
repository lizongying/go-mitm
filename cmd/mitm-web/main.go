package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/lizongying/go-mitm/proxy"
	"github.com/lizongying/go-mitm/web/api"
	"net/http"
)

func main() {
	midPortPtr := flag.Int("mid-port", 8082, "-mid-port proxyPort")
	webPortPtr := flag.Int("web-port", 8083, "-web-port webPort")
	includePtr := flag.String("include", "", "-include include")
	excludePtr := flag.String("exclude", "localhost;127.0.0.1", "-exclude exclude")
	proxyPtr := flag.String("proxy", "", "-proxy proxy")
	flag.Parse()

	var err error
	messageChan := make(chan *api.Message, 255)

	p, err := proxy.NewProxy(*includePtr, *excludePtr, *proxyPtr)
	if err != nil {
		panic(err)
	}
	p.SetMessageChan(messageChan)
	midSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", *midPortPtr),
		Handler:      p,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	fmt.Printf("Mid: http://localhost:%d\n", *midPortPtr)
	go func() {
		err = midSrv.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	handler := api.NewApi(messageChan, p.Include(), p.SetInclude, p.ClearInclude, p.Exclude(), p.SetExclude, p.ClearExclude, p.Proxy(), p.SetProxy, p.ClearProxy, p.Replay).Handler()
	handler = api.CrossDomain(handler)
	handler = api.Print(handler)
	srvApi := &http.Server{
		Addr:    fmt.Sprintf(":%d", *webPortPtr),
		Handler: handler,
	}
	fmt.Printf("Web: http://localhost:%d\n", *webPortPtr)
	go func() {
		err = srvApi.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	select {}
}
