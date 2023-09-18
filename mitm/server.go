package mitm

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/lizongying/go-mitm/static"
	"github.com/lizongying/go-mitm/web/api"
	"io"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Proxy struct {
	rootCert     *x509.Certificate
	rootKey      *rsa.PrivateKey
	privateKey   *rsa.PrivateKey
	listener     *Listener
	srv          *http.Server
	serialNumber int64
	replace      bool
	proxy        *url.URL
	filter       *regexp.Regexp
	messageChan  chan *api.Message
	exclude      []string
	logger       *slog.Logger
}

func (p *Proxy) SetMessageChan(messageChan chan *api.Message) {
	p.messageChan = messageChan
}
func (p *Proxy) getCertificate(domain string) (cert *tls.Certificate, err error) {
	atomic.AddInt64(&p.serialNumber, 1)
	serverTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(p.serialNumber),
		Subject: pkix.Name{
			CommonName: domain,
		},
		NotBefore: time.Now().AddDate(0, 0, -1),
		NotAfter:  time.Now().AddDate(1, 0, 0),
	}
	ip := net.ParseIP(domain)
	if ip != nil {
		serverTemplate.IPAddresses = []net.IP{ip}
	} else {
		serverTemplate.DNSNames = []string{domain}
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, serverTemplate, p.rootCert, &p.privateKey.PublicKey, p.rootKey)
	if err != nil {
		return
	}

	cert = &tls.Certificate{
		PrivateKey:  p.privateKey,
		Certificate: [][]byte{certBytes},
	}
	return
}
func (p *Proxy) doReplace(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>%s</title>
</head>
<body>
  %s
</body>
</html>
`, r.Host, r.URL.String())))
}
func (p *Proxy) doRequest(w http.ResponseWriter, r *http.Request) {
	if p.proxy != nil {
		r.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(p.proxy.String())))
	}

	if r.URL.Host == "" {
		r.URL.Host = r.Host
	}
	if r.URL.Scheme == "" {
		r.URL.Scheme = "https"
	}

	//fmt.Println(strings.Repeat("#", 100))
	//fmt.Println("Request:")
	//requestDump, err := httputil.DumpRequest(r, true)
	//if err != nil {
	//	fmt.Println("Error dumping request:", err)
	//	return
	//}
	//fmt.Println(string(bytes.TrimSpace(requestDump)))

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

	//getBody, err := r.GetBody()
	//if err != nil {
	//	return
	//}
	//reqBody, err := io.ReadAll(getBody)

	begin := time.Now()
	response, err := http.DefaultTransport.RoundTrip(r)
	spend := uint16(time.Now().Sub(begin).Milliseconds())
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	copyHeader(w.Header(), response.Header)
	w.WriteHeader(response.StatusCode)

	//fmt.Println(strings.Repeat("#", 100))
	//fmt.Println("Response:")
	//responseDump, err := httputil.DumpResponse(response, true)
	//if err != nil {
	//	fmt.Println("Error dumping response:", err)
	//	return
	//}
	//fmt.Println(string(bytes.TrimSpace(responseDump)))

	var size int64
	var respBody string
	contentTypes := response.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(contentTypes), "image") {
		size, _ = io.Copy(w, response.Body)
	} else {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			http.Error(w, "Failed to read response body", http.StatusInternalServerError)
			return
		}

		s, _ := w.Write(bodyBytes)
		size = int64(s)

		if response.Header.Get("Content-Encoding") == "deflate" {
			reader := flate.NewReader(bytes.NewReader(bodyBytes))
			defer func() {
				err = reader.Close()
				if err != nil {
					return
				}
			}()

			bodyBytes, err = io.ReadAll(reader)
			if err != nil {
				return
			}
		}
		if response.Header.Get("Content-Encoding") == "br" {
			bodyBytes, err = io.ReadAll(brotli.NewReader(bytes.NewReader(bodyBytes)))
			if err != nil {
				return
			}
		}
		if response.Header.Get("Content-Encoding") == "deflate" {
			reader := flate.NewReader(bytes.NewReader(bodyBytes))
			defer func() {
				if reader != nil {
					err = reader.Close()
					if err != nil {
						return
					}
				}
			}()
			bodyBytes, err = io.ReadAll(reader)
			if err != nil {
				return
			}
		}
		if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(bytes.NewReader(bodyBytes))
			defer func() {
				if reader != nil {
					err = reader.Close()
					if err != nil {
						return
					}
				}
			}()
			if err != nil {
				return
			}
			bodyBytes, err = io.ReadAll(reader)
			if err != nil {
				return
			}
		}

		respBody = string(bodyBytes)
	}

	go func(r *http.Request, response *http.Response) {
		reqHeader := make(map[string]string)
		for k := range r.Header {
			reqHeader[k] = r.Header.Get(k)
		}
		respHeader := make(map[string]string)
		for k := range response.Header {
			respHeader[k] = response.Header.Get(k)
		}

		reqTrailer := make(map[string]string)
		for k := range r.Trailer {
			reqTrailer[k] = r.Trailer.Get(k)
		}
		respTrailer := make(map[string]string)
		for k := range response.Trailer {
			respTrailer[k] = response.Trailer.Get(k)
		}

		reqCookie := make(map[string]string)
		for _, v := range r.Cookies() {
			reqCookie[v.Name] = v.Raw
		}
		respCookie := make(map[string]string)
		for _, v := range response.Cookies() {
			respCookie[v.Name] = v.Raw
		}

		reqTls := make(map[string]string)
		if r.TLS != nil {
			reqTls["ServerName"] = r.TLS.ServerName
			reqTls["NegotiatedProtocol"] = r.TLS.NegotiatedProtocol
			reqTls["Version"] = fmt.Sprintf("%d", r.TLS.Version)
			reqTls["Unique"] = string(r.TLS.TLSUnique)
			reqTls["CipherSuite"] = fmt.Sprintf("%d", r.TLS.CipherSuite)
		}

		respTls := make(map[string]string)
		if response.TLS != nil {
			respTls["ServerName"] = response.TLS.ServerName
			respTls["NegotiatedProtocol"] = response.TLS.NegotiatedProtocol
			version := "Unknown"
			switch response.TLS.Version {
			case tls.VersionTLS10:
				version = "1.0"
			case tls.VersionTLS11:
				version = "1.1"
			case tls.VersionTLS12:
				version = "1.2"
			case tls.VersionTLS13:
				version = "1.3"
			}
			respTls["Version"] = version
			respTls["Unique"] = base64.StdEncoding.EncodeToString(response.TLS.TLSUnique)
			respTls["CipherSuite"] = fmt.Sprintf("%d", response.TLS.CipherSuite)
		}

		contentType := contentTypes
		for _, v := range strings.Split(contentTypes, ";") {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			if strings.Contains(strings.ToLower(v), "charset=") {
				continue
			}
			contentType = v
			break
		}

		p.logger.Info("Response", "StatusCode", response.StatusCode, r.Method, r.URL.String(), "contentType", contentType)

		p.messageChan <- &api.Message{
			Url:         r.URL.String(),
			Method:      r.Method,
			Type:        contentType,
			Time:        spend,
			Size:        uint16(size),
			Status:      uint16(response.StatusCode),
			ReqHeader:   reqHeader,
			ReqTrailer:  reqTrailer,
			ReqCookie:   reqCookie,
			ReqBody:     string(reqBody),
			RespHeader:  respHeader,
			RespTrailer: respTrailer,
			RespCookie:  respCookie,
			RespBody:    respBody,
			RespTls:     respTls,
		}
	}(r, response)
}

func (p *Proxy) handleHttps(w http.ResponseWriter, _ *http.Request) {
	client, server := net.Pipe()
	defer func() {
		_ = client.Close()
	}()

	p.listener.AddConn(server)
	_, _ = w.Write([]byte("HTTP/1.1 200 Connection Established\n\n"))

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	hijack, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	defer func() {
		_ = hijack.Close()
	}()

	var g sync.WaitGroup
	g.Add(2)
	go func() {
		defer g.Done()
		transfer(client, hijack)
	}()
	go func() {
		defer g.Done()
		transfer(hijack, client)
	}()
	g.Wait()
}
func (p *Proxy) start() error {
	return p.srv.ServeTLS(p.listener, "", "")
}
func (p *Proxy) close() (err error) {
	err = p.srv.Close()
	if err != nil {
		return
	}
	err = p.listener.Close()
	if err != nil {
		return
	}
	return
}
func (p *Proxy) forward(w http.ResponseWriter, r *http.Request) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	hijack, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	defer func() {
		_ = hijack.Close()
	}()

	r.Header.Del("Proxy-Connection")
	if r.URL.Host == "" {
		r.URL.Host = r.Host
	}
	if r.URL.Scheme == "" {
		r.URL.Scheme = "https"
	}

	conn, err := new(net.Dialer).Dial("tcp", r.Host)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	if r.Method == "CONNECT" {
		_, err = fmt.Fprint(hijack, "HTTP/1.1 200 Connection established\r\n\r\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	var g sync.WaitGroup
	g.Add(2)
	go func() {
		defer g.Done()
		transfer(conn, hijack)
	}()
	go func() {
		defer g.Done()
		transfer(hijack, conn)
	}()
	g.Wait()
}
func (p *Proxy) SetFilter(filter string) {
	p.filter, _ = regexp.Compile(filter)
}
func (p *Proxy) ClearFilter() {
	p.filter = nil
}
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(p.exclude) > 0 {
		u := r.URL
		if u.Host == "" {
			u.Host = r.Host
		}
		if u.Scheme == "" {
			u.Scheme = "https"
		}
		for _, v := range p.exclude {
			if strings.Contains(u.String(), v) {
				p.forward(w, r)
				return
			}
		}
	}
	if p.filter != nil {
		u := r.URL
		if u.Host == "" {
			u.Host = r.Host
		}
		if u.Scheme == "" {
			u.Scheme = "https"
		}
		if len(p.filter.FindStringIndex(u.String())) == 0 {
			p.forward(w, r)
			return
		}
	}

	if r.Method == http.MethodConnect {
		p.handleHttps(w, r)
	} else {
		if p.replace {
			p.doReplace(w, r)
		} else {
			p.doRequest(w, r)
		}
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer func() {
		_ = destination.Close()
	}()
	defer func() {
		_ = source.Close()
	}()
	_, _ = io.Copy(destination, source)
}

func NewProxy(filter string, proxy string, replace bool) (p *Proxy, err error) {
	p = new(Proxy)
	p.logger = slog.Default()
	p.SetFilter(filter)
	p.exclude = []string{"localhost", "127.0.0.1"}

	if P, err := url.Parse(proxy); err != nil {
		p.proxy = P
	}
	p.replace = replace

	// ca.cert
	block, _ := pem.Decode(static.CaCert)
	if block == nil {
		return
	}
	p.rootCert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return
	}

	// ca.key
	block, _ = pem.Decode(static.CaKey)
	if block == nil {
		return
	}
	if err != nil {
		return
	}
	p.rootKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	// server.key
	p.privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}

	p.listener, _ = NewListener()
	p.srv = &http.Server{
		Handler: p,
		TLSConfig: &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				return p.getCertificate(info.ServerName)
			},
		},
	}
	go func() {
		_ = p.start()
	}()
	return
}
