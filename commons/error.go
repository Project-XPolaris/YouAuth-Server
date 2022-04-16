package commons

const (
	TokenExpire = "1001"
)

type APIError struct {
	Err  error
	Code string
	Desc string
}

func (e *APIError) Error() string {
	return e.Err.Error()
}
