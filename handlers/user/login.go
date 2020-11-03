package user

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/0987363/vsub/middleware"
	"github.com/0987363/vsub/models"
)

func loginValidate(c *gin.Context) (*models.User, error) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		return nil, models.Error("Unable to parse and decode the request.")
	}

	user.UserName = strings.TrimSpace(user.UserName)
	if !models.RegexpUser.MatchString(user.UserName) {
		return nil, models.Error("username must be a valid value.")
	}

	user.Password = strings.TrimSpace(user.Password)
	if !models.RegexpPwd.MatchString(user.Password) {
		return nil, models.Error("password must be a valid value.")
	}

	return &user, nil
}

// Login in with username and sha256 password
func Login(c *gin.Context) {
	db := middleware.GetDB(c)
	logger := middleware.GetLogger(c)

	reqUser, err := loginValidate(c)
	if err != nil {
		logger.Error("Validate request failed!", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := models.FindUserByUserName(db, reqUser.UserName)
	if err != nil {
		logger.Warnf("Could not find user by name: %s, err:%s", reqUser.UserName, err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		logger.Error("Find user failed: ", reqUser.UserName)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := models.PasswordVerify(reqUser.Password, user.HashedPassword); err != nil {
		logger.Errorf("Password is invalid: %s, err:%v", reqUser.UserName, err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	token := models.NewToken(user.ID, viper.GetString("authentication.secret"))
	user.Token = token.String()
	user.Expiry = &token.Expiry

	if err = user.UpdateToken(db); err != nil {
		logger.Errorf("User:%s, Update token failed:%v", reqUser.UserName, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	logger.Infof("User:%s update new token in mongo.", user.UserName)

	middleware.SetToken(c, token)
	c.JSON(http.StatusOK, user)
}
