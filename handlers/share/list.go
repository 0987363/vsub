package share

import (
	"github.com/0987363/mgo/bson"
	"github.com/0987363/vsub/middleware"
	"github.com/0987363/vsub/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	db := middleware.GetDB(c)
	logger := middleware.GetLogger(c)

	res, err := models.ListShareByUserID(db, bson.ObjectIdHex(middleware.GetUserID(c)))
	if err != nil {
		logger.Error("List share by user id failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}
