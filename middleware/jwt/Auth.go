package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"tiktok_demo/config"
	"tiktok_demo/service"
	"time"
)

var Conf = config.GetConfig()
var SecretKey = Conf.JWT.Secret

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Claims struct {
	UserId   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

// GenToken Generate token based on username
func GenToken(userName string) (string, error) {
	fmt.Printf("generatetoken: %v\n", userName)
	u := service.UserService.GetTableUserByUserName(new(service.UserImpl), userName)
	expiresTime := time.Now().Unix() + Conf.OneDayOfHours.OneDayOfHours
	fmt.Printf("expiresTime: %v\n", expiresTime)
	claims := Claims{
		UserId:   u.Id,
		UserName: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "tiktok_demo",
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresTime, expiresTime).Local()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SecretKey))
	if err == nil {
		log.Println("generate token success!\n")
		return signedToken, nil
	} else {
		log.Println("generate token fail\n")
		return "", err
	}
}

// ParseToken Parse the JWT token
func ParseToken(token string) (*Claims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := jwtToken.Claims.(*Claims); ok && jwtToken.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// VerifyToken Verify the JWT token
func VerifyToken(token string) (int64, error) {
	if token == "" {
		return int64(0), nil
	}
	claims, err := ParseToken(token)
	if err != nil {
		return int64(0), err
	}
	return claims.UserId, nil
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
		if err != nil || userId == int64(0) {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Token Error",
			})
		} else {
			log.Println("token good")
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
	log.Println("Result: " + sha)
	return sha
}
