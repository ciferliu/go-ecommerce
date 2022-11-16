package common

import "encoding/json"

var (
	FORBIDDEN   *ErrorResponse = &ErrorResponse{Code: 403, Msg: "permission denied"}
	UNAVAILABLE *ErrorResponse = &ErrorResponse{Code: 500, Msg: "service is unavailable for now, please try again later"}
)

type ErrorResponse struct {
	Code int16  `json:"code"`
	Msg  string `json:"msg"`
}

func (r *ErrorResponse) String() string {
	bytes, _ := json.Marshal(r)
	return string(bytes)
}
