package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strconv"
	"tiktok_demo/config"
	"tiktok_demo/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Claims struct {
	jwt.RegisteredClaims
}

// GenToken Generate token based on username
func GenToken(userName string) (string, error) {
	var Conf = config.GetConfig()
	var SecretKey = Conf.JWT.Secret
	log.Println("generatorUserName: ", userName)
	u := service.UserService.GetTableUserByUserName(new(service.UserImpl), userName)
	expiresTime := time.Now().Unix() + Conf.OneDayOfHours.OneDayOfHours
	expiresTimeUnix := time.Unix(expiresTime, 0).UTC()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "tiktok_demo",
			ExpiresAt: jwt.NewNumericDate(expiresTimeUnix),
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour * time.Duration(1))), // 过期时间12小时,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        strconv.FormatInt(u.Id, 10),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SecretKey))
	if err == nil {
		log.Println("generate token succeed")
		return signedToken, nil
	} else {
		log.Println("generate token failed")
		return "", err
	}
}

// ParseToken Parse the JWT token
func ParseToken(token string) (*Claims, error) {
	var SecretKey = config.GetConfig().JWT.Secret
	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if jwtToken.Claims.(*Claims).ExpiresAt.Unix() < time.Now().Unix() {
		log.Println("token has expired")
	}
	if err != nil {
		return nil, err
	}
	if claims, ok := jwtToken.Claims.(*Claims); ok && jwtToken.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// VerifyToken Verify the JWT token
func VerifyToken(token string) (string, error) {
	if token == "" {
		return string(rune(0)), nil
	}
	claims, err := ParseToken(token)
	if err != nil {
		return string(rune(0)), err
	}
	return claims.ID, nil
}

// Auth Actions that require login
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		if token == "" {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
		}
		userId, err := VerifyToken(token)
		if err != nil || userId == "0" {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Token Error",
			})
		} else {
			log.Println("token good, userId: ", userId)
		}
		ctx.Set("userId", userId)
		ctx.Next()
	}
}

// AuthWithoutLogin Actions that do not require login
func AuthWithoutLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		if token == "" {
			userId := "0"
			ctx.Set("userId", userId)
		} else {
			userId, err := VerifyToken(token)
			if err != nil {
				ctx.Abort()
				ctx.JSON(http.StatusUnauthorized, Response{
					StatusCode: -1,
					StatusMsg:  "Token Error",
				})
			}
			ctx.Set("userId", userId)
		}
		ctx.Next()
	}
}

// PswEnCode Encode the password of User
func PswEnCode(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
