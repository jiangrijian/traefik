package auth

import (
    "github.com/containous/traefik/v2/pkg/middlewares/sfAuth"
    "net/http"
)


type Plugin interface {
    Auth(s sfAuth.SfAuth, rw http.ResponseWriter, req *http.Request) (authSuccess bool, err error)
}
