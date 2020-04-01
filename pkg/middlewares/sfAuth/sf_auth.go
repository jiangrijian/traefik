package sfAuth

import (
    "fmt"
    "github.com/containous/traefik/v2/pkg/middlewares/sfAuth/auth"
    "net/http"
    "sync"
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

type SfAuth struct {
    SvcPath string
    AuthPlugin []auth.Plugin
    AuthType AuthType
    Next    http.Handler
    Mu             sync.Mutex
}

type UserInfo struct {
    Token   string
    UserName string
    ExpireTime int64
    Tenant string
}

func (s *SfAuth)GetSvcPolicyBySvc(svc string) (AuthType, error)  {
    s.Mu.Lock()
    defer s.Mu.Unlock()
    if s.AuthType != 0 {
        return s.AuthType, nil
    }
    
    
    return 1, nil
}

func GetUserInfoByToken(token string) (UserInfo, error) {
    return UserInfo{}, nil
}

func NewHandler(next http.Handler) (*SfAuth, error) {
    openAuthPlugin := auth.OpenAuthPlugin{}
    loginAuthPlugin := auth.LoginAuthPlugin{}
    authorizationAuthPlugin := auth.AuthorizationAuthPlugin{}
    approvalAuthPlugin := auth.ApprovalAuthPlugin{}
    recordAuthPlugin := auth.RecordAuthPlugin{}
    
    plugins := []auth.Plugin{openAuthPlugin, loginAuthPlugin,
        authorizationAuthPlugin, approvalAuthPlugin, recordAuthPlugin}
    
    sfAuth := &SfAuth{
        SvcPath: "",
        AuthType: 0,
        AuthPlugin: plugins,
        Next: next,
    }
    return sfAuth, nil
}

func (s *SfAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    if s.SvcPath == "" {
        s.SvcPath = req.Header.Get("X-Forwarded-Prefix")
    }
    if s.AuthType == 0 {
        authType, err := s.GetSvcPolicyBySvc(s.SvcPath)
        if err != nil {
            fmt.Errorf("get svc policy error: %s", err)
            return
        }
        s.AuthType = authType
    }
    
    for _, plugin := range s.AuthPlugin {
        authSuccess, err := plugin.Auth(*s, rw, req)
        if err != nil {
            fmt.Errorf("%s exec auth plugin error: %s", s.SvcPath, err)
            return
        }
        if authSuccess {
            if s.Next != nil {
                s.Next.ServeHTTP(rw, req)
            }
            return
        }
    }
}
