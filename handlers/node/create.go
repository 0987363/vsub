package node

import (
	"github.com/0987363/mgo/bson"
	"github.com/0987363/vsub/middleware"
	"github.com/0987363/vsub/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

func createValidate(c *gin.Context) (*models.NodeV2ray, error) {
	var node models.NodeV2ray
	if err := c.BindJSON(&node); err != nil {
		return nil, models.Error("Unable to parse and decode the request.")
	}

	return &node, nil
}

func CreateV2ray(c *gin.Context) {
	db := middleware.GetDB(c)
	logger := middleware.GetLogger(c)

	v, err := createValidate(c)
	if err != nil {
		logger.Error("Validate request failed!", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	node := models.Node{
		ID:     bson.NewObjectId(),
		UserID: bson.ObjectIdHex(middleware.GetUserID(c)),
		Class:  "v2ray",
		V2ray:  v,
	}
	if err = node.Create(db); err != nil {
		logger.Error("Create node failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, &node)
}
