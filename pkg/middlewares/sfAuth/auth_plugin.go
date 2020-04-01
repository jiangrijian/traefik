package sfAuth

import (
    "net/http"
)

type AuthType int

const (
    _ AuthType = iota
    Open          = 1   // 开放
    Login         = 2   // 登录
    Authorization = 3   // 授权
    Approval      = 4   // 审批
    Record        = 5   // 备案
)

type Plugin interface {
    Auth(authType AuthType, rw http.ResponseWriter, req *http.Request) (authSuccess bool, err error)
}
