package saavuu

import (
	"errors"
)

var ErrJWT error = errors.New("JWT error")
var ErrParm error = errors.New("parameter error")
var ErrInvalidData error = errors.New("invalid data")
var ErrInvalidInput error = errors.New("invalid input")
