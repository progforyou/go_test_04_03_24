package tools

type WSError struct {
	Code int    `json:"-"`
	Err  string `json:"message"`
}

func (e *WSError) Error() string {
	return e.Err
}

func NewWSError(code int, err string) *WSError {
	return &WSError{Code: code, Err: err}
}
