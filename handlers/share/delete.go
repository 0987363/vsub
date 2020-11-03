package share

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

	shareID := c.Param("id")
	if !bson.IsObjectIdHex(shareID) {
		logger.Error("Share id invalid:", shareID)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	share, err := models.FindShareByID(db, bson.ObjectIdHex(shareID))
	if err != nil {
		logger.Error("Find share by id failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if share.UserID != bson.ObjectIdHex(middleware.GetUserID(c)) {
		logger.Error("The share is not user's.")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if err = share.Delete(db); err != nil {
		logger.Error("Delete share failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
