package auth

import (
	"net/http"
	"time"

	"github.com/alifudin-a/go-echo-jwt/database/psql"
	"github.com/alifudin-a/go-echo-jwt/helpers"
	"github.com/alifudin-a/go-echo-jwt/models"
	"github.com/labstack/echo/v4"
)

// LoginV2 login handler version 2
func LoginV2(c echo.Context) (err error) {
	username := c.FormValue("username")
	password := c.FormValue("password")

	var u models.Users
	var resp helpers.Response

	// bind struct
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

	createToken, err := helpers.CreateToken(u.ID, u.Username, u.Password, u.Email, u.FullName)
	if err != nil {
		resp.Code = http.StatusUnprocessableEntity
		resp.Message = err.Error()
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	// Save Access token to cookie
	accessTokenCookie := new(http.Cookie)
	accessTokenCookie.Name = createToken.AccessToken
	accessTokenCookie.Expires = time.Now().Add(time.Minute * 15)
	c.SetCookie(accessTokenCookie)

	// Save Refresh token to cookie
	refreshTokenCookie := new(http.Cookie)
	refreshTokenCookie.Name = createToken.RefreshToken
	refreshTokenCookie.Expires = time.Now().Add(time.Hour * 24 * 7)
	c.SetCookie(refreshTokenCookie)

	resp.Code = http.StatusOK
	resp.Message = "Successfully Create Tokens!"
	resp.Data = map[string]interface{}{
		"access_token":  createToken.AccessToken,
		"refresh_token": createToken.RefreshToken,
	}

	return c.JSON(http.StatusOK, resp)
}
