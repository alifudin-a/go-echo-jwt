package helpers

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/joho/godotenv/autoload"
	"github.com/twinj/uuid"
)

// TokenDetails struct for access token and refresh token
type TokenDetails struct {
	AccessToken         string
	RefreshToken        string
	AccessUUID          string
	RefreshUUID         string
	AccessTokenExpires  int64
	RefreshTokenExpires int64
}

// CreateToken generate jwt
func CreateToken(id int64, username, password, email, fullname string) (*TokenDetails, error) {
	var err error
	tokeDetails := &TokenDetails{}
	tokeDetails.AccessTokenExpires = time.Now().Add(time.Minute * 15).Unix()
	tokeDetails.AccessUUID = uuid.NewV4().String()
	tokeDetails.RefreshTokenExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokeDetails.RefreshUUID = uuid.NewV4().String()

	// Access Token
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["id"] = id
	accessTokenClaims["username"] = username
	accessTokenClaims["password"] = password
	accessTokenClaims["email"] = email
	accessTokenClaims["fullname"] = fullname
	accessTokenClaims["access_uuid"] = tokeDetails.AccessUUID
	accessTokenClaims["expires"] = tokeDetails.AccessTokenExpires
	acessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	tokeDetails.AccessToken, err = acessToken.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		log.Println("An error occured: ", err)
		return nil, err
	}

	// Refresh Token
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["refresh_uuid"] = uuid.NewV4().String()
	refreshTokenClaims["id"] = id
	refreshTokenClaims["expires"] = tokeDetails.RefreshTokenExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tokeDetails.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("JWT_REFRESH_TOKEN")))
	if err != nil {
		log.Println("An error occured: ", err)
		return nil, err
	}

	return tokeDetails, nil
}
