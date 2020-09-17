package token

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
)

type JWTClaim struct {
	Property string
	Value    interface{}
}

// Generate JWT token with kid header.
func GenerateJWT(claims []JWTClaim, tokenSecret string, tokenExpiration int) string {

	token := setClaims(jwt.New(), claims, tokenExpiration)
	key:= []byte(tokenSecret)
	// encodedToken, err := jwt.Sign(token, jwa.HS256, key)

	buf, err := json.Marshal(token)

	headers := jws.NewHeaders()
	_ = headers.Set(jws.TypeKey, "JWT")
	_ = headers.Set(jws.KeyIDKey, os.Getenv("KID"))

	encodedToken, err := jws.Sign( buf, jwa.HS256, key, jws.WithHeaders(headers) )

	if err != nil {
		fmt.Println(err)
	}

	return string(encodedToken)

}

// set custom and reserved claims
func setClaims(token jwt.Token, claims []JWTClaim, expiration int) jwt.Token {

	jti := uuid.New().String()

	for _, claim := range claims {
		_ = token.Set(claim.Property, claim.Value)
	}
	_ = token.Set(jwt.IssuedAtKey, time.Now())
	_ = token.Set(jwt.ExpirationKey, time.Now().Local().Add(time.Duration(expiration) * time.Second))
	_ = token.Set(jwt.JwtIDKey, jti)
	return token
}

func GetTokenID(token string) string {
	decodedToken, _ := jwt.Parse(strings.NewReader(token))
	val, _ := decodedToken.Get("jti")
	jti := fmt.Sprintf("%v", val )
	return jti
}

