package v1Test_analyst

import (
  "fmt"
  "testing"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "gitlab.jiangxingai.com/jxserving/components/analyst/router"
  "gitlab.jiangxingai.com/jxserving/e2eTest"
)

func TestQuery(t *testing.T) {
  condition := gin.H{
    "key":   "",
    "logic": "",
    "field": "",
    "value": "",
  }
  fmt.Println(condition)
  body := gin.H{
    "timestamp": []int64{time.Now().Unix() - 7200*24, time.Now().Unix()},
    "value":     "reliability",
    "cond":      [0]map[string]interface{}{},
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/analyst/query", reader)
  data, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, "query analyst data success", msg)
  dataArr := data.(map[string]interface{})
  for _, item := range dataArr {
    assert.NotEmpty(t, item)
  }
}
