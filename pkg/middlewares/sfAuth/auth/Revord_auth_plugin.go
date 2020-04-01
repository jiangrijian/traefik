package auth

import (
    "fmt"
    "github.com/containous/traefik/v2/pkg/middlewares/sfAuth"
    "net/http"
)

type RecordAuthPlugin struct {

}

func (o RecordAuthPlugin) Auth(s sfAuth.SfAuth, rw http.ResponseWriter, req *http.Request) (authSuccess bool, err error){
    if s.AuthType != sfAuth.Record {
        return false, nil
    }
    
    fmt.Println("record auth...")
    
    
    return true, nil
}
