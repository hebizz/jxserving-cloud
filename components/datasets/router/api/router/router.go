package router

import (
  "github.com/casbin/casbin"
  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/datasets/router/api/v1"
  "gitlab.jiangxingai.com/jxserving/components/middleware"
)

func InitRouter(r *gin.Engine) *gin.Engine {
  e, _ := casbin.NewEnforcer("../middleware/casbin/model.conf", "../middleware/casbin/policy.csv")
  r.Use(middleware.Cors())
  r.Use(middleware.Jwt())
  r.Use(middleware.Interceptor(e))
  apiV1 := r.Group("/api/v1")
  apiV1.POST("/reportData/upload", v1.UploadReportDataHandler)
  apiV1.GET("/dataSet/query", v1.QueryDataSetsHandler)
  apiV1.POST("/dataSet/download", v1.DownloadDataSetHandler)
  apiV1.POST("/dataSet/upload", v1.UploadDataSetsForVocHandler)
  apiV1.POST("/dataSet/delete", v1.DeleteDataSetsHandler)
  apiV1.POST("/dataSet/manage/query", v1.QueryAllDataSetsForReMarkerHandler)
  apiV1.POST("/dataSet/manage/update", v1.UploadHandleMarkerImageHandler)
  return r
}

