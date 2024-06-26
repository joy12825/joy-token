package gtoken

import (
	"github.com/joy12825/gf/v2/net/ghttp"
	"strings"
)

const FailedAuthCode = 401

type AuthFailed struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}

func (m *GfToken) GetRequestToken(r *ghttp.Request) (token string) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			return
		} else if parts[1] == "" {
			return
		}
		token = parts[1]
	} else {
		authHeader = r.Get("token").String()
		if authHeader == "" {
			return
		}
		token = authHeader
	}
	return
}

func (m *GfToken) GetToken(r *ghttp.Request) (tData *TokenData, err error) {
	token := m.GetRequestToken(r)
	tData, _, err = m.GetTokenData(r.GetCtx(), token)
	return
}

func (m *GfToken) IsLogin(r *ghttp.Request) (b bool) {
	b = true
	urlPath := r.URL.Path
	if !m.AuthPath(urlPath) {
		// 如果不需要认证，继续
		return
	}
	token := m.GetRequestToken(r)
	if len(token) == 0 {
		b = false
		return
	}
	b = m.IsEffective(r.GetCtx(), token)
	return
}
