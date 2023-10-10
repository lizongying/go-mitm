package proxy

import "encoding/json"

type Message struct {
	Id         uint64            `json:"id,omitempty"`
	Method     string            `json:"method,omitempty"`
	Type       string            `json:"type,omitempty"`
	Time       uint16            `json:"time,omitempty"` // ms
	Size       uint16            `json:"size,omitempty"`
	Status     uint16            `json:"status,omitempty"`
	Url        string            `json:"url,omitempty"`
	RemoteAddr string            `json:"remote_addr,omitempty"`
	ReqHeader  map[string]string `json:"req_header,omitempty"`
	ReqCookie  map[string]string `json:"req_cookie,omitempty"`
	ReqTls     map[string]string `json:"req_tls,omitempty"`
	ReqBody    string            `json:"req_body,omitempty"`
	RespHeader map[string]string `json:"resp_header,omitempty"`
	RespCookie map[string]string `json:"resp_cookie,omitempty"`
	RespTls    map[string]string `json:"resp_tls,omitempty"`
	RespBody   string            `json:"resp_body,omitempty"`
}

func (m *Message) String() string {
	message, _ := json.Marshal(m)
	return string(message)
}
