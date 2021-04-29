package router

import (
  "github.com/casbin/casbin"
  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/middleware"
  "gitlab.jiangxingai.com/jxserving/components/target/router/api/v1"
)

func InitRouter(r *gin.Engine) *gin.Engine {
  e, _ := casbin.NewEnforcer("../middleware/casbin/model.conf", "../middleware/casbin/policy.csv")
  r.Use(middleware.Cors())
  r.Use(middleware.Jwt())
  r.Use(middleware.Interceptor(e))
  apiV1 := r.Group("/api/v1/target")
  apiV1.POST("/upload", v1.UploadTargetHandler)
  apiV1.POST("/update", v1.UpdateTargetHandler)
  apiV1.POST("/delete", v1.DeleteTargetHandler)
  apiV1.GET("/query", v1.QueryTargetHandler)
  return r
}
