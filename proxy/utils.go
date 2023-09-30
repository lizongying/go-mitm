package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

func LanIp() (lanIp string) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			addr, err := iface.Addrs()
			if err != nil {
				fmt.Println(err)
				continue
			}

			for _, addr := range addr {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					lanIp = ipNet.IP.String()
					break
				}
			}
		}
	}
	return
}

func InternetIp() (internetIp string) {
	resp, err := http.Get("https://api64.ipify.org?format=text")
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if err == nil {
		_, _ = fmt.Fscanf(resp.Body, "%s", &internetIp)
	}
	return
}
