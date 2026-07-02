package errors

type Kind string

const (
	Invalid      Kind = "invalid"
	NotFound     Kind = "not_found"
	Conflict     Kind = "conflict"
	Forbidden    Kind = "forbidden"
	Unauthorized Kind = "unauthorized"
	Internal     Kind = "internal"
)

type AppError struct {
	Kind    Kind
	Code    string
	Message string
	Details map[string]string
	Err     error
}

func (e AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Code
}
