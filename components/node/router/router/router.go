package router

import (
  "net/http"

  "github.com/gin-gonic/gin"

  "gitlab.jiangxingai.com/jxserving/components/node/router/api/v1"
)

func InitRouter(r *gin.Engine) *gin.Engine {

  r.Use(Cors())
  apiV1 := r.Group("/api/v1/node")
  {
    apiV1.POST("/add", v1.AddNodeHandler)
    apiV1.GET("/query", v1.QueryNodeListHandler)
    apiV1.POST("/delete", v1.DeleteNodeHandler)
    apiV1.POST("/update", v1.UpdateNodeInfoHandler)
    apiV1.POST("/jxServing/backend/query", v1.QueryTargetNodeInfoHandler)
    apiV1.POST("/jxServing/modelInfo/query", v1.QuerySubNodeModelInfo)
  }
  return r
}

func Cors() gin.HandlerFunc {
  return func(c *gin.Context) {
    method := c.Request.Method

    c.Header("Access-Control-Allow-Origin", "*")
    c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
    c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
    c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
    c.Header("Access-Control-Allow-Credentials", "true")

    if method == "OPTIONS" {
      c.AbortWithStatus(http.StatusNoContent)
    }
    c.Next()
  }
}
