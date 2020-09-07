package cam

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

var codeMessageMap = map[int]string{
	402: "bad request",
	500: "bad gateway",
	200: "success",
}

// ErrResp .
type ErrResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// AuthUser auth user data
type AuthUser struct {
	Name    string `json:"name"`
	ID      int    `json:"id"`
	IsSuper int    `json:"is_super"`
}

// AuthResponse  response for auth
type AuthResponse struct {
	Code    int      `json:"code"`
	Data    AuthUser `json:"data,omitempty"`
	Message string   `json:"message"`
}

func (ar *AuthResponse) makeData(code int) {
	ar.Code = code
	ar.Message = codeMessageMap[code]
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

func verifyToken(authEndpoint string, token string, url string) *AuthResponse {
	cli := http.Client{}
	ar := new(AuthResponse)
	ar.makeData(200)
	req, err := http.NewRequest("GET", authEndpoint, nil)
	if err != nil {
		ar.makeData(402)
		return ar
	}
	req.Header.Add("token", token)
	req.Header.Add("uri", url)
	resp, err := cli.Do(req)
	if err != nil {
		ar.makeData(402)
		return ar
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		ar.makeData(500)
		return ar
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ar.makeData(500)
		return ar
	}
	if err := json.Unmarshal(body, ar); err != nil {
		ar.makeData(500)
		return ar
	}
	if ar.Code == 0 {
		ar.makeData(resp.StatusCode)
	}
	return ar
}
