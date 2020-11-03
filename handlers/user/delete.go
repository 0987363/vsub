package user

import (
	"github.com/0987363/mgo/bson"
	"github.com/0987363/vsub/middleware"
	"github.com/0987363/vsub/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {
	db := middleware.GetDB(c)
	logger := middleware.GetLogger(c)

	if err := models.DeleteShareByUserID(db, bson.ObjectIdHex(middleware.GetUserID(c))); err != nil {
		logger.Error("Delete share failed:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := models.DeleteNodeByUserID(db, bson.ObjectIdHex(middleware.GetUserID(c))); err != nil {
		logger.Error("Delete node failed:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := models.DeleteUserByID(db, bson.ObjectIdHex(middleware.GetUserID(c))); err != nil {
		logger.Error("Delete myself failed:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
