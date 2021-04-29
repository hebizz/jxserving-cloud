package v1

import (
  "github.com/gin-gonic/gin"

  "gitlab.jiangxingai.com/jxserving/components/analyst/pkg/config"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
)

func Ping(c *gin.Context) {
  response.Success(c, "alive(apiV1)", config.Version())
}
