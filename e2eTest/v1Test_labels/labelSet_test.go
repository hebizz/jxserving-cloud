package v1Test_labels

import (
  "encoding/json"
  "fmt"
  "net/http"
  "testing"

  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "gitlab.jiangxingai.com/jxserving/components/target/router/api/router"
  "gitlab.jiangxingai.com/jxserving/e2eTest"
  "k8s.io/klog"
)

type Target struct {
  Id      string
  LabelId string   `json:"id"`
  Name    string   `json:"name"`
  D       []string `json:"label"`
}

var labelSetTestId string

/**
  上传 target
*/
func TestUploadLabelSetsHandler(t *testing.T) {
  body := gin.H{
    "name":  "sgcc",
    "label": [3]string{"kk", "xx", "yy"},
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/target/upload", reader)
  data, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("upload target %s success", body["name"]), msg)
  res := data.(map[string]interface{})
  labelSetTestId = res["id"].(string)
  klog.Info("insert target id is:",labelSetTestId)
}

/**
  查询target 列表
*/
func TestQueryLabelSetSHandler(t *testing.T) {
  r := e2eTest.GenerateGetRequestBody()
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "GET", "/api/v1/target/query", nil)
  data, msg := e2eTest.AssertResponseOk(t, w, nil)
  assert.Equal(t, "query target list success", msg)
  dataArr := data.([]interface{})
  var dataBody Target
  for _, item := range dataArr {
    arr, err := json.Marshal(item)
    assert.Nil(t, err)
    err = json.Unmarshal(arr, &dataBody)
    assert.NotEmpty(t, dataBody.Id)
    assert.NotEmpty(t, dataBody.Name)
    assert.NotEmpty(t, dataBody.D)
  }
}

/**
  更新 target
*/
func TestUpdateLabelSetsHandler(t *testing.T) {
  body := gin.H{
    "name":  "yanhuoTest",
    "id":    labelSetTestId,
    "label": [3]string{"kk", "xx", "yy"},
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/target/update", reader)
  _, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("update target %s success", body["name"]), msg)
}

/**
  删除 target
*/
func TestDeleteLabelSetsHandler(t *testing.T) {
  body := gin.H{
    "id": labelSetTestId,
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/target/delete", reader)
  switch w.Code {
  case http.StatusOK:
    _, msg := e2eTest.AssertResponseOk(t, w, err)
    assert.Equal(t, fmt.Sprintf("delete target (%s) success", body["id"]), msg)
  case http.StatusBadRequest:
    e2eTest.AssertResponseError(t, w, "request params error", err)
  case http.StatusForbidden:
    e2eTest.AssertResponseError(t, w, "this label has been associated to model and not be deleted", err)
  case http.StatusNotFound:
    e2eTest.AssertResponseError(t, w, "this target id not exist", err)
  case http.StatusInternalServerError:
    e2eTest.AssertResponseError(t, w, "delete target by id error", err)
  }
}
