package huaweilogin

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/middlewares"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
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

	// 如果cookie中没有JSESSIONID,获取 JSESSIONID ,添加到 header 中
	// 发送登陆请求
	jessesinId := sendLoginToGetSessionId(s)
	fmt.Println("after sendLoginToGetSessionId JSESSIONID:" + jessesinId)
	if jessesinId != "" {
		cookie := &http.Cookie{Name: "JSESSIONID", Value: jessesinId}
		req.AddCookie(cookie)
	} else {
		println("JSESSIONID is nil")
	}
	s.next.ServeHTTP(rw, req)

}

func sendLoginToGetSessionId(s *huaweiLogin) string {

	// post 表单接口示例
	fmt.Println("huaweiLogin.url:" + s.loginUrl + ", huaweiLogin.user:" + s.user + ", huaweiLogin.pwd:" + s.password + ", huaweiLogin.name:" + s.name)
	var jessionId string
	data := make(url.Values)
	data["account"] = []string{s.user}
	data["pwd"] = []string{s.password}
	fmt.Println("huaweiLogin getJessionId begin...")

	var resp2 *http.Response

	if strings.HasPrefix(s.loginUrl, "http") {

		fmt.Println("start send http://...")

		resp2, _ = http.PostForm(s.loginUrl, data)

		fmt.Println(resp2.StatusCode, resp2.Status)

	}
	if strings.HasPrefix(s.loginUrl, "https") {
		fmt.Println("start send https://...")

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		cookieJar, _ := cookiejar.New(nil)
		c := &http.Client{
			Jar:       cookieJar,
			Transport: tr,
		}
		resp2, _ = c.PostForm(s.loginUrl, data)
	}
	if resp2 == nil {
		fmt.Println("response is nil")
	} else {
		cookies := resp2.Cookies()
		fmt.Println("response's cookies: ", cookies)
		for _, v := range cookies {
			if v.Name == "JSESSIONID" {
				jessionId = v.Value
			}
		}
		fmt.Println("huaweiLogin getJessionId begin...")
		defer resp2.Body.Close()
	}

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
