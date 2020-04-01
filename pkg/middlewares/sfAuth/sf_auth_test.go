package sfAuth

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "net/url"
    "strings"
    "testing"
)

func TestPlugin(t *testing.T) {
    req := &http.Request{
        Header: map[string][]string{
            "User-Agent": {"test1"},
            "Referer":    {"test2"},
            "X-Forwarded-Prefix": {"/a/"},
            "Tssotoken" : {"xxxxxxxx"},
        },
        Proto:      "http",
        Host:       "a.com",
        Method:     "POST",
        RemoteAddr: fmt.Sprintf("%s:%d", "a.com", 90),
        URL: &url.URL{
            User: url.UserPassword("admin", "123"),
            Path: "/a/b/?query=test",
        },
        Body:       ioutil.NopCloser(strings.NewReader("I'm body")),
    }
    
    auth, err := NewHandler(nil)
    require.NoError(t, err)
    
    test := httptest.NewRecorder()
    
    auth.AuthType = Open
    auth.ServeHTTP(test, req)
    assert.Equal(t, 200, test.Code, genneralMsg(200, test.Code, test.Body.String()))
    
    auth.AuthType = Login
    auth.ServeHTTP(test, req)
    assert.Equal(t, 200, test.Code, genneralMsg(200, test.Code, test.Body.String()))
    
    auth.AuthType = Authorization
    auth.ServeHTTP(test, req)
    assert.Equal(t, 200, test.Code, genneralMsg(200, test.Code, test.Body.String()))
    
    auth.AuthType = Approval
    auth.ServeHTTP(test, req)
    assert.Equal(t, 200, test.Code, genneralMsg(200, test.Code, test.Body.String()))
    
    auth.AuthType = Record
    auth.ServeHTTP(test, req)
    assert.Equal(t, 200, test.Code, genneralMsg(200, test.Code, test.Body.String()))
}

func genneralMsg(expected int, returnValue int, body string) string{
    return fmt.Sprintf(`
            Expected: %v
            Actual:   %v
            Body: %s`, expected, returnValue, body)
}
