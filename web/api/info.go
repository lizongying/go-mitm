package api

import "encoding/json"

type Info struct {
	Record  bool     `json:"record,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
	Include []string `json:"include,omitempty"`
}

func (i *Info) String() string {
	info, _ := json.Marshal(i)
	return string(info)
}
