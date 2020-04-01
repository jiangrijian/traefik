package auth

import (
    "errors"
    "fmt"
    "github.com/containous/traefik/v2/pkg/middlewares/sfAuth"
    "net/http"
    "time"
)

type AuthorizationAuthPlugin struct {

}

func (o AuthorizationAuthPlugin) Auth(s sfAuth.SfAuth, rw http.ResponseWriter, req *http.Request) (authSuccess bool, err error){
    if s.AuthType != sfAuth.Authorization {
        return false, nil
    }

    fmt.Println("authorization auth...")
    
    token, err := GetToken(rw, req)
    if err != nil {
        return false, err
    }
    
    // TODO: 检查用户是否是授权用户
    userInfo, err := sfAuth.GetUserInfoByToken(token)
    if err != nil {
        return false, err
    }
    if userInfo.ExpireTime > time.Now().Unix() {
        rw.WriteHeader(401)
        err := errors.New(fmt.Sprintln("%s token is expire", token))
        rw.Write([]byte(err.Error()))
        fmt.Errorf("token is expire for %s ", token)
        return false, err
    }
    
    return true, nil
}
