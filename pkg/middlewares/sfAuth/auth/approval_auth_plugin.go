package auth

import (
    "fmt"
    "github.com/containous/traefik/v2/pkg/middlewares/sfAuth"
    "net/http"
)

type ApprovalAuthPlugin struct {

}

func (o ApprovalAuthPlugin) Auth(s sfAuth.SfAuth, rw http.ResponseWriter, req *http.Request) (authSuccess bool, err error){
    if s.AuthType != sfAuth.Approval {
        return false, nil
    }
    
    fmt.Println("approval auth...")
    
    
    return true, nil
}
