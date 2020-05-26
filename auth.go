package pie

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"strings"
	"time"
)

type JwtClaims struct {
	Hasura HasuraClaims `json:"https://hasura.io/jwt/claims"`
	jwt.StandardClaims
}

type HasuraClaims struct {
	AllowedRoles []string `json:"x-hasura-allowed-roles"`
	DefaultRole  string   `json:"x-hasura-default-role"`
	Id           string   `json:"x-hasura-user-id"`
}

// generate jwt token
func AuthJwt(id, role string, duration time.Duration) (string, error) {
	claims := &JwtClaims{
		Hasura: HasuraClaims{
			AllowedRoles: []string{role},
			DefaultRole:  role,
			Id:           id,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("APP_JWT_KEY")))
}

// parse jwt token to jwt claims
func DecodeJwt(tokenString string) (*JwtClaims, error) {
	claims := &JwtClaims{}

	parts := strings.Split(tokenString, " ")
	tokenString = parts[len(parts)-1]

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(os.Getenv("APP_JWT_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
