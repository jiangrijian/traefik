package sfAuth

import (
    "errors"
    "fmt"
    "github.com/containous/traefik/v2/pkg/log"
    "net/http"
    "time"
)

type LoginAuthPlugin struct {

}

func GetToken(req *http.Request) (string, error)  {
    token := req.Header.Get("Tssotoken")
    if token == "" {
        token = req.FormValue("Tssotoken")
    }
    if token == "" {
        return token, errors.New("token is null")
    }
    return token, nil
}



func (o LoginAuthPlugin) Auth(authType AuthType, rw http.ResponseWriter, req *http.Request) (authSuccess bool, err error){
    if authType != Login {
        return false, nil
    }
    
    log.Debug(fmt.Sprintf("%s login auth access...", req.Header.Get("X-Forwarded-Prefix")))
    
    token, err := GetToken(req)
    if err != nil {
        return false, err
    }
    
    userInfo, err := GetUserInfoByToken(token)
    if err != nil {
        return false, err
    }
    if userInfo.ExpireTime > time.Now().Unix() {
        rw.WriteHeader(401)
        err := errors.New(fmt.Sprintf("%s token is expire", token))
        rw.Write([]byte(err.Error()))
        fmt.Errorf("token is expire for %s ", token)
        return false, err
    }
    
    log.Debug(fmt.Sprintf("%s login auth access success...", req.Header.Get("X-Forwarded-Prefix")))
    return true, nil
}
