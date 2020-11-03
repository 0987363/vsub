package middleware

import (
	"github.com/0987363/vsub/models"

	"net/http"

	"github.com/spf13/viper"

	"github.com/0987363/mgo/bson"

	"github.com/gin-gonic/gin"
)

/*
const HeaderAuthentication = "X-Heifeng-Authentication"
const HeaderAccessID = "X-Heifeng-AccessID"
const HeaderSecret = "X-Heifeng-Secret"
*/

func Authenticator() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := GetDB(c)

		token, err := GetToken(c)
		if err != nil {
			writeError(c, http.StatusUnauthorized, models.Error("Load token from cookie failed:", err))
			return
		}
		if !token.Validate(viper.GetString("authentication.secret")) {
			writeError(c, http.StatusUnauthorized, models.Error("Validate token: %+v failed.", token))
			return
		}

		user := models.FindUserByID(db, bson.ObjectIdHex(token.UserID))
		if user == nil {
			writeError(c, http.StatusUnauthorized, models.Error("Could not found user:", token.UserID))
			return
		}

		c.Set(models.MiddwareKeyUserID, token.UserID)
		c.Next()
	}
}

/*
func Authenticator() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := GetLogger(c)
		db := GetDB(c)

		access := c.Request.Header.Get(HeaderAccessID)
		secret := c.Request.Header.Get(HeaderSecret)
		if access == "" || secret == "" {
			logger.Error("Authorization found empty.")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ah := c.Request.Header.Get(HeaderAuthentication)
		if ah == "" {
			ah, _ = url.QueryUnescape(c.Query(HeaderAuthentication))
		}
		if ah == "" {
			logger.Error("Authentication token is nil.")
			writeError(c, http.StatusUnauthorized, "Token is nil.")
			return
		}

		token, err := models.ParseToken(ah)
		if err != nil {
			if logger != nil {
				logger.Infof(
					"Failed to parse authentication token: %s due to: %s.",
					ah,
					err,
				)
			}
			writeError(c, http.StatusUnauthorized, "Token is invalid.")
			return
		}

		if !token.Validate(viper.GetString("authentication.secret")) {
			if logger != nil {
				logger.Errorf("Authentication token: %s is not valid.", ah)
			}
			writeError(c, http.StatusUnauthorized, "Token is in valid.")
			return
		}

		user := models.FindUserByID(db, bson.ObjectIdHex(token.UserID))
		if user == nil {
			logger.Error("User is invalid:", token.UserID)
			writeError(c, http.StatusUnauthorized, "Token is invalid.")
			return
		}

		if user.Token != ah {
			logger.Errorf("Authentication token: %s is expiry. Relogin please", ah)
			writeError(c, http.StatusUnauthorized, "Token is expiry.")
			return
		}

		c.Set(userIDKey, token.UserID)
		c.Next()
	}
}
*/

func GetUserID(c *gin.Context) string {
	if id, ok := c.Get(models.MiddwareKeyUserID); ok {
		return id.(string)
	}

	return ""
}

func writeError(c *gin.Context, code int, err error) {
	logger := GetLogger(c)
	logger.Error(err)

	c.AbortWithStatus(code)
}
