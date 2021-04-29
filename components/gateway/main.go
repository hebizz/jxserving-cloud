package main

import (
  "fmt"

  "github.com/gin-gonic/gin"

  "gitlab.jiangxingai.com/jxserving/components/gateway/pkg/config"
  "gitlab.jiangxingai.com/jxserving/components/gateway/router"
)

func main() {
  gin.SetMode(config.RunMode())
  engine := gin.Default()
  router.InitRouter(engine)
  if err := engine.Run(config.HttpPort()); err != nil {
    fmt.Print(err)
  }
}
