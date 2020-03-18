package dflerr

import "errors"

const (
	RequestFailure = "request_failure"
	AccessDenied   = "access_denied"
	NotFound       = "not_found"
	Unknown        = "unknown"
)

var (
	ErrNotFound = errors.New(NotFound)
)

type M map[string]interface{}

type E struct {
	Code    string `json:"code"`
	Meta    M      `json:"meta,omitempty"`
	Reasons []E    `json:"reasons,omitempty"`
}

func New(code string, meta M, reasons ...E) E {
	return E{
		Code:    code,
		Meta:    meta,
		Reasons: reasons,
	}
}

func (e E) Error() string {
	return e.Code
}

func Parse(err error) E {
	return E{
		Code: err.Error(),
	}
}
