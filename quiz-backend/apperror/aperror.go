package apperror

// AppError carries everything a handler needs to build a correct,
// consistent response: which HTTP status to use, a stable machine-readable
// code, and a human-readable message. Services return these instead of
// plain errors.New(...) so handlers never have to guess a status from
// message text.

type AppError struct {
	Status  int
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func NotFound(code, message string) *AppError     { return New(404, code, message) }
func BadRequest(code, message string) *AppError   { return New(400, code, message) }
func Unauthorized(code, message string) *AppError { return New(401, code, message) }
func Forbidden(code, message string) *AppError    { return New(403, code, message) }
func Conflict(code, message string) *AppError     { return New(409, code, message) }
func Internal(code, message string) *AppError     { return New(500, code, message) }
