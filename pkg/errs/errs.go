package errs

import (
	"github.com/pkg/errors"
)

var (
	NotFoundErr = errors.New("resource not found")
	ParamsError = errors.New("param err")
)
