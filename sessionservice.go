package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"bloom-filter-example/blacklisthandler"
	jwkutil "bloom-filter-example/jwk"
	"bloom-filter-example/token"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

type Session struct {
	SessionId string `json:"session_id"`
}

var bfc *blacklisthandler.Conn

func main() {
	var err error
	bfc, err = blacklisthandler.Connect(os.Getenv("BLOOM_FILTER_SERVER"))
	if err !=nil {
		log.Println("error while connecting: ", err)
	}

	r := setupRoutes()
	_ = http.ListenAndServe(":3000", r)

}

func setupRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Post("/sessions", sessionsPostHandler)
	r.Delete("/sessions", sessionsDeleteHandler)

	r.Get("/keys", keysGetHandler)

	r.Get("/messages", messagesGetHandler)
	return r
}

// --------------- API handlers --------------------------

// this api is secured at the gateway level:  see config/krakend.json
func messagesGetHandler(w http.ResponseWriter, r *http.Request) {
	message := []byte(`{"message":"secure data"}`)
	_, _ = w.Write(message)
}

// Used for krakend plugin for the jwk security: see config/krakend.json
func keysGetHandler(w http.ResponseWriter, r *http.Request) {
	jwkKeys := jwkutil.GetJWKKeys(os.Getenv("TOKEN_SECRET"), os.Getenv("KID"))
	_ = json.NewEncoder(w).Encode(jwkKeys)
}

// Create session -- returns jwt token
func sessionsPostHandler(w http.ResponseWriter, r *http.Request) {

	jwtClaims := []token.JWTClaim{
		token.JWTClaim{ Property: "username", Value: "joe",},
	}

	session := Session{ SessionId: createToken(jwtClaims), }
	_ = json.NewEncoder(w).Encode(session)
}

// secured api - requires authorization token
// this is needed to retrieve the jti key.
func sessionsDeleteHandler(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header["Authorization"]

	if len(authHeader) != 0 {
		data := strings.Split(authHeader[0], " ")
		jwt := data[1]
		blacklistToken(jwt, bfc)
	}
	_, _ = w.Write([]byte(`{}`))

}

// --------------- helper functions --------------------------

// Creates the token
func createToken(jwtClaims []token.JWTClaim) string {

	tokenSecret := os.Getenv("TOKEN_SECRET")
	tokenExpiration, _ := strconv.Atoi( os.Getenv("TOKEN_EXPIRATION_IN_SECS"))
	jwtToken := token.GenerateJWT(jwtClaims, tokenSecret, tokenExpiration)

	return jwtToken
}

func blacklistToken(jwt string, bfc *blacklisthandler.Conn) {

	jti := token.GetTokenID(jwt)
	log.Println("Blacklisting token with id: ", jti)
	bfc.Add(jti)

}

func init() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}
