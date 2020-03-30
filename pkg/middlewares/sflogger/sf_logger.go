package sflogger

import (
    "errors"
    "fmt"
    "github.com/containous/traefik/v2/pkg/middlewares/accesslog"
    "github.com/containous/traefik/v2/pkg/types"
    "github.com/sirupsen/logrus"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "strconv"
    "strings"
    "sync"
    "time"
)

type logData struct {
    Header          map[string]interface{}
    Body            string
    Service         string
    Time            time.Time
    HostName        string
    Path            string
}

type handlerParams struct {
    logDataTable *logData
    crr          *accesslog.CaptureRequestReader
    crw          *accesslog.CaptureResponseWriter
}

type Handler struct {
    config         *types.SfLogger
    logger         *logrus.Logger
    file           io.WriteCloser
    mu             sync.Mutex
    logHandlerChan chan handlerParams
    wg             sync.WaitGroup
    next        http.Handler
}

type noopCloser struct {
    *os.File
}

func (h *Handler) SetDefaults() {
    bufferSize := TransUnit(h.config.BufferingSize)
    bodyMaxSize := TransUnit(h.config.BodyMaxSize)
    
    if bufferSize < 0 {
        h.config.BufferingSize = "1K" // 1k
    }
    if bufferSize > 16 * 1024 * 1024 {
        h.config.BufferingSize = "16M" // 16M
    }
    if bodyMaxSize <= 0 {
        h.config.BodyMaxSize = "4K" // 4k
    }
    
    if bodyMaxSize > 16 * 1024 * 1024 {
        h.config.BodyMaxSize = "16M" // 16M
    }
}

func (n noopCloser) Write(p []byte) (int, error) {
    return n.File.Write(p)
}

func (n noopCloser) Close() error {
    // noop
    return nil
}

func (h *Handler) logTheRoundTrip(logDataTable *logData, crr *accesslog.CaptureRequestReader, crw *accesslog.CaptureResponseWriter) {
    fields := make(map[string]interface{})
    headers := logrus.Fields{}
    
    for k, v := range logDataTable.Header {
        headers[k] = v
    }
    fields["headers"] = headers
    fields["body"] = logDataTable.Body
    fields["service"] = logDataTable.Service
    fields["hostName"] = logDataTable.HostName
    fields["path"] = logDataTable.Path
    
    h.mu.Lock()
    defer h.mu.Unlock()
    h.logger.WithFields(fields).Println()
}

func NewHandler(config *types.SfLogger, next http.Handler) (*Handler, error) {
    if config.FilePath == "" {
        return &Handler{}, errors.New("未设置filepath参数！")
    }
    
    
    var file io.WriteCloser = noopCloser{os.Stdout}
    if len(config.FilePath) > 0 {
        f, err := accesslog.OpenAccessLogFile(config.FilePath)
        if err != nil {
            return &Handler{}, fmt.Errorf("error opening access log file: %s", err)
        }
        file = f
    }
    bufferSize := TransUnit(config.BufferingSize)
    logHandlerChan := make(chan handlerParams, bufferSize)
    
    formatter := new(logrus.JSONFormatter)
    
    logger := &logrus.Logger{
        Out:       file,
        Formatter: formatter,
        Hooks:     make(logrus.LevelHooks),
        Level:     logrus.InfoLevel,
    }
    
    logHandler := &Handler{
        config:         config,
        logger:         logger,
        file:           file,
        logHandlerChan: logHandlerChan,
        next:           next,
    }
    
    logHandler.SetDefaults()
    
    if bufferSize > 0 {
        logHandler.wg.Add(1)
        go func() {
            defer logHandler.wg.Done()
            for handlerParams := range logHandler.logHandlerChan {
                logHandler.logTheRoundTrip(handlerParams.logDataTable, handlerParams.crr, handlerParams.crw)
            }
        }()
    }
    
    return logHandler, nil
}

func TransUnit(strNumber string) int64 {
    if strNumber == "" {
        return 0
    }
    
    var number int64
    var err error
    if strings.HasSuffix(strNumber, "K") {
        number, err = strconv.ParseInt(strNumber[: len(strNumber) - 1], 10, 64)
        number = number * 1024
    } else if strings.HasSuffix(strNumber, "M") {
        number, err = strconv.ParseInt(strNumber[: len(strNumber) - 1], 10, 64)
        number = number * 1024 * 1024
    } else if strings.HasSuffix(strNumber, "G") {
        number, err = strconv.ParseInt(strNumber[: len(strNumber) - 1], 10, 64)
        number = number * 1024 * 1024 * 1024
    } else {
        number, err = strconv.ParseInt(strNumber, 10, 64)
    }
    
    if err != nil {
        fmt.Errorf("unit transport error: %s", err)
        return 0
    }
    return number
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    now := time.Now().UTC()
    
    headers := make(map[string]interface{})
    for k, v := range req.Header {
        headers[k] = v
    }
    
    body := ""
    if h.config.BodyEnable && req.Body != nil {
        bodyByte, err := ioutil.ReadAll(req.Body)
        if err != nil {
            fmt.Errorf("error read body: %s", err)
            return
        }
        bodyMaxSize := TransUnit(h.config.BodyMaxSize)
        
        len := int64(len(bodyByte))
        if len > bodyMaxSize {
            len = bodyMaxSize
        }
        body = string(bodyByte[:len])
    }

    logDataTable := &logData{
        Header: headers,
        Time: now.Local(),
        Service: h.config.Service,
        Body: body,
        HostName: req.Host,
        Path: req.URL.Path,
    }
    
    var crr *accesslog.CaptureRequestReader
    if req.Body != nil {
        crr = &accesslog.CaptureRequestReader{Source: req.Body, Count: 0}
    }
    
    crw := &accesslog.CaptureResponseWriter{Rw: rw}
    bufferSize := TransUnit(h.config.BufferingSize)
    
    if bufferSize > 0 {
        h.logHandlerChan <- handlerParams{
            logDataTable: logDataTable,
            crr:          crr,
            crw:          crw,
        }
    } else {
        h.logTheRoundTrip(logDataTable, crr, crw)
    }
    
    if h.next != nil {
        h.next.ServeHTTP(rw, req)
    }
}

func (h *Handler) Close() error {
    close(h.logHandlerChan)
    h.wg.Wait()
    return h.file.Close()
}