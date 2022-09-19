package saavuu

import (
	"errors"
	. "net/http"
)

var ErrJWT error = errors.New("JWT error")

func InternalError(w ResponseWriter, err error) {
	w.WriteHeader(StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
