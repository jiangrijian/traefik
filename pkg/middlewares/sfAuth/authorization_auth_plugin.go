package sfAuth

import (
    "errors"
    "fmt"
    "github.com/containous/traefik/v2/pkg/log"
    "net/http"
    "time"
)

type AuthorizationAuthPlugin struct {

}

func (o AuthorizationAuthPlugin) Auth(authType AuthType, rw http.ResponseWriter, req *http.Request) (authSuccess bool, err error){
    if authType != Authorization {
        return false, nil
    }

    log.Debug(fmt.Sprintf("%s authorization auth access...", req.Header.Get("X-Forwarded-Prefix")))
    
    token, err := GetToken(req)
    if err != nil {
        return false, err
    }
    
    // TODO: 检查用户是否是授权用户
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
    
    log.Debug(fmt.Sprintf("%s authorization auth access success...", req.Header.Get("X-Forwarded-Prefix")))
    return true, nil
}
