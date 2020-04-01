package sfAuth

import (
    "fmt"
    "github.com/containous/traefik/v2/pkg/log"
    "net/http"
)

type OpenAuthPlugin struct {
}

func (o OpenAuthPlugin) Auth(authType AuthType, rw http.ResponseWriter, req *http.Request) (bool, error){
    if authType == Open {
        log.Debug(fmt.Sprintf("%s open access...", req.Header.Get("X-Forwarded-Prefix")))
        return true, nil
    }
    
    return false, nil
}