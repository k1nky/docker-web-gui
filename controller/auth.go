package controller

import (
	"log"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

type UserCredentials struct {
	Username string `json:"userName" form:"userName" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type User struct {
	Username string `json:"userName"`
	Role     string `json:"role"`
}

var users map[string]string = map[string]string{
	"rw": "454ee5f",
	"ro": "123Qwer",
}

func NewAuthMiddleware(secret string) *jwt.GinJWTMiddleware {
	mw, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "Dockerboard zone",
		Key:             []byte(secret),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     "id",
		PayloadFunc:     jwtPayload,
		IdentityHandler: jwtIdentityHandler,
		Authenticator: func(ctx *gin.Context) (interface{}, error) {
			var uc UserCredentials
			if err := ctx.ShouldBind(&uc); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			p, ok := users[uc.Username]
			if !ok || p != uc.Password {
				return nil, jwt.ErrFailedAuthentication
			}
			return &User{
				Username: uc.Username,
				Role:     uc.Username,
			}, nil
		},
		Authorizator: func(data interface{}, ctx *gin.Context) bool {
			if strings.Contains(ctx.Request.URL.Path, "command") {
				if v, ok := data.(*User); ok && v.Role == "rw" {
					return true
				}
				return false
			}
			return true
		},
		Unauthorized:  jwtUnauthorizedHandler,
		TokenLookup:   "header: Authorization, param: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	return mw
}

func jwtPayload(data interface{}) jwt.MapClaims {
	if v, ok := data.(*User); ok {
		return jwt.MapClaims{
			"id":   v.Username,
			"role": v.Role,
		}
	}
	return jwt.MapClaims{}
}

func jwtIdentityHandler(ctx *gin.Context) interface{} {
	claims := jwt.ExtractClaims(ctx)
	return &User{
		Username: claims["id"].(string),
		Role:     claims["id"].(string),
	}
}

func jwtUnauthorizedHandler(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{
		"code":    code,
		"message": msg,
	})
}
