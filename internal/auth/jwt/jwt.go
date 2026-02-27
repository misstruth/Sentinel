package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// Claims JWT 声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}

var secret = []byte("fo-sentinel-secret-key")

// SetSecret 设置密钥
func SetSecret(s string) {
	secret = []byte(s)
}

// Generate 生成 token
func Generate(userID uint, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Exp:      time.Now().Add(24 * time.Hour).Unix(),
	}
	return encode(claims)
}

// Parse 解析 token
func Parse(token string) (*Claims, error) {
	claims, err := decode(token)
	if err != nil {
		return nil, err
	}
	if claims.Exp < time.Now().Unix() {
		return nil, errors.New("token expired")
	}
	return claims, nil
}

// encode 编码JWT
func encode(c *Claims) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload, _ := json.Marshal(c)
	payloadEnc := base64.RawURLEncoding.EncodeToString(payload)
	sig := sign(header + "." + payloadEnc)
	return header + "." + payloadEnc + "." + sig, nil
}

// decode 解码JWT
func decode(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}
	if sign(parts[0]+"."+parts[1]) != parts[2] {
		return nil, errors.New("invalid signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

func sign(data string) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
