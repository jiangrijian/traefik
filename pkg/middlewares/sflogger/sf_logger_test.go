package sflogger

import (
    "encoding/json"
    "fmt"
    "github.com/containous/traefik/v2/pkg/types"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "net/url"
    "os"
    "strings"
    "testing"
)

func logWriterTestHandlerFunc(rw http.ResponseWriter, r *http.Request) {
    if _, err := rw.Write([]byte("I'm reBody")); err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }
    
    rw.WriteHeader(200)
}

func TestBodyLog(t *testing.T) {
    testBody(true, "I'm body", t, "1K")
}

func TestNoBodyLog(t *testing.T) {
    testBody(false, "", t, "1K")
}

func TestSomeBodyLog(t *testing.T) {
    testBody(true, "I", t, "1")
}

func GetBodyData(bufferingSize string, bodyEnable bool, bodyMaxSize string, t *testing.T) []byte {
    req := &http.Request{
        Header: map[string][]string{
            "User-Agent": {"test1"},
            "Referer":    {"test2"},
        },
        Proto:      "http",
        Host:       "a.com",
        Method:     "POST",
        RemoteAddr: fmt.Sprintf("%s:%d", "a.com", 90),
        URL: &url.URL{
            User: url.UserPassword("admin", "123"),
            Path: "/a/b/",
        },
        Body:       ioutil.NopCloser(strings.NewReader("I'm body")),
    }
    
    
    filePath := "/tmp/test.log"
    defer os.RemoveAll(filePath)
    config := &types.SfLogger{
        FilePath: filePath,
        BodyEnable: bodyEnable,
        BodyMaxSize: bodyMaxSize,
        BufferingSize: bufferingSize,
        Service: "testService",
    }
    
    logger, err := NewHandler(config, http.HandlerFunc(logWriterTestHandlerFunc))
    require.NoError(t, err)
    defer logger.Close()
    
    logger.ServeHTTP(httptest.NewRecorder(), req)
    
    logData, err := ioutil.ReadFile(filePath)
    require.NoError(t, err)
    
    return logData
    
}

func genneralLog(logData []byte, expected interface{}, returnValue interface{}) string{
    return fmt.Sprintf(`
            Expected: %s
            Actual:   %s
            Log: %s`, expected, returnValue, string(logData))
}

func testBody(bodyEnable bool, expectedBody string, t *testing.T, bodyMaxSize string)  {
    logData := GetBodyData("1K", bodyEnable, bodyMaxSize, t)
    var data = make(map[string]interface{})
    if err := json.Unmarshal(logData, &data); err == nil {
        assert.Equal(t, expectedBody, data["body"],
            genneralLog(logData, "I'm body", data["body"].(string)))
        assert.Equal(t, "a.com", data["hostName"],
            genneralLog(logData, "a.com", data["hostName"].(string)))
        assert.Equal(t, "/a/b/", data["path"],
            genneralLog(logData, "/a/b/", data["path"].(string)))
        assert.Equal(t, "testService", data["service"],
            genneralLog(logData, "testService", data["service"].(string)))
    }
}

func TestHeaders(t *testing.T) {
    logData := GetBodyData("1K", false, "", t)
    
    var data = make(map[string]interface{})
    if err := json.Unmarshal(logData, &data); err == nil {
        headers := data["headers"]
        assert.Equal(t, 2, len(headers.(map[string]interface{})),
            genneralLog(logData, 2, len(headers.(map[string]interface{}))))
    }
}

func TestDirectWriteLog(t *testing.T)  {
    logData := GetBodyData("", false, "", t)
    assert.NotEqual(t, logData, "",
        genneralLog(logData, "", logData))
}