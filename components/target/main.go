package main

import (
  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/target/pkg/config"
  "gitlab.jiangxingai.com/jxserving/components/target/router/api/router"
  log "k8s.io/klog"
)

func main() {
  log.InitFlags(nil)
  defer log.Flush()
  gin.SetMode(config.RunMode())
  engine := gin.Default()
  router.InitRouter(engine)
  if err := engine.Run(config.HttpPort()); err != nil {
    log.Error(err)
  }
}
