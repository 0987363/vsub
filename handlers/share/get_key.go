package share

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/0987363/mgo/bson"
	"github.com/0987363/vsub/middleware"
	"github.com/0987363/vsub/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

func GetKey(c *gin.Context) {
	db := middleware.GetDB(c)
	logger := middleware.GetLogger(c)

	key := c.Param("key")
	if len(key) < 10 {
		logger.Error("Key invalid:", key)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	share, err := models.FindShareByKey(db, key)
	if err != nil {
		logger.Error("Find share by key failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if len(share.NodeID) == 0 {
		logger.Info("Req share node empty.")
		c.JSON(http.StatusOK, []int{})
		return
	}

	nodes, err := models.ListNodeByFilter(db, bson.M{"_id": bson.M{"$in": share.NodeID}})
	if err != nil {
		logger.Error("List nodes by filter failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var b bytes.Buffer
	buf := bufio.NewWriter(&b)
	for _, node := range nodes {
		if node.Class == "v2ray" {
			j, _ := json.Marshal(node.V2ray)
			buf.WriteString(fmt.Sprintf("vmess://%s\n", base64.StdEncoding.EncodeToString(j)))
		}
	}
	buf.Flush()

	c.JSON(http.StatusOK, base64.StdEncoding.EncodeToString(b.Bytes()))
}
