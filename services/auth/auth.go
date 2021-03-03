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

	// create token
	token := jwt.New(jwt.SigningMethodHS256)

	// set claims
	tokenClaims := token.Claims.(jwt.MapClaims)
	tokenClaims["id"] = u.ID
	tokenClaims["username"] = u.Username
	tokenClaims["password"] = u.Password
	tokenClaims["email"] = u.Email
	tokenClaims["fullname"] = u.FullName
	tokenClaims["expired"] = time.Now().Add(time.Minute * 15).Unix()

	// create a cookie
	cookie := new(http.Cookie)
	cookie.Name = u.Username
	cookie.Value = u.Password
	cookie.Expires = time.Now().Add(10 * time.Minute)
	// save cookie
	c.SetCookie(cookie)

	// generate encoded token and send it as response
	encodeToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println("Error", err)
		return err
	}

	resp.Code = http.StatusOK
	resp.Message = "Successfully Create Token!"
	resp.Data = map[string]interface{}{
		"token": encodeToken,
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
