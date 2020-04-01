package auth

import (
    "fmt"
    "github.com/containous/traefik/v2/pkg/middlewares/sfAuth"
    "net/http"
)

type OpenAuthPlugin struct {
}

func (o OpenAuthPlugin) Auth(s sfAuth.SfAuth, rw http.ResponseWriter, req *http.Request) (bool, error){
    if s.AuthType == sfAuth.Open {
        fmt.Println("open auth...")
        return true, nil
    }
    
    return false, nil
}