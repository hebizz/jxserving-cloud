package router

import (
  "github.com/casbin/casbin"
  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/assistant/router/api/v1"
  "gitlab.jiangxingai.com/jxserving/components/middleware"
)

func InitRouter(r *gin.Engine) *gin.Engine {
  e, _ := casbin.NewEnforcer("../middleware/casbin/model.conf", "../middleware/casbin/policy.csv")
  r.Use(middleware.Cors())
  r.Use(middleware.Jwt())
  r.Use(middleware.Interceptor(e))
  apiV1 := r.Group("/api/v1/assistant")
  apiV1.POST("/verifyInfo", v1.PostVerifyInfoHandler)
  apiV1.GET("/verifyInfo", v1.GetVerifyInfoHandler)
  apiV1.POST("/update", v1.UpdateVerifyInfoHandler)
  apiV1.POST("/create", v1.CreateInfoHandler)
  apiV1.POST("/default/update", v1.UploadDataToDataSetHandler)

  apiV1.GET("/project/query", v1.QueryProjectListHandler)
  apiV1.POST("/project/recovery", v1.RecoveryAssistantDataHandler)

  apiV1.POST("/jxserving/switch", v1.JxServingSwitchModelHandler)
  apiV1.POST("/jxserving/detect", v1.JxServingDetectImageHandler)

  apiV1.GET("/ping", v1.Ping)
  return r
}

