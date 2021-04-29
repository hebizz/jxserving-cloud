package main

import (
  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/assistant/pkg/config"
  "gitlab.jiangxingai.com/jxserving/components/assistant/router/api/router"
  log "k8s.io/klog"
)

func main() {
  gin.SetMode(config.RunMode())
  engine := gin.Default()
  router.InitRouter(engine)
  if err := engine.Run(config.HttpPort()); err != nil {
    log.Info(err)
  }
}
