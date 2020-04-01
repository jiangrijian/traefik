package sfAuth

import (
    "fmt"
    "github.com/containous/traefik/v2/pkg/log"
    "net/http"
)

type ApprovalAuthPlugin struct {

}

func (o ApprovalAuthPlugin) Auth(authType AuthType, rw http.ResponseWriter, req *http.Request) (authSuccess bool, err error){
    if authType != Approval {
        return false, nil
    }
    
    log.Debug(fmt.Sprintf("%s approval auth access...", req.Header.Get("X-Forwarded-Prefix")))
    
    log.Debug(fmt.Sprintf("%s approval auth access success...", req.Header.Get("X-Forwarded-Prefix")))
    return true, nil
}
