package token

import (
	"github.com/middlepartedhairstyle/HiWe/pkg/token/encrypt"
	"github.com/middlepartedhairstyle/HiWe/pkg/token/jwt_rsa"
	"math/rand"
	"time"
)

//Token
/*
双token认证
access_token用于资源请求,应设置较短的过期时间(如1-2小时)
refresh_token用于access_token的续签,
access_token过期后使用refresh_token获取新的access_token
refresh_token应设置较长的过期时间(如30天)
*/
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

const (
	ACCESS_TOKEN  = "access_token"
	REFRESH_TOKEN = "refresh_token"
)

//NewToken
/* 该函数用于创建一对token包含accessToken和refreshToken
需要传入accessToken有效载荷部分和refreshToken有效载荷部分，还需传入一对RSA密钥
*/
func NewToken(accessTokenPayload map[string]interface{}, refreshTokenPayload map[string]interface{}, privateKey *encrypt.RSAKey) *Token {
	token := &Token{}
	accessTokenPayload["type"] = ACCESS_TOKEN
	refreshTokenPayload["type"] = REFRESH_TOKEN
	token.AccessToken = jwt_rsa.NewJwtRsa(accessTokenPayload, privateKey).JwtRsaToToken()
	token.RefreshToken = jwt_rsa.NewJwtRsa(refreshTokenPayload, privateKey).JwtRsaToToken()
	return token
}

//NewAccessToken
/* 该函数用于创建单个token包含accessToken
需要传入accessToken有效载荷部分，还需传入一对RSA密钥
*/
func NewAccessToken(accessTokenPayload map[string]interface{}, privateKey *encrypt.RSAKey) *Token {
	token := &Token{}
	accessTokenPayload["type"] = ACCESS_TOKEN
	token.AccessToken = jwt_rsa.NewJwtRsa(accessTokenPayload, privateKey).JwtRsaToToken()
	return token
}

//NewRefreshToken
/* 该函数用于创建单个token包含refreshToken
需要传入refreshToken有效载荷部分，还需传入一对RSA密钥
*/
func NewRefreshToken(refreshTokenPayload map[string]interface{}, privateKey *encrypt.RSAKey) *Token {
	token := &Token{}
	refreshTokenPayload["type"] = REFRESH_TOKEN
	token.RefreshToken = jwt_rsa.NewJwtRsa(refreshTokenPayload, privateKey).JwtRsaToToken()
	return token
}

//NewPayload
/*
构建token的payload字段
输入自定义字段和可选默认字段
"iat":默认当前时间字段
"jti":随机数字段
*/
func NewPayload(payload map[string]interface{}, a ...string) map[string]interface{} {
	for _, val := range a {
		switch val {
		case "iat":
			payload["iat"] = time.Now().Unix()
		case "jti":
			payload["jti"] = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000000)
		default:
		}
	}
	return payload
}

//GetToken
/*获取token，传入为token类型的字符串切片并解析为Token结构体
该函数仅会解析出单个Token结构体，传入切片应注意不能为多个同种类型token
如都为accessToken(该情况仅会获取最后一个access Toke)
*/
func GetToken(a ...string) *Token {
	token := &Token{}
	for _, val := range a {
		jwt := jwt_rsa.TokenToJwtRsa(val)
		payload := jwt.DecodePayloadToMap()
		if payload["type"] == ACCESS_TOKEN {
			token.AccessToken = val
		}
		if payload["type"] == REFRESH_TOKEN {
			token.RefreshToken = val
		}
	}
	return token
}

//VerifyAccessToken
/*对AccessToken进行验证，验证通过返回true反之返回false
 */
func (token *Token) VerifyAccessToken(key *encrypt.RSAKey) bool {
	jwt := jwt_rsa.TokenToJwtRsa(token.AccessToken)
	return jwt.JwtRsaVerify(key)
}

//VerifyRefreshToken
/*对RefreshToken进行验证，验证通过返回true反之返回false
 */
func (token *Token) VerifyRefreshToken(key *encrypt.RSAKey) bool {
	jwt := jwt_rsa.TokenToJwtRsa(token.RefreshToken)
	return jwt.JwtRsaVerify(key)
}

//NVUpdateAccessToken
/*
无验证跟新AccessToken,确保Token结构体中包含的RefreshToken验证通过
输入新的AccessToken荷载字段,以构建新的AccessToken
*/
func (token *Token) NVUpdateAccessToken(payload map[string]interface{}, privateKey *encrypt.RSAKey) {
	t := NewAccessToken(payload, privateKey)
	token.AccessToken = t.AccessToken
}

//NVUpdateRefreshToken
/*
无验证跟新RefreshToken,确保Token结构体中包含的RefreshToken验证通过
输入新的RefreshToken荷载字段,以构建新的RefreshToken
*/
func (token *Token) NVUpdateRefreshToken(payload map[string]interface{}, privateKey *encrypt.RSAKey) {
	t := NewRefreshToken(payload, privateKey)
	token.RefreshToken = t.RefreshToken
}

//HVUpdateAccessToken
/*
有验证跟新AccessToken,确保Token结构体中包含RefreshToken
输入新的AccessToken荷载字段,以构建新的AccessToken
*/
func (token *Token) HVUpdateAccessToken(payload map[string]interface{}, key *encrypt.RSAKey) {
	if token.VerifyRefreshToken(key) {
		t := NewAccessToken(payload, key)
		token.AccessToken = t.AccessToken
	}
}

//HVUpdateRefreshToken
/*
有验证跟新RefreshToken,确保Token结构体中包含RefreshToken
输入新的RefreshToken荷载字段,以构建新的RefreshToken
*/
func (token *Token) HVUpdateRefreshToken(payload map[string]interface{}, key *encrypt.RSAKey) {
	if token.VerifyRefreshToken(key) {
		t := NewRefreshToken(payload, key)
		token.RefreshToken = t.RefreshToken
	}
}

//GetPayloadToMap
/*获取token的Payload并转换为Map类型
 */
func (token *Token) GetPayloadToMap(tokenType string) map[string]interface{} {
	switch tokenType {
	case ACCESS_TOKEN:
		jwt := jwt_rsa.TokenToJwtRsa(token.AccessToken)
		payload := jwt.DecodePayloadToMap()
		return payload
	case REFRESH_TOKEN:
		jwt := jwt_rsa.TokenToJwtRsa(token.RefreshToken)
		payload := jwt.DecodePayloadToMap()
		return payload
	default:
		return nil
	}
}

//GetPayloadToString
/*获取token的Payload并转换为String类型
 */
func (token *Token) GetPayloadToString(tokenType string) string {
	switch tokenType {
	case ACCESS_TOKEN:
		jwt := jwt_rsa.TokenToJwtRsa(token.AccessToken)
		return jwt.Payload
	case REFRESH_TOKEN:
		jwt := jwt_rsa.TokenToJwtRsa(token.RefreshToken)
		return jwt.Payload
	default:
		return ""
	}
}

//GetSignature
/*获取token的Signature
 */
func (token *Token) GetSignature(tokenType string) string {
	switch tokenType {
	case ACCESS_TOKEN:
		jwt := jwt_rsa.TokenToJwtRsa(token.AccessToken)
		return jwt.Signature
	case REFRESH_TOKEN:
		jwt := jwt_rsa.TokenToJwtRsa(token.RefreshToken)
		return jwt.Signature
	default:
		return ""
	}
}
