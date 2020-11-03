package user

import (
	"github.com/0987363/vsub/middleware"
	"github.com/0987363/vsub/models"

	"net/http"
	"strings"

	"github.com/0987363/mgo/bson"
	"github.com/gin-gonic/gin"
)

func createValidate(c *gin.Context) (*models.User, error) {
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

// Create a user
func Create(c *gin.Context) {
	db := middleware.GetDB(c)
	logger := middleware.GetLogger(c)

	user, err := createValidate(c)
	if err != nil {
		logger.Error("Validate request failed!", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	origUser, err := models.FindUserByUserName(db, user.UserName)
	if err != nil {
		logger.Error("Find username failed!", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if origUser != nil {
		logger.Error("The user already exists: ", origUser.UserName)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if user.Password == "" || !user.HashPassword() {
		logger.Error("The password is invalid.", user.Password)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user.ID = bson.NewObjectId()
	if err = user.Create(db); err != nil {
		logger.Error("Create user failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}
