package api

import (
	"encoding/json"
	"fmt"
	"github.com/lizongying/go-mitm/proxy"
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
	messageChan chan *proxy.Message
	lanIp       string
	internetIp  string
	proxyPort   int
	server      *proxy.Proxy
}

func NewApi(messageChan chan *proxy.Message,
	lanIp string,
	internetIp string,
	proxyPort int,
	server *proxy.Proxy,
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
	a.lanIp = lanIp
	a.internetIp = internetIp
	a.proxyPort = proxyPort
	a.server = server
	return
}

func (a *Api) Handler() http.Handler {
	return http.Handler(a.mux)
}
func (a *Api) info(w http.ResponseWriter, _ *http.Request) {
	info := Info{
		Record:     a.record,
		Proxy:      a.server.Proxy(),
		Exclude:    a.server.Exclude(),
		Include:    a.server.Include(),
		LanIp:      a.lanIp,
		InternetIp: a.internetIp,
		ProxyPort:  a.proxyPort,
		Replace:    a.server.Replace(),
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
		var message proxy.Message
		_ = json.Unmarshal([]byte(replay), &message)
		a.server.Replay(message)
		return
	}
	exclude := r.URL.Query().Get("exclude")
	if exclude != "" {
		if exclude == "-" {
			a.server.ClearExclude()
			return
		}
		a.server.SetExclude(exclude)
		return
	}
	include := r.URL.Query().Get("include")
	if include != "" {
		if include == "-" {
			a.server.ClearInclude()
			return
		}
		a.server.SetInclude(include)
		return
	}
	p := r.URL.Query().Get("proxy")
	if p != "" {
		if p == "-" {
			a.server.ClearProxy()
			return
		}
		a.server.SetProxy(p)
		return
	}
	replace := r.URL.Query().Get("replace")
	if replace != "" {
		if replace == "-" {
			a.server.ClearReplace()
			return
		}
		a.server.SetReplace(replace)
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
