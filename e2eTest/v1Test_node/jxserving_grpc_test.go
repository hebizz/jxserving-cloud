package v1Test_node

import (
  "fmt"
  "testing"

  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "gitlab.jiangxingai.com/jxserving/components/node/router/router"
  "gitlab.jiangxingai.com/jxserving/e2eTest"
  log "k8s.io/klog"
)

/**
   检查node status and listStorageBackend
 */
func TestQueryNodeInfoHandler(t *testing.T) {
  body := gin.H{
    "cidr": "10.55.1.180/32",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/node/jxServing/backend/query", reader)
  data, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("query node info success"), msg)
  log.Info("node machine info: ", data)
}


/**
  检查node status and listStorageBackend
*/
func TestQuerySubNodeModelInfoHandler(t *testing.T) {
  r:= e2eTest.GenerateGetRequestBody()
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/node/jxServing/modelInfo/query", nil)
  data, msg := e2eTest.AssertResponseOk(t, w, nil)
  assert.Equal(t, fmt.Sprintf("query model list info success"), msg)
  log.Info("node machine info: ", data)
}
