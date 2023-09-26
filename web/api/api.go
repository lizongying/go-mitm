package api

import (
	"fmt"
	"github.com/lizongying/go-mitm/static"
	"io/fs"
	"log"
	"net/http"
	"sync/atomic"
)

type Api struct {
	record      bool
	id          uint64
	client      uint64
	mux         *http.ServeMux
	messageChan chan *Message

	include        []string
	fnSetInclude   func(include string) []string
	fnClearInclude func() []string
	exclude        []string
	fnSetExclude   func(exclude string) []string
	fnClearExclude func() []string
	proxy          string
	fnSetProxy     func(proxy string) string
	fnClearProxy   func() string
}

func NewApi(messageChan chan *Message,
	include []string,
	fnSetInclude func(include string) []string,
	fnClearInclude func() []string,
	exclude []string,
	fnSetExclude func(exclude string) []string,
	fnClearExclude func() []string,
	proxy string,
	fnSetProxy func(proxy string) string,
	fnClearProxy func() string,
) (a *Api) {
	a = new(Api)
	a.mux = http.NewServeMux()
	files, _ := fs.Sub(static.Dist, "dist")
	a.mux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(files))))
	a.mux.HandleFunc("/info", a.info)
	a.mux.HandleFunc("/event", a.event)
	a.mux.HandleFunc("/action", a.action)
	a.record = true
	a.messageChan = messageChan
	a.include = include
	a.fnSetInclude = fnSetInclude
	a.fnClearInclude = fnClearInclude
	a.exclude = exclude
	a.fnSetExclude = fnSetExclude
	a.fnClearExclude = fnClearExclude
	a.proxy = proxy
	a.fnSetProxy = fnSetProxy
	a.fnClearProxy = fnClearProxy
	return
}

func (a *Api) Handler() http.Handler {
	return http.Handler(a.mux)
}
func (a *Api) info(w http.ResponseWriter, _ *http.Request) {
	info := Info{
		Record:  a.record,
		Proxy:   a.proxy,
		Exclude: a.exclude,
		Include: a.include,
	}
	_, _ = w.Write([]byte(info.String()))
	return
}
func (a *Api) action(_ http.ResponseWriter, r *http.Request) {
	record := r.URL.Query().Get("record")
	if record != "" {
		a.record = record == "true"
		return
	}
	replay := r.URL.Query().Get("replay")
	if replay != "" {
		fmt.Println("replay", replay)
		return
	}
	exclude := r.URL.Query().Get("exclude")
	if exclude != "" {
		if exclude == "-" {
			a.exclude = a.fnClearExclude()
			return
		}
		a.exclude = a.fnSetExclude(exclude)
		return
	}
	include := r.URL.Query().Get("include")
	if include != "" {
		if include == "-" {
			a.exclude = a.fnClearInclude()
			return
		}
		a.include = a.fnSetInclude(include)
		return
	}
	proxy := r.URL.Query().Get("proxy")
	if proxy != "" {
		if proxy == "-" {
			a.proxy = a.fnClearProxy()
			return
		}
		a.proxy = a.fnSetProxy(proxy)
		return
	}
	return
}
func (a *Api) event(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Panic("server not support")
	}

	client := atomic.AddUint64(&a.client, 1)
out:
	for {
		if !a.record {
			break
		}

		select {
		case message := <-a.messageChan:
			atomic.AddUint64(&a.id, 1)
			message.Id = atomic.LoadUint64(&a.id)
			_, err := fmt.Fprintf(w, "data: %s\n\n", message.String())
			if err != nil {
				fmt.Println(err)
				break out
			}
			flusher.Flush()
		default:
			if client != atomic.LoadUint64(&a.client) {
				break out
			}
		}
	}

	_, _ = fmt.Fprintf(w, "event: close\ndata: close\n\n")
	return
}
