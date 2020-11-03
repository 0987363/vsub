package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/0987363/vsub/models"

	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logrus.New()
		log.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: time.RFC3339Nano}

		start := time.Now()

		c.Set(models.MiddwareKeyLogger, log)
		c.Next()

		log.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"user_agent": c.Request.UserAgent(),
			"remote":     c.ClientIP(),
			"status":     c.Writer.Status(),
			"spent":      int(time.Now().Sub(start) / time.Millisecond),
			"user_id":    GetUserID(c),
		}).Infof("Responded %03d in %s", c.Writer.Status(), time.Now().Sub(start))
	}
}

func GetLogger(c *gin.Context) *logrus.Logger {
	if logger, ok := c.Get(models.MiddwareKeyLogger); ok {
		return logger.(*logrus.Logger)
	}

	return nil
}
