package goauth // มีหน้าที่ generate token
import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Vodka479/go-shop-tutorial/config"
	"github.com/Vodka479/go-shop-tutorial/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

type goAuth struct {
	mapClaims *goMapClaims
	cfg       config.IJwtConfig
}

type goAdmin struct {
	*goAuth
}

type goMapClaims struct {
	Claims               *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims                   // รายละเอียด token พื้นฐาน
}

type IGoAuth interface {
	SignToken() string
}

type IGoAdmin interface {
	SignToken() string
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0)) //time.Unix มี type เป็น time.time
}

func (a *goAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims) //Sign token พร้อม payload // SigningMethodHS256 key เดียว decrypt encrypt
	ss, _ := token.SignedString(a.cfg.SecrectKey())
	return ss
}

func (a *goAdmin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims) //Sign token พร้อม payload // SigningMethodHS256 key เดียว decrypt encrypt
	ss, _ := token.SignedString(a.cfg.AdminKey())
	return ss
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*goMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &goMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SecrectKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has been expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*goMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*goMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &goMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.AdminKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has been expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*goMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &goAuth{
		cfg: cfg,
		mapClaims: &goMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "goshop-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

func NewGoAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IGoAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	default:
		return nil, fmt.Errorf("unknow token type")
	}
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IGoAuth {
	return &goAuth{
		cfg: cfg,
		mapClaims: &goMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "goshop-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IGoAuth {
	return &goAuth{
		cfg: cfg,
		mapClaims: &goMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "goshop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newAdminToken(cfg config.IJwtConfig) IGoAuth {
	return &goAdmin{
		goAuth: &goAuth{
			cfg: cfg,
			mapClaims: &goMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "goshop-api",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}
