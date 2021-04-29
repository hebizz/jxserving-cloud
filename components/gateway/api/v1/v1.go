package v1

import (
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/gateway/pkg/network"
  log "k8s.io/klog"

  "gitlab.jiangxingai.com/jxserving/components/gateway/pkg/config"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
)

func Ping(c *gin.Context) {
  response.Success(c, "alive(apiV1)", config.Version())
}

/**
  jxserving 告警数据入口
*/
func Import(c *gin.Context) {
  rawDataId, err := network.DatasetStoreRawDataHandler(c)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  var status = "success"
  if err := network.IsAliveHandler(); err == nil {
    err = network.PublishRawDataToAssistantHandler(c, rawDataId)
    if err != nil {
      log.Error("publish data to assistant error: ", err)
      status = "failed"
    }
  } else {
    log.Error("assistant server not connected")
    status = "failed"
  }
  response.Success(c, fmt.Sprintf("push raw data to dataSet success & assistant (%s)", status), nil)
}
