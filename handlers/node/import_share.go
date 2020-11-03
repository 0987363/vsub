package node

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"strings"

	"github.com/0987363/mgo/bson"
	"github.com/0987363/vsub/middleware"
	"github.com/0987363/vsub/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

func importShareValidate(c *gin.Context) ([]string, error) {
	data, err := c.GetRawData()
	if err != nil {
		return nil, models.Error("Unable to get request data.")
	}
	if string(data[:4]) == "http" {
		data, err = models.LoadShareFromRemote(string(data))
		if err != nil {
			return nil, err
		}
	}

	vs, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	strs := []string{}
	scanner := bufio.NewScanner(bytes.NewBuffer(vs))
	for scanner.Scan() {
		v := strings.Split(scanner.Text(), "://")
		if len(v) != 2 {
			return nil, models.Error("Split request data failed.")
		}
		if v[0] != "vmess" {
			continue
		}

		strs = append(strs, v[1])
	}

	return strs, nil
}

func ImportShare(c *gin.Context) {
	db := middleware.GetDB(c)
	logger := middleware.GetLogger(c)

	ss, err := importShareValidate(c)
	if err != nil {
		logger.Error("Validate request failed!", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	for _, s := range ss {
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
		logger.Infof("Import node %+v: ", v)
	}

	c.Status(http.StatusAccepted)
}
