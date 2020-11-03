package node

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

	nodeID := c.Param("id")
	if !bson.IsObjectIdHex(nodeID) {
		logger.Error("Node id invalid: ", nodeID)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res, err := models.FindNodeByID(db, bson.ObjectIdHex(nodeID))
	if err != nil {
		logger.Error("List node by user id failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if res.UserID != bson.ObjectIdHex(middleware.GetUserID(c)) {
		logger.Error("The node is not user's.")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if err := models.RemoveNodeFromShare(db, res.UserID, res.ID); err != nil {
		logger.Error("Remove node from share failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := res.Delete(db); err != nil {
		logger.Error("Delete node failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
