package auth

import (
    "errors"
    "fmt"
    "github.com/containous/traefik/v2/pkg/middlewares/sfAuth"
    "net/http"
    "time"
)

type LoginAuthPlugin struct {

}

func GetToken(rw http.ResponseWriter, req *http.Request) (token string, err error)  {
    token = req.Header.Get("x-sf-token")
    if token == "" {
        token = req.FormValue("x-sf-token")
    }

    if token == "" {
        rw.WriteHeader(401)
        err := errors.New("x-sf-token is null")
        rw.Write([]byte(err.Error()))
        return "", err
    }
    return token, nil
}

func (o LoginAuthPlugin) Auth(s sfAuth.SfAuth, rw http.ResponseWriter, req *http.Request) (authSuccess bool, err error){
    if s.AuthType != sfAuth.Login {
        return false, nil
    }
    
    fmt.Println("login auth...")
    
    token, err := GetToken(rw, req)
    if err != nil {
        return false, err
    }
    
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
