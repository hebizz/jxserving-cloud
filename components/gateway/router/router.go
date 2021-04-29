package router

import (
  "github.com/casbin/casbin"
  "github.com/gin-gonic/gin"

  "gitlab.jiangxingai.com/jxserving/components/gateway/api/v1"
  "gitlab.jiangxingai.com/jxserving/components/middleware"
)

func InitRouter(r *gin.Engine) *gin.Engine {
  e, _ := casbin.NewEnforcer("../middleware/casbin/model.conf", "../middleware/casbin/policy.csv")
  r.Use(middleware.Cors())
  r.Use(middleware.Jwt())
  r.Use(middleware.Interceptor(e))
  apiV1 := r.Group("/api/v1/gateway")
  apiV1.GET("/ping", v1.Ping)

  apiV1.POST("/import", v1.Import)
  return r
}
