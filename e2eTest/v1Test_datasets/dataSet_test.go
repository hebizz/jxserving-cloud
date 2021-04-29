package v1Test_datasets

import (
  "encoding/json"
  "fmt"
  "net/http/httptest"
  "testing"

  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "gitlab.jiangxingai.com/jxserving/components/datasets/router/api/router"
  "gitlab.jiangxingai.com/jxserving/e2eTest"
)

/**
  数据集上传_test
*/
func TestUploadDataSetsForVocHandler(t *testing.T) {
  extraParams := map[string]string{
    "type": "1",
    "name": "sgcc",
  }
  r := e2eTest.GenerateGetRequestBody()
  r = router.InitRouter(r)
  request, err := e2eTest.NewFileUploadRequest("/api/v1/dataSet/upload", extraParams, "file", "/home/zhouyou/Desktop/10g_voc_dataset.zip")
  if err != nil {
    t.Error("generate request body error: ", err)
    return
  }
  w := httptest.NewRecorder()
  r.ServeHTTP(w, request)
  _, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, "upload file success", msg)
}

/**
  数据集下载_test
*/
func TestDownloadDataSetHandler(t *testing.T) {
  body := gin.H{
    "name":      "sgcc",
    "type":      1, // type=1 : voc / type=2: coco
    "downType":  2, // downType=1 : all / downType=2 : apart
    "StartTime": "2019-10-12",
    "EndTime":   "2019-11-15",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/dataSet/download", reader)
  _, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("download %s file success", body["name"]), msg)
}

/**
  查询数据集列表
*/
func TestQueryDataSetsHandler(t *testing.T) {
  type DataSetBody struct {
    Name   string `json:"name"`
    Marker int64  `json:"marker"`
    Total  int64  `json:"total"`
  }
  r := e2eTest.GenerateGetRequestBody()
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "GET", "/api/v1/dataSet/query", nil)
  data, msg := e2eTest.AssertResponseOk(t, w, nil)
  assert.Equal(t, "query dataset list success", msg)
  dataArr := data.([]interface{})
  var dataBody DataSetBody
  for _, item := range dataArr {
    arr, err := json.Marshal(item)
    assert.Nil(t, err)
    err = json.Unmarshal(arr, &dataBody)
    assert.NotEmpty(t, dataBody.Name)
    assert.NotEmpty(t, dataBody.Marker)
    assert.NotEmpty(t, dataBody.Total)
  }
}

/**
  删除数据集 by name
*/
func TestDeleteDataSetsHandler(t *testing.T) {
  body := gin.H{
    "name": "sgcc",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/dataSet/delete", reader)
  _, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("delete dataSets name=%s success", body["name"]), msg)
}
