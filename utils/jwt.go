package utils

import (
	"errors"
	"realWorld/global"
	"realWorld/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	signingKey []byte
}

// Jwt`结构体：signingKey`是一个私有属性，用来存储 JWT 签名密钥。
func NewJwt() *Jwt {
	return &Jwt{signingKey: []byte(global.CONFIG.Jwt.Secret)}
	// signingKey 小写的目的是为了防止外界对它进行修改
	// 每次要用到`Jwt`的地方都需要新建一个`Jwt`对象，这些为了防止多线程的竞争问题
}

// CreateCliams 方法，根据传入用户信息，构造一个自定义的 Claims 对象，结合业务需求和标准字段。
func (j *Jwt) CreateClaims(baseClaims model.BaseClaims) model.GoShopClaims {
	claims := model.GoShopClaims{
		BaseClaims: baseClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    global.CONFIG.Jwt.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.CONFIG.Jwt.ExpireTime) * time.Second)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Duration(global.CONFIG.Jwt.NotBefore) * time.Second)),
		},
	}
	return claims
}

// `GenerateToken`方法，根据 Claims 生成签名后的 Token
func (j *Jwt) GenerateToken(claims *model.GoShopClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.signingKey)
}

var (
	TokenExpired     = errors.New("令牌过期，请重新登录")
	TokenNotValidYet = errors.New("令牌尚未生效，请稍后再试")
	TokenMalformed   = errors.New("非法的令牌")
	TokenInvalid     = errors.New("无效令牌")
)

func (j *Jwt) ParseToken(tokenString string) (*model.GoShopClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.GoShopClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.signingKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, TokenMalformed
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, TokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, TokenNotValidYet
		} else {
			return nil, TokenInvalid
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*model.GoShopClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}
