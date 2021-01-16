package apierror

import "encoding/json"

type ErrorMessage struct {
	Message string `json:"message"`
	R       List   `json:"r"`
}

type Error struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type List struct {
	Errors []*Error `json:"errors"`
}

func (l *List) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Errors    []*Error `json:"errors"`
		HasAny    bool     `json:"has_any"`
		HasErrors bool     `json:"has_errors"`
	}{
		Errors:    l.Errors,
		HasErrors: l.HasErrors(),
	})
}

func (l *List) HasErrors() bool {
	return len(l.Errors) > 0
}

func (l *List) AddError(err error) {
	l.Errors = append(l.Errors, &Error{Message: err.Error()})
}
