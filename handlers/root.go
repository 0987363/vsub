package handlers

import (
	"github.com/0987363/vsub/middleware"

	"github.com/0987363/vsub/handlers/node"
	"github.com/0987363/vsub/handlers/share"
	"github.com/0987363/vsub/handlers/user"

	"github.com/gin-gonic/gin"
)

var RootMux = gin.New()

func init() {
	gin.SetMode(gin.ReleaseMode)

	RootMux.Use(middleware.Logger())
	RootMux.Use(middleware.Recoverer())
	RootMux.Use(middleware.DBConnector())

	v1Mux := RootMux.Group("/v1")
	{
		{
			v1Mux.POST("/user", user.Create)
			v1Mux.POST("/user/login", user.Login)
		}

		shareMux := v1Mux.Group("/share")
		{
			shareMux.GET("/key/:key", share.GetKey)
		}

		v1Mux.Use(middleware.Authenticator())

		userMux := v1Mux.Group("/user")
		{
			userMux.DELETE("/me", user.Delete)
		}

		nodeMux := v1Mux.Group("/node")
		{
			nodeMux.POST("/v2ray", node.CreateV2ray)
			nodeMux.POST("/node", node.ImportNode)
			nodeMux.POST("/share", node.ImportShare)

			nodeMux.PUT("/v2ray/:id", node.Update)

			nodeMux.DELETE("/id/:id", node.Delete)

			nodeMux.GET("/", node.List)
			nodeMux.GET("/id/:id", node.Get)
		}

		shareMux = v1Mux.Group("/share")
		{
			shareMux.POST("/", share.Create)

			shareMux.PUT("/id/:id", share.Update)

			shareMux.GET("/", share.List)
			shareMux.GET("/id/:id/nodes", share.ListNodes)

			shareMux.DELETE("/id/:id", share.Delete)
		}
	}
}
