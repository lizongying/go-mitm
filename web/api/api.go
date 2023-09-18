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
}

func NewApi(messageChan chan *Message) (a *Api) {
	a = new(Api)
	a.mux = http.NewServeMux()
	files, _ := fs.Sub(static.Dist, "dist")
	a.mux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(files))))
	a.mux.HandleFunc("/info", a.info)
	a.mux.HandleFunc("/event", a.event)
	a.mux.HandleFunc("/action", a.action)
	a.record = true
	a.messageChan = messageChan
	return
}

func (a *Api) Mux() *http.ServeMux {
	return a.mux
}
func (a *Api) info(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())
	if a.record {
		_, _ = w.Write([]byte(`true`))
	} else {
		_, _ = w.Write([]byte(`false`))
	}
	return
}
func (a *Api) action(_ http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())
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

	return
}
func (a *Api) event(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())
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
	fmt.Println("close")
	return
}
