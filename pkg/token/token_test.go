package token

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/middlepartedhairstyle/HiWe/pkg/token/jwt_rsa"
	"testing"
	"time"
)

var path = "E:/golang/go/HiWe/config/key"
var pr, pu = "pr.pem", "pu.pem"
var access = "eyJpZCI6MTIzMTI0MzQsIm5hbWUiOiJoZWxsbyIsInR5cGUiOiJhY2Nlc3NfdG9rZW4ifQ.K5zief3bdeXvo8mCIjJ1f2GQ45upTWnmfdhZBl0jiDasAZ8VOQGEiHgS-Zcrl4BkrHVG9CU5YJuDJaFYPhI5Gg"
var refresh = "eyJpZCI6MTIzMTI0MzQsIm5hbWUiOiJoZWxsbyIsInR5cGUiOiJyZWZyZXNoX3Rva2VuIn0.hkmFx9Svk-6GKy2JNl5fm89JyWtvVkjSMP-CGWitvkBFvOBBXX0l1dWissd90loy1DGYOxoN62JL9KsUMR9wEA"

func TestNewToken(t *testing.T) {
	k := jwt_rsa.LoadRSAKeyFile(path, pr, pu)
	token := NewToken(map[string]interface{}{"name": "hello", "id": 12312434}, map[string]interface{}{"name": "hello", "id": 12312434}, k)
	if token.RefreshToken == "" || token.AccessToken == "" {
		t.Errorf("New Token Err")
	}
	fmt.Println(token)
}

func TestNewAccessToken(t *testing.T) {
	k := jwt_rsa.LoadRSAKeyFile(path, pr, pu)
	token := NewAccessToken(map[string]interface{}{"name": "hello", "id": 123}, k)
	if token.AccessToken == "" {
		t.Errorf("New Token Err")
	}
	fmt.Println(token)
}

func TestNewRefreshToken(t *testing.T) {
	k := jwt_rsa.LoadRSAKeyFile(path, pr, pu)
	token := NewRefreshToken(map[string]interface{}{"name": "hello", "id": 123}, k)
	if token.RefreshToken == "" {
		t.Errorf("New Token Err")
	}
	fmt.Println(token)
}

func TestGetToken(t *testing.T) {
	access := "eyJpZCI6MTIzMTI0MzQsIm5hbWUiOiJoZWxsbyIsInR5cGUiOiJhY2Nlc3NfdG9rZW4ifQ.K5zief3bdeXvo8mCIjJ1f2GQ45upTWnmfdhZBl0jiDasAZ8VOQGEiHgS-Zcrl4BkrHVG9CU5YJuDJaFYPhI5Gg"
	refresh := "eyJpZCI6MTIzMTI0MzQsIm5hbWUiOiJoZWxsbyIsInR5cGUiOiJyZWZyZXNoX3Rva2VuIn0.hkmFx9Svk-6GKy2JNl5fm89JyWtvVkjSMP-CGWitvkBFvOBBXX0l1dWissd90loy1DGYOxoN62JL9KsUMR9wEA"
	token := GetToken(access, refresh)
	if token.AccessToken != access || token.RefreshToken != refresh {
		t.Errorf("GetToken Err")
	}
	fmt.Println(token)
}

func TestToken_VerifyAccessToken(t *testing.T) {
	key := jwt_rsa.LoadRSAKeyFile(path, pr, pu)
	access := "eyJpZCI6MTIzMTI0MzQsIm5hbWUiOiJoZWxsbyIsInR5cGUiOiJhY2Nlc3NfdG9rZW4ifQ.K5zief3bdeXvo8mCIjJ1f2GQ45upTWnmfdhZBl0jiDasAZ8VOQGEiHgS-Zcrl4BkrHVG9CU5YJuDJaFYPhI5Gg"
	token := GetToken(access)
	b := token.VerifyAccessToken(key)
	if !b {
		t.Errorf("VerifyAccessToken Err")
	}
	fmt.Println(token)
}

func TestToken_VerifyRefreshToken(t *testing.T) {
	key := jwt_rsa.LoadRSAKeyFile(path, pr, pu)
	refresh := "eyJpZCI6MTIzMTI0MzQsIm5hbWUiOiJoZWxsbyIsInR5cGUiOiJyZWZyZXNoX3Rva2VuIn0.hkmFx9Svk-6GKy2JNl5fm89JyWtvVkjSMP-CGWitvkBFvOBBXX0l1dWissd90loy1DGYOxoN62JL9KsUMR9wEA"
	token := GetToken(refresh)
	b := token.VerifyRefreshToken(key)
	if !b {
		t.Errorf("VerifyRefreshToken Err")
	}
	fmt.Println(token)
}

func TestToken_HVUpdateAccessToken(t *testing.T) {
	access := "eyJpZCI6MTIzMTI0MzQsIm5hbWUiOiJoZWxsbyIsInR5cGUiOiJhY2Nlc3NfdG9rZW4ifQ.K5zief3bdeXvo8mCIjJ1f2GQ45upTWnmfdhZBl0jiDasAZ8VOQGEiHgS-Zcrl4BkrHVG9CU5YJuDJaFYPhI5Gg"
	refresh := "eyJpZCI6MTIzMTI0MzQsIm5hbWUiOiJoZWxsbyIsInR5cGUiOiJyZWZyZXNoX3Rva2VuIn0.hkmFx9Svk-6GKy2JNl5fm89JyWtvVkjSMP-CGWitvkBFvOBBXX0l1dWissd90loy1DGYOxoN62JL9KsUMR9wEA"
	token := GetToken(access, refresh)
	k := jwt_rsa.LoadRSAKeyFile(path, pr, pu)
	payload := NewPayload(map[string]interface{}{"name": "a"}, "iat", "jti")
	fmt.Println(payload)
	token.HVUpdateAccessToken(payload, k)
	if token.AccessToken == access || token.RefreshToken != refresh {
		t.Errorf("HVUpdateAccessToken Err")
	}
	fmt.Println(token.AccessToken)
}

func TestToken_GetPayloadToMap(t *testing.T) {
	token := GetToken(access, refresh)
	payload := token.GetPayloadToMap(ACCESS_TOKEN)
	if payload == nil {
		t.Errorf("GetPayloadToMap Err")
	}
	fmt.Println(payload)
}

func TestToken_GetPayloadToString(t *testing.T) {
	token := GetToken(access, refresh)
	payload := token.GetPayloadToString(ACCESS_TOKEN)
	if payload == "" {
		t.Errorf("GetPayloadToString Err")
	}
	fmt.Println(payload)
}

func TestToken_GetSignature(t *testing.T) {
	token := GetToken(access, refresh)
	signature := token.GetSignature(ACCESS_TOKEN)
	if signature == "" {
		t.Errorf("GetSignature Err")
	}
	fmt.Println(signature)
}

func TestToken_Redis(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "23.95.15.178:6379",
		Password: "luogan",
		DB:       0,
	})
	token := GetToken(access, refresh)
	accessSignature := token.GetSignature(ACCESS_TOKEN)
	refreshSignature := token.GetSignature(REFRESH_TOKEN)
	rdb.Set(context.Background(), "adi123", accessSignature, time.Second*30)
	rdb.Set(context.Background(), "rdi123", refreshSignature, time.Second*30)
	m := map[string]interface{}{
		"access_token":  accessSignature,
		"refresh_token": refreshSignature,
	}
	rdb.HMSet(context.Background(), "userid12314", m)
	if s, _ := rdb.Get(context.Background(), "adi123").Result(); s == "" {
		t.Errorf("Redis Get Err")
	}
}
