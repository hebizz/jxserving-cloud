package v1Test_node

import (
  "encoding/json"
  "testing"

  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "gitlab.jiangxingai.com/jxserving/components/node/router/router"
  "gitlab.jiangxingai.com/jxserving/e2eTest"
)

type NodeInfo struct {
  Id        string `bson:"_id" json:"id"`
  NodeName  string `json:"name" binding:"required" `
  Cidr      string `json:"cidr" binding:"required" `
  PublicKey string `json:"publicKey" binding:"required" `
  Status    string `json:"status"`
  NetSpeed  string `json:"netSpeed"`
}

type NodeList struct {
  Master   NodeInfo   `json:"master"`
  NodeList []NodeInfo `json:"nodeList"`
}

func TestQueryNodeListHandler(t *testing.T) {
  r := e2eTest.GenerateGetRequestBody()
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "GET", "/api/v1/node/query", nil)
  data, msg := e2eTest.AssertResponseOk(t, w, nil)
  assert.Equal(t, "query node list success", msg)
  var dataBody NodeList
  arr, err := json.Marshal(data)
  assert.Nil(t, err)
  err = json.Unmarshal(arr, &dataBody)
  assert.Nil(t, err)
  assert.NotEmpty(t, dataBody.Master.Cidr)
  assert.NotEmpty(t, dataBody.Master.PublicKey)
  assert.NotEmpty(t, dataBody.Master.Status)
  for _, item := range dataBody.NodeList {
    assert.NotEmpty(t, item.NodeName)
    assert.NotEmpty(t, item.Id)
    assert.NotEmpty(t, item.PublicKey)
    assert.NotEmpty(t, item.Cidr)
  }
}

func TestAddNodeHandler(t *testing.T) {
  body := gin.H{
    "name":      "AI-Main-Node",
    "cidr":      "10.55.1.180/32",
    "publicKey": "0nkARa85uY1AB9wQCcdnTyy+QoOjUV6Rxf/kC6slyQo=",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/node/add", reader)
  _, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, "add node list success", msg)
}

func TestDeleteNodeHandler(t *testing.T) {
  body := gin.H{
    "id": "5db69effa0c32ffb7e1e3478",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/node/delete", reader)
  _, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, "delete node success", msg)
}

func TestUpdateNodeInfo(t *testing.T) {
  body := gin.H{
    "id":        "5db69ec2776c00392f7bf627",
    "name":      "AI-Main-Node",
    "cidr":      "9.101.0.3/32",
    "publicKey": "xcsssdd=",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/node/update", reader)
  _, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, "update node success", msg)
}
