package middleware

import (
	"fastgin/boost/config"
	"fastgin/common/cache"
	"fastgin/modules/sys/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

func GenerateJWTToken(uid uint64) (string, error) {
	claims := jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Configs.Jwt.Key))
}

func AuthMiddleware() gin.HandlerFunc {
	jwtSecret := []byte(config.Configs.Jwt.Key)
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			tokenString = c.Request.Header.Get("Bearer")
			if tokenString == "" {
				tokenString = c.Param("token")
				if tokenString == "" {
					tokenString = c.Query("token")
					if tokenString == "" {
						tokenString = c.PostForm("token")
					}
				}
			}
		}
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token " + err.Error()})
			c.Abort()
			return
		}

		uidObj := token.Claims.(jwt.MapClaims)["uid"]
		uid := uint64(uidObj.(float64))
		c.Set("uid", uid)
		ctxUser := cache.GetUser(uid) //cache.UserCache.Get(fmt.Sprintf("%v", uid))
		if ctxUser == nil {
			ctxUser, err = service.NewUserService().GetUserById(uid)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Find user info fail  " + err.Error()})
				c.Abort()
				return
			}
			cache.UserCache.Set(fmt.Sprintf("%v", uid), ctxUser, 0)
		}
		c.Set("user", ctxUser)
		c.Next()
	}
}
