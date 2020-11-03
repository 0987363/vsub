package node

import (
	"strings"

	"github.com/0987363/mgo/bson"
	"github.com/0987363/vsub/middleware"
	"github.com/0987363/vsub/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

func importNodeValidate(c *gin.Context) (string, error) {
	data, err := c.GetRawData()
	if err != nil {
		return "", models.Error("Unable to get request data.")
	}
	v := strings.Split(string(data), "://")
	if len(v) != 2 {
		return "", models.Error("Split request data failed.")
	}
	if v[0] != "vmess" {
		return "", models.Error("Recv unknown protocol.")
	}

	return v[1], nil
}

func ImportNode(c *gin.Context) {
	db := middleware.GetDB(c)
	logger := middleware.GetLogger(c)

	s, err := importNodeValidate(c)
	if err != nil {
		logger.Error("Validate request failed!", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	v, err := models.DecodeV2ray(s)
	if err != nil {
		logger.Error("Decode request failed!", err)
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
