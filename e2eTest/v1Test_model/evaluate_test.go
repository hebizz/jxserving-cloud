package v1Test_model

import (
  "fmt"
  "testing"

  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "gitlab.jiangxingai.com/jxserving/components/model/router/router"
  "gitlab.jiangxingai.com/jxserving/e2eTest"
  "k8s.io/klog"
)

/**
  模型评价
*/
func TestEvaluateHandler(t *testing.T) {
  body := gin.H{
    "type":     1,
    "modelmd5": "91b918ad09ef4b592b77a4febb9b5db7",
    "name":     "yanhuo",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/model/evaluate", reader)
  data, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("upload target %s success", body["name"]), msg)
  klog.Info("result: ", data)
}
