package router

import (
  "github.com/casbin/casbin"
  "github.com/gin-gonic/gin"

  "gitlab.jiangxingai.com/jxserving/components/middleware"
  "gitlab.jiangxingai.com/jxserving/components/model/router/api/v1"
)

func InitRouter(r *gin.Engine) *gin.Engine {
  e, _ := casbin.NewEnforcer("../middleware/casbin/model.conf", "../middleware/casbin/policy.csv")
  r.Use(middleware.Cors())
  r.MaxMultipartMemory = 500
  r.Use(middleware.Jwt())
  r.Use(middleware.Interceptor(e))
  apiV1 := r.Group("/api/v1/model")
  {
    apiV1.GET("/framework/query", v1.GetFramework)
    apiV1.GET("/query", v1.QueryModelsHandler)
    apiV1.GET("/key/query", v1.QueryModelKeyHandler)
    apiV1.POST("/subList/query", v1.QuerySubModelListHandler)
    apiV1.POST("/subModel/update", v1.UpdateSubModelHandler)
    apiV1.POST("/upload", v1.UploadModelHandler)
    apiV1.POST("/delete", v1.DeleteModelHandler)
    apiV1.POST("/offline/download", v1.DownloadModelOffLineHandler)
    apiV1.POST("/publish", v1.PublishModelHandler)
    apiV1.POST("/evaluate", v1.EvaluateModel)
    apiV1.POST("/evaluate/history", v1.HistoryEvaluate)

    apiV1.POST("/jxServing/online/distribute", v1.DownloadModelOnLineHandler)
    apiV1.POST("/jxServing/listModel", v1.QueryTargetNodeListModelHandler)
    apiV1.POST("/jxServing/isAlive", v1.QueryActiveNodeHandler)
    apiV1.POST("/jxServing/machine/info", v1.QueryTargetNodeMachineInfoHandler)
  }
  return r
}