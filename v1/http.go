package core

type Response struct {
	Status   string      `json:"status"`
	Error    error       `json:"error"`
	Data     interface{} `json:"data"`
	Metadata metadata    `json:"metadata"`
}

type metadata struct {
	SessionToken string `json:"session_token"`
	RequestId    string `json:"request_id"`
	RequestTime  string `json:"request_time"`
}
