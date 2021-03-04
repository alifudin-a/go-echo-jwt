package helpers

// TokenDetails struct for access token and refresh token
type TokenDetails struct {
	AcessToken          string
	RefreshToken        string
	AccessUUID          string
	RefreshUUID         string
	AccessTokenExpires  int64
	RefreshTokenExpires int64
}
