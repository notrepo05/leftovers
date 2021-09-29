package nimbus

type nimbusLaunchRequest struct {
	User string            `json:"user"`
	Cmd  []string          `json:"cmd"`
	Env  map[string]string `json:"env"`
}

type nimbusLaunchResponse struct {
	ID      string `json:"id"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type nimbusStatusResponse struct {
	Status   string `json:"status"`
	State    string `json:"state"`
	Reason   string `json:"reason"`
	Message  string `json:"msg"`
	Exitcode int    `json:"exitcode"`
}

type nimbusLogsResponse struct {
	Logs string `json:"logs"`
}
