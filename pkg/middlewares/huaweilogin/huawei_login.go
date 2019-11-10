package huaweilogin

import (
	"context"
	"fmt"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/middlewares"
	"net/http"
	"net/url"
)

const (
	typeName = "HuaweiLogin"
)

type huaweiLogin struct {
	next     http.Handler
	user     string
	password string
	loginUrl string
	name     string
}

func New(ctx context.Context, next http.Handler, config dynamic.HuaweiLogin, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, typeName)).Debug("Creating huawei-login middleware")
	return &huaweiLogin{
		user:     config.User,
		password: config.Password,
		loginUrl: config.LoginUrl,
		name:     name,
		next:     next,
	}, nil
}

func (s *huaweiLogin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	fmt.Println("huaweiLogin ServerHttp func() start")

	var jsessionIdValue string
	cookies := req.Cookies()
	for _, v := range cookies {
		if v.Name == "JSESSIONID" {
			jsessionIdValue = v.Value
		}
	}

	if jsessionIdValue == "" {
		// 如果cookie中没有JSESSIONID,获取 JSESSIONID ,添加到 header 中
		// 发送登陆请求
		jessesinId := sendLoginToGetSessionId(s)
		if jessesinId != "" {
			cookie := &http.Cookie{Name: "JSESSIONID", Value: jessesinId}
			req.AddCookie(cookie)
		}

	} else {
		fmt.Errorf("jessionId为空")
	}
	s.next.ServeHTTP(rw, req)

}

func sendLoginToGetSessionId(s *huaweiLogin) string {

	// post 表单接口示例
	var jessionId string
	data := make(url.Values)
	data["account"] = []string{s.user}
	data["pwd"] = []string{s.password}
	resp2, _ := http.PostForm(s.loginUrl, data)

	cookies2 := resp2.Cookies()
	for _, v := range cookies2 {
		if v.Name == "JSESSIONID" {
			jessionId = v.Value
		}
	}
	fmt.Println("huaweiLogin getJessionId begin...")
	defer resp2.Body.Close()

	// GET 接口
	/*resp, err := http.Get("http://localhost:12345/")

	var jessionId string
	cookies := resp.Cookies()
	for _, v := range cookies {
		if(v.Name == "JSESSIONID") {
			jessionId = v.Value
		}
	}

	fmt.Println(jessionId)

	if err != nil {
		fmt.Errorf(err.Error())
	}*/

	return jessionId
}
