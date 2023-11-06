package api

import "encoding/json"

type Info struct {
	Record     bool       `json:"record,omitempty"`
	Proxy      string     `json:"proxy,omitempty"`
	Exclude    []string   `json:"exclude,omitempty"`
	Include    []string   `json:"include,omitempty"`
	LanIp      string     `json:"lan_ip,omitempty"`
	InternetIp string     `json:"internet_ip,omitempty"`
	ProxyPort  int        `json:"proxy_port,omitempty"`
	Replace    [][]string `json:"replace,omitempty"`
}

func (i *Info) String() string {
	info, _ := json.Marshal(i)
	return string(info)
}
