package common

type LError struct {
	HttpCode  int    `json:"status"`
	ErrorCode int64  `json:"error_code"`
	Msg       string `json:"msg"`
	Err       error  `json:"error"`
}
