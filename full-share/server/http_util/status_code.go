package http_util

// StatusCode contains the int and string associated with an http status Status (i.e. 200 OK)
type StatusCode struct {
	Code int
	Msg  string
}

func (s *StatusCode) GetCodeAndMsg() (int, string) {
	return s.Code, s.Msg
}

// These structs should be considered constants. Do not change.
var (
	// ************ 2xx Success ************

	/*
		Status200Ok is the standard response for successful HTTP requests.
	*/
	Status200Ok = StatusCode{200, "Ok"}

	// ************ 4xx Client errors ************

	/*
		Status400BadRequest is the Status when the server cannot or will not process the request due to an apparent client
		error (e.g., malformed request syntax, size too large, invalid request message framing, or deceptive request routing).
	*/
	Status400BadRequest = StatusCode{400, "Bad Request"}

	// ************ 5xx Server errors ************

	/*
		Status500InternalServerError is a generic error message, given when an unexpected condition
		was encountered and no more specific message is suitable.
	*/
	Status500InternalServerError = StatusCode{500, "Internal Server Error"}
)

/*
ErrWithStatus contains a StatusCode and an error. It can be returned from functions where you
need a certain StatusCode associated with an error.
*/
type ErrWithStatus struct {
	Status StatusCode
	Err    error
}

/*
Error is implemented so that ErrWithStatus implements the error interface.
*/
func (e *ErrWithStatus) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}
