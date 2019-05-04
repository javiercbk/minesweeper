package security

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	contextKey = "jwtUser"
	userID     = "id"
	userName   = "name"
)

// ErrUserNotFound is returned when a jwt token was not found in the request context
var ErrUserNotFound = errors.New("user was not found in the request context")

// JWTUser has all the data that the JWT encodes
type JWTUser struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// JWTMiddlewareFactory creates a JWTMiddleware
func JWTMiddlewareFactory(jwtSecret string) echo.MiddlewareFunc {
	// TODO: make this middleware respond with the api response format
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(jwtSecret),
		ContextKey: contextKey,
	})
}

// JWTEncode encodes a user into a jwt.MapClaims
func JWTEncode(user JWTUser, d time.Duration) jwt.MapClaims {
	claims := jwt.MapClaims{}
	claims[userID] = user.ID
	claims[userName] = user.Name
	// session lasts only 20 minutes
	claims["exp"] = time.Now().Add(d).Unix()
	return claims
}

// JWTDecode attempt to decode a user
func JWTDecode(c echo.Context) (JWTUser, error) {
	var err error
	jwtUser := JWTUser{}
	user := c.Get(contextKey).(*jwt.Token)
	if user == nil {
		err = ErrUserNotFound
	} else {
		claims := user.Claims.(jwt.MapClaims)
		jwtUser.ID = int64(claims[userID].(float64))
		jwtUser.Name = claims[userName].(string)
	}
	return jwtUser, err
}
