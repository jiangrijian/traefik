package sfAuth

import (
    "fmt"
    "github.com/containous/traefik/v2/pkg/log"
    "net/http"
)

type RecordAuthPlugin struct {

}

func (o RecordAuthPlugin) Auth(authType AuthType, rw http.ResponseWriter, req *http.Request) (authSuccess bool, err error){
    if authType != Record {
        return false, nil
    }
    
    log.Debug(fmt.Sprintf(fmt.Sprintf("%s record auth access...", req.Header.Get("X-Forwarded-Prefix"))))
    
    log.Debug(fmt.Sprintf(fmt.Sprintf("%s record auth access success...", req.Header.Get("X-Forwarded-Prefix"))))
    return true, nil
}
