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
	json.NewEncoder(w).Encode(er)
	w.WriteHeader(code)
}
