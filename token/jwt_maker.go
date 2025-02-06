package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	// 对称加密的共享密钥
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	// 对称算法HS256,此处NewWithClaims的第二个参数是一个接口 Claims
	// Claims接口包含了Valid()方法,用于验证token是否有效,我们在payload.go中实现
	// 此时的jwtToken还没有签名
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// 使用密钥签名token,并返回最终生成的完整JWT token字符串
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	// keyFunc用于提供签名密钥
	// The function receives the parsed, but unverified Token
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法是否有效:SigningMethodHS256是SigningMethodHMAC的一个instance
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	// &Payload{} 表示传递了一个空的 Payload 结构体的指针，用于存储解析后的 Claims 数据
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	// 区分是Token过期还是Token无效
	if err != nil {
		// 类型断言
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// 从token中提取payload
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
