package cam

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ErrResp .
type ErrResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func include(data string, lst []string) bool {
	for _, item := range lst {
		if strings.Contains(data, item) {
			return true
		}
	}
	return false
}

func makeErrResp(w http.ResponseWriter, code int, message string) {
	er := ErrResp{
		Code:    code,
		Message: message,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(er)
}

func jointURL(host, url string) string {
	return host + url
}

func verifyToken(authURL string, token string, url string) bool {
	cli := http.Client{}
	req, err := http.NewRequest("GET", authURL, nil)
	if err != nil {
		return false
	}
	req.Header.Add("token", token)
	req.Header.Add("verify", url)
	resp, err := cli.Do(req)
	if err != nil {
		return false
	}
	if resp.StatusCode == 200 {
		return true
	}
	return false
}
