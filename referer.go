package referer

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	Domains      []string
	EmptyReferer bool
	Type         string // white/black
}

func CreateConfig() *Config {
	return &Config{
		Domains:      []string{},
		EmptyReferer: false,
		Type:         "white",
	}
}

type Referer struct {
	next      http.Handler
	name      string
	whiteList bool
	config    *Config
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	return &Referer{
		next:      next,
		name:      name,
		config:    config,
		whiteList: config.Type[0] == 'w' || config.Type[0] == 'W',
	}, nil
}

func (a *Referer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	referer := req.Header.Get("Referer")
	if a.whiteList {
		// 白名单允许空referer，可以访问
		if a.config.EmptyReferer && referer == "" {
			a.next.ServeHTTP(rw, req)
			return
		}
		// 白名单，如果匹配到，则可以访问
		if a.isMatch(referer) {
			a.next.ServeHTTP(rw, req)
			return
		}
		// 其他情况无法访问
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte{})
	} else {
		// 黑名单拒绝空referer，拒绝访问
		if a.config.EmptyReferer && referer == "" {
			rw.WriteHeader(http.StatusForbidden)
			rw.Write([]byte{})
			return
		}

		// 黑名单，如果匹配到，则拒绝访问
		if a.isMatch(referer) {
			rw.WriteHeader(http.StatusForbidden)
			rw.Write([]byte{})
			return
		}

		a.next.ServeHTTP(rw, req)
	}
}

// 检查referer是否匹配
func (a *Referer) isMatch(referer string) bool {
	u, err := url.Parse(referer)
	if err != nil || referer == "" { // 如果referer为空，则取决于设置
		return a.config.EmptyReferer
	}

	for _, item := range a.config.Domains {
		if strings.HasPrefix(item, "*.") {
			if strings.HasSuffix(u.Host, item[1:]) {
				return true
			}
		} else if item == u.Host {
			return true
		}
	}
	return false
}
