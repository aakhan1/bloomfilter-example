package jwkutil

import (
	"encoding/base64"
)

type JWK struct {
	K   string `json:"k"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

func GetJWKKeys(tokenSecret, kid string) interface{} {
	data := []byte(tokenSecret)
	var jwk JWK

	jwk.K = base64.RawURLEncoding.EncodeToString(data)
	jwk.Alg = "HS256"
	jwk.Kid = kid
	jwk.Kty = "oct"

	var jwks JWKS
	jwks.Keys = append(jwks.Keys, jwk)

	return jwks

}
