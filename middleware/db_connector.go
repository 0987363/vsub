package middleware

import (
	"github.com/0987363/mgo"

	"github.com/gin-gonic/gin"
)

const dbKey = "Db"

var db *mgo.Session

func ConnectDB(dataURL string) (err error) {
	db, err = mgo.Dial(dataURL)
	if err == nil {
		db.SetMode(mgo.Monotonic, true)
	}
	return
}

func DBConnector() gin.HandlerFunc {
	return func(c *gin.Context) {
		d := db.Copy()
		d.SetMode(mgo.Monotonic, true)
		c.Set(dbKey, d)
		c.Next()
		defer d.Close()
	}
}

func GetDB(c *gin.Context) *mgo.Session {
	if db, ok := c.Get(dbKey); ok {
		return db.(*mgo.Session)
	}

	return nil
}
