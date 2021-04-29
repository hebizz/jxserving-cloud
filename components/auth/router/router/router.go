package router

import (
  "github.com/gin-gonic/gin"

  "gitlab.jiangxingai.com/jxserving/components/auth/router/api"
  "gitlab.jiangxingai.com/jxserving/components/middleware"
)

func InitRouter(r *gin.Engine) {
  r.Use(middleware.Cors())

  r.POST("/login", api.Login)
  r.POST("/register", api.Register)
}
