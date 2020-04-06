package auth

import (
	"net/http"
	"time"

	"github.com/yibozhuang/go-gql-server/internal/orm"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/yibozhuang/go-gql-server/internal/logger"
	"github.com/yibozhuang/go-gql-server/pkg/utils"
)

// Claims JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Begin login with the auth provider
func Begin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// You have to add value context with provider name to get provider name in GetProviderName method
		c.Request = addProviderToContext(c, c.Param("provider"))
		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err != nil {
			gothic.BeginAuthHandler(c.Writer, c.Request)
		} else {
			logger.Infof("user: %#v", gothUser)
		}
	}
}

// Callback callback to complete auth provider flow
func Callback(cfg *utils.ServerConfig, orm *orm.ORM) gin.HandlerFunc {
	return func(c *gin.Context) {
		// You have to add value context with provider name to get provider name in GetProviderName method
		c.Request = addProviderToContext(c, c.Param("provider"))
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		u, err := orm.FindUserByJWT(user.Email, user.Provider, user.UserID)

		if err != nil {
			if u, err = orm.UpsertUserProfile(&user); err != nil {
				logger.Errorf("[Auth.CallBack.UserLoggedIn.UpsertUserProfile.Error]: %q", err)
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}

		logger.Info("[Auth.CallBack.UserLoggedIn]: ", u.ID)
		jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod(cfg.JWT.Algorithm), Claims{
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				Id:        user.UserID,
				Issuer:    user.Provider,
				IssuedAt:  time.Now().UTC().Unix(),
				NotBefore: time.Now().UTC().Unix(),
				ExpiresAt: user.ExpiresAt.UTC().Unix(),
			},
		})

		token, err := jwtToken.SignedString([]byte(cfg.JWT.Secret))
		if err != nil {
			logger.Error("[Auth.Callback.JWT] error: ", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		logger.Info("token: ", token)
		json := gin.H{
			"type":          "Bearer",
			"token":         token,
			"refresh_token": user.RefreshToken,
		}
		c.JSON(http.StatusOK, json)
	}
}

// Logout logs out of the auth provider
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = addProviderToContext(c, c.Param("provider"))
		gothic.Logout(c.Writer, c.Request)
		c.Writer.Header().Set("Location", "/")
		c.Writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}
