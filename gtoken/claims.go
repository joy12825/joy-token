package gtoken

import "github.com/golang-jwt/jwt/v4"

// token部分
const (
	ErrorsParseTokenFail    string = "解析Token失败"
	ErrorsTokenInvalid      string = "Token已失效"
	ErrorsTokenNotActiveYet string = "Token尚未激活"
	ErrorsTokenMalFormed    string = "Token格式不正确"

	JwtTokenOK            int = 200100  // 有效
	JwtTokenInvalid       int = -400100 // 无效
	JwtTokenExpired       int = -400101 // 过期
	JwtTokenFormatErrCode int = -400102 // 格式错误
)

type CustomClaims struct {
	Data interface{}
	jwt.RegisteredClaims
}
