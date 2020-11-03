package share

import (
	"github.com/0987363/mgo/bson"
	"github.com/0987363/vsub/middleware"
	"github.com/0987363/vsub/models"
	"github.com/gofrs/uuid"

	"net/http"

	"github.com/gin-gonic/gin"
)

func createValidate(c *gin.Context) (*models.Share, error) {
	var share models.Share
	if err := c.BindJSON(&share); err != nil {
		return nil, models.Error("Unable to parse and decode the request.")
	}

	return &share, nil
}

func Create(c *gin.Context) {
	db := middleware.GetDB(c)
	logger := middleware.GetLogger(c)

	share, err := createValidate(c)
	if err != nil {
		logger.Error("Validate request failed!", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if len(share.NodeID) > 0 {
		res, err := models.ListNodeByFilter(db, bson.M{
			"_id":     bson.M{"$in": share.NodeID},
			"user_id": bson.ObjectIdHex(middleware.GetUserID(c)),
		})
		if err != nil {
			logger.Error("List nodes by filter failed: ", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if len(share.NodeID) != len(res) {
			logger.Error("Some node id is not yours.")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}

	share.ID = bson.NewObjectId()
	share.UserID = bson.ObjectIdHex(middleware.GetUserID(c))
	uid, err := uuid.NewV4()
	if err != nil {
		logger.Error("Generate uuid failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	share.Key = uid.String()

	if err := share.Create(db); err != nil {
		logger.Error("Create share failed: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &share)
}
