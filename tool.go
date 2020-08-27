package cam

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// ErrResp .
type ErrResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// AuthResponse  response for auth
type AuthResponse struct {
	Code    int    `json:"code"`
	Data    bool   `json:"data"`
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

func verifyToken(authEndpoint string, token string, url string) bool {
	cli := http.Client{}
	req, err := http.NewRequest("GET", authEndpoint, nil)
	if err != nil {
		return false
	}
	req.Header.Add("token", token)
	req.Header.Add("uri", url)
	resp, err := cli.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	ar := new(AuthResponse)
	if err := json.Unmarshal(body, ar); err != nil {
		return false
	}
	return ar.Data
}
