package middleware

import (
	"User_Api/handler"
	"User_Api/payloads"
	"User_Api/pkg/models"
	"User_Api/pkg/weberror"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const identityKey = "user"

var MW *jwt.GinJWTMiddleware = nil

func GetJWTMiddle() *jwt.GinJWTMiddleware {
	if MW == nil {
		authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
			Realm:       "lgdSearch",
			Key:         []byte("bu ban meng ma"),
			IdentityKey: identityKey,
			PayloadFunc: func(data interface{}) jwt.MapClaims {
				if v, ok := data.(*models.User); ok {
					return jwt.MapClaims{
						identityKey: v,
					}
				}
				return jwt.MapClaims{}
			},
			Authenticator: func(c *gin.Context) (interface{}, error) {
				var req payloads.LoginReq
				if err := c.ShouldBindJSON(&req); err != nil {
					return "", jwt.ErrMissingLoginValues
				}
				username := req.Username
				password := req.Password
				user, err := handler.QueryUser(0, username)
				if (err == nil) && (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil) {
					return user, nil
				}
				return nil, jwt.ErrFailedAuthentication
			},
			Unauthorized: func(c *gin.Context, code int, message string) {
				c.JSON(code, weberror.Info{Error: message})
			},
			LogoutResponse: func(c *gin.Context, code int) {
				c.JSON(http.StatusNoContent, nil)
			},
		})
		if err != nil {
			panic("JWTMiddle initialization failed")
		}
		MW = authMiddleware
	}
	return MW
}
func ProdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("mytoken") != "\"lgdSearch\"" {
			c.JSON(500, gin.H{"status": fmt.Sprintf("%s", "禁止外部访问!")})
			c.Abort()
		}
		c.Next()
	}
}
