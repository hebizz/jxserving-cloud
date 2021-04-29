package v1Test_model

import (
  "fmt"
  "net/http"
  "net/http/httptest"
  "testing"

  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "gitlab.jiangxingai.com/jxserving/components/model/router/router"
  "gitlab.jiangxingai.com/jxserving/e2eTest"
  log "k8s.io/klog"
)

/**
  上传模型
*/
func TestUploadModelHandler(t *testing.T) {
  extraParams := map[string]string{
    "name":      "sgcc",
    "version":   "10",
    "framework": "tensorflow",
    "notes":     "data test",
    "targetId":  "5dc93268c9ac77cb348f69d4",
    "threshold": "1.02,1.05",
    "mapping":   "aa, bb",
  }
  r := e2eTest.GenerateGetRequestBody()
  r = router.InitRouter(r)
  request, err := e2eTest.NewFileUploadRequest("/api/v1/model/upload", extraParams, "file", "/home/zhouyou/Desktop/test.tar.gz")
  if err != nil {
    t.Error("generate request body error: ", err)
    return
  }
  w := httptest.NewRecorder()
  r.ServeHTTP(w, request)
  switch w.Code {
  case http.StatusOK:
    data, msg := e2eTest.AssertResponseOk(t, w, err)
    assert.Equal(t, fmt.Sprintf("upload model success"), msg)
    log.Info("upload model md5:", data)
  case http.StatusBadRequest:
    e2eTest.AssertResponseError(t, w, "request params error", err)
  case http.StatusForbidden:
    e2eTest.AssertResponseError(t, w, "model already exists", err)
  case http.StatusInternalServerError:
    e2eTest.AssertResponseError(t, w, "xxxxxxxx", err)
  }
}

/**
  删除model
*/
func TestDeleteModelHandler(t *testing.T) {
  body := gin.H{
    "modelMd5": "f5fa83ffb4ca3a70cde3394d2a8fc52d",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/model/delete", reader)
  data, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("delete model success"), msg)
  log.Info("result: ", data)
}

/**
  downLoad model
*/
func TestDownloadModelOffLineHandler(t *testing.T) {
  body := gin.H{
    "modelMd5": "f5fa83ffb4ca3a70cde3394d2a8fc52d",
    "labelMd5": "8fcc094f0f67e4372eec347430faa73b",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/model/offline/download", reader)
  data, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("download model success"), msg)
  log.Info("result: ", data)
}

/**
  query model list
*/
func TestQueryModelListHandler(t *testing.T) {
  r := e2eTest.GenerateGetRequestBody()
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "GET", "/api/v1/model/query", nil)
  data, msg := e2eTest.AssertResponseOk(t, w, nil)
  assert.Equal(t, fmt.Sprintf("query model list success"), msg)
  log.Info("model list: ", data)
}

/**
  query sub model list
*/
func TestQuerySubModelListHandler(t *testing.T) {
  body := gin.H{
    "modelId": "5dc4e7d0edab76c76de79ab6",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/model/subList/query", reader)
  data, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("query sub model list success"), msg)
  log.Info("result: ", data)

}

/**
  update sub model configs
*/
func TestUpdateSubModelHandler(t *testing.T) {
  body := gin.H{
    "modelId":   "5dc4e7d0edab76c76de79ab6",
    "threshold": []float64{1.02, 1.05, 1.34, 1.21},
    "mapping":   []string{"aa, bb,cc,dd"},
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/model/subModel/update", reader)
  data, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("update sub model success"), msg)
  log.Info("result: ", data)
}

func TestQueryModelKeyHandler(t *testing.T) {
  body := gin.H{
    "username": "test",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/model/key/query", reader)
  data, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, "query  model key success", msg)
  log.Info("result: ", data)
}
