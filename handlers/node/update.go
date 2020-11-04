package node

import (
	"github.com/0987363/mgo/bson"
	"github.com/0987363/vsub/middleware"
	"github.com/0987363/vsub/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

func updateV2rayValidate(c *gin.Context) (*models.Node, error) {
	var node models.Node
	if err := c.BindJSON(&node.V2ray); err != nil {
		return nil, models.Error("Unable to parse and decode the request.")
	}

	nodeID := c.Param("id")
	if !bson.IsObjectIdHex(nodeID) {
		return nil, models.Error("Node id invalid:", nodeID)
	}
	node.ID = bson.ObjectIdHex(nodeID)

	return &node, nil
}

func UpdateV2ray(c *gin.Context) {
	db := middleware.GetDB(c)
	logger := middleware.GetLogger(c)

	node, err := updateV2rayValidate(c)
	if err != nil {
		logger.Error("Validate request failed!", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	node.UserID = bson.ObjectIdHex(middleware.GetUserID(c))
	node.Class = "v2ray"
	if err = node.Update(db); err != nil {
		logger.Error("Update node failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}
