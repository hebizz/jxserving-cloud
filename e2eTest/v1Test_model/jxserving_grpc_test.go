package v1Test_model

import (
  "fmt"
  "testing"

  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "gitlab.jiangxingai.com/jxserving/components/model/router/router"
  "gitlab.jiangxingai.com/jxserving/e2eTest"
  log "k8s.io/klog"
)

/**
  查询backend运行的ai model
*/
func TestQueryTargetNodeListModelHandler(t *testing.T) {
  body := gin.H{
    "targetIp": "10.55.1.180",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/model/jxServing/listModel", reader)
  data, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("query active node list success"), msg)
  log.Info("node list: ", data)
}

/**
   检查backend的状态
 */
func TestQueryActiveNodeHandler(t *testing.T) {
  body := gin.H{
    "targetIp": "10.55.1.180",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/model/jxServing/isAlive", reader)
  data, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("this ip is alive"), msg)
  log.Info("node list: ", data)
}
/**
   检查backend machine status
 */
func TestQueryNodeMachineInfoHandler(t *testing.T) {
  body := gin.H{
    "targetIp": "10.55.1.180",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/model/jxServing/machine/info", reader)
  data, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("query machine info success"), msg)
  log.Info("node machine info: ", data)
}
