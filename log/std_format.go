package log

import "github.com/json-iterator/go"

type HttpRequestLog struct {
	Ip            string                         `json:"ip"`
	UserId        int                            `json:"userId"`
	Status        int                            `json:"status"`
	Path          string                         `json:"api"`
	Method        string                         `json:"method"`
	Authorization string                         `json:"authorization"`
	AppVersion    string                         `json:"appVersion"`
	AppClient     string                         `json:"appClient"`
	Headers       map[string]jsoniter.RawMessage `json:"headers"`
	Request       map[string]jsoniter.RawMessage `json:"request"`
	Response      map[string]interface{} `json:"response"`
}
