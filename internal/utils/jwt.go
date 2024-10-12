package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var signatureKEY []byte

func InitJWT(secretjwt string) {
	signatureKEY = []byte(secretjwt)
}

type DataClaims struct {
	ID         int        `json:"id"`
	RoleID     int        `json:"role_id"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	Name       string     `json:"name"`
	LastAccess *time.Time `json:"last_access"`
}

type Claims struct {
	DataClaims
	jwt.StandardClaims
}

func NewToken(params DataClaims) *Claims {
	return &Claims{
		DataClaims: params,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
}

func (c *Claims) Create() (string, error) {
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signedStr, err := tokens.SignedString(signatureKEY)
	if err != nil {
		return "", fmt.Errorf("failed to SignedString : %v", err)
	}
	return signedStr, nil
}

func CheckToken(token string) (claim *Claims, err error) {
	tokens, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return signatureKEY, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to ParseWithClaims")
	}
	if tokens == nil {
		return nil, fmt.Errorf("token is nil")
	}

	claim, ok := tokens.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("failed to assert claims token")
	}

	return
}
