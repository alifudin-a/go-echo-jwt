package auth

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alifudin-a/go-echo-jwt/database/psql"
	"github.com/alifudin-a/go-echo-jwt/helpers"
	"github.com/alifudin-a/go-echo-jwt/models"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

// Login login
func Login(c echo.Context) (err error) {
	username := c.FormValue("username")
	password := c.FormValue("password")

	var u models.Users
	var resp helpers.Response

	if err = c.Bind(&u); err != nil {
		resp.Code = http.StatusUnprocessableEntity
		resp.Message = "Invalid JSON"
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	//db
	db := psql.OpenDB()
	query := `SELECT * FROM users WHERE username = $1;`
	err = db.Get(&u, query, username)
	if err != nil {
		resp.Code = http.StatusNotFound
		resp.Message = "User Not Found"
		return c.JSON(http.StatusNotFound, resp)
	}

	// check user and password
	if username != u.Username || password != u.Password {
		resp.Code = http.StatusUnauthorized
		resp.Message = "Unauthorized"
		return c.JSON(http.StatusUnauthorized, resp)
	}

	//create token
	accessToken := jwt.New(jwt.SigningMethodHS256)
	// ACCESS TOKEN
	// set claims
	accessTokenClaims := accessToken.Claims.(jwt.MapClaims)
	accessTokenClaims["id"] = u.ID
	accessTokenClaims["username"] = u.Username
	accessTokenClaims["password"] = u.Password
	accessTokenClaims["email"] = u.Email
	accessTokenClaims["fullname"] = u.FullName
	accessTokenClaims["expires"] = time.Now().Add(time.Minute * 15).Unix()

	// generate encoded token and send it as response
	encodeAccessToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		log.Println("Error", err)
		return err
	}

	// create a cookie
	cookie := new(http.Cookie)
	cookie.Name = encodeAccessToken
	// cookie.Name = encodeRefreshToken
	cookie.Expires = time.Now().Add(10 * time.Minute)
	// save cookie
	c.SetCookie(cookie)

	resp.Code = http.StatusOK
	resp.Message = "Successfully Create Token!"
	resp.Data = map[string]interface{}{
		"access_token": encodeAccessToken,
		// "refresh_token": encodeRefreshToken,
	}

	return c.JSON(http.StatusOK, resp)
}

// Accessible page
func Accessible(c echo.Context) (err error) {
	return c.String(http.StatusOK, "Accessible")
}

// Restricted page
func Restricted(c echo.Context) error {
	var resp helpers.Response
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["fullname"].(string)

	resp.Code = http.StatusOK
	resp.Message = "Welcome " + name + "!"
	return c.JSON(http.StatusOK, resp)
}
