package sfAuth

import (
    "fmt"
    "net/http"
    "sync"
)

type SfAuth struct {
    SvcPath    string
    AuthPlugin []Plugin
    AuthType   AuthType
    Next       http.Handler
    Mu         sync.Mutex
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
    openAuthPlugin := OpenAuthPlugin{}
    loginAuthPlugin := LoginAuthPlugin{}
    authorizationAuthPlugin := AuthorizationAuthPlugin{}
    approvalAuthPlugin := ApprovalAuthPlugin{}
    recordAuthPlugin := RecordAuthPlugin{}
    
    plugins := []Plugin{
        openAuthPlugin,
        loginAuthPlugin,
        authorizationAuthPlugin,
        approvalAuthPlugin,
        recordAuthPlugin,
    }
    
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
            rw.WriteHeader(500)
            rw.Write([]byte(err.Error()))
            fmt.Errorf("get svc policy error: %s", err)
            return
        }
        s.AuthType = authType
    }
    
    for _, plugin := range s.AuthPlugin {
        authSuccess, err := plugin.Auth(s.AuthType, rw, req)
        if err != nil {
            rw.WriteHeader(401)
            rw.Write([]byte(err.Error()))
            fmt.Errorf("%s exec auth plugin fail: %s", s.SvcPath, err)
            return
        }
        if authSuccess {
            if s.Next != nil {
                s.Next.ServeHTTP(rw, req)
            }
            return
        }
    }
    
    err := fmt.Errorf("%s auth fail", s.SvcPath)
    rw.WriteHeader(401)
    rw.Write([]byte(err.Error()))
    return
}
