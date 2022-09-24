package saavuu

import (
	"errors"
)

var ErrJWT error = errors.New("JWT error")
var ErrParm error = errors.New("parameter error")
var ErrInvalidData error = errors.New("invalid data")
var ErrInvalidInput error = errors.New("invalid input")
var ErrInvalidField error = errors.New("invalid field")
var ErrInvalidJwtField error = errors.New("invalid jwt field")
var ErrInvalidJwt error = errors.New("invalid jwt")
var ErrInvalidKey error = errors.New("invalid key")
var ErrInvalidValue error = errors.New("invalid value")
var ErrInvalidType error = errors.New("invalid type")
var ErrInvalidMethod error = errors.New("invalid method")
var ErrInvalidAuth error = errors.New("invalid auth")
var ErrInvalidUserOrPassword error = errors.New("invalid user or password")
