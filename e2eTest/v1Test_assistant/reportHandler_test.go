package v1Test_assistant

import (
  "encoding/json"
  "testing"

  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "gitlab.jiangxingai.com/jxserving/components/assistant/router/api/router"
  "gitlab.jiangxingai.com/jxserving/e2eTest"
)

type VerifyBody struct {
  ClusterId     string          `json:"cluster_id" binding:"required" `
  Title         string          `json:"title" binding:"required"`
  AlertType     string          `json:"alert_type" binding:"required"`
  Description   string          `json:"description" binding:"required"`
  AlertPosition []AlertPosition `json:"alert_position" binding:"required"`
  EventId       string          `json:"event_id" binding:"required"`
  DeviceId      string          `json:"device_id" binding:"required"`
  ImagePath     string          `json:"image_path" binding:"required"`
  Image         string          `json:"image" binding:"required"`
  Time          int64           `json:"time"  binding:"required"`
  CreatedTime   int64           `json:"created_time"  binding:"required"`
  AppInfo       AppInfo         `json:"app_info" binding:"required"`
  Reliability   float64         `json:"reliability" binding:"required"`
  DataID        string          `json:"data_id" binding:"required"`
}

type AlertPosition struct {
  LeftX  string `json:"left_x" binding:"required"`
  LeftY  string `json:"left_y" binding:"required"`
  RightX string `json:"right_x" binding:"required"`
  RightY string `json:"right_y" binding:"required"`
}

type AppInfo struct {
  AppName    string `json:"app_name" binging:"required"`
  AppVersion string `json:"app_version" binding:"required"`
}

/**
  test data
*/
var position = map[string]string{
  "left_x":  "143",
  "left_y":  "609",
  "right_x": "503",
  "right_y": "674",
}
var appInfo = map[string]string{
  "app_name":    "通道可视化智能报警程序",
  "app_version": "1.0.0",
}
var vb = gin.H{
  "cluster_id":     "sgcc",
  "title":          "test 数据",
  "alert_type":     "yanhuo",
  "description":    "reliability 0.69938915371894836",
  "reliability":    0.29938915371894836,
  "alert_position": []map[string]string{position},
  "app_info":       appInfo,
  "event_id":       "5fc8752c-c599-11e9-94fe-0242ac1100012",
  "app_id":         "5fc8752c-c599-11e9-94fe-0242acxxxx5",
  "device_id":      "5d4a9c1d97e91b6667bcec93",
  "image_path":     "media/5d4a9c1d97e91b6667bcec93/image/2019-08-23/11-30-00.jpg",
  "image":          "/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCABLAHkDASIAAhEBAxEB/8QAGwAAAgMBAQEAAAAAAAAAAAAAAQIDBAUABgf/xAAuEAACAgIAAwgBAwUBAAAAAAAAAQIDBBEFEjETISJBUWFxgTMGIzIUUnLB0fD/xAAZAQADAQEBAAAAAAAAAAAAAAAAAQIDBAX/xAAfEQACAgIDAQEBAAAAAAAAAAAAAQIRAyESIjEEMkH/2gAMAwEAAhEDEQA/APOpDqJyQ6Wz3DzgKLfRDcrHgmh9bYrAjUSRIZIbQrACQyOSCABCBBQxDoZCIdAAyHQqGQgJEMIhgAwFFskiipxDNWBhTyHBz5fJepJi8c4JmRjzX2YNj7uW+PNDf+S/2ZzzwhLjJmsccpK0W0OkMobTlBqyG9KcHuL+GPGDet9y9iuSe0RTXommMovzIsK63Ix3OyMHFSahKK02k31XQv49Mrfxw5pPoiY5FJWVKDTorJHcppSwLVD8T3v0K1lEqpalFplKSZLTK3KFIkcQaKEBDIGhkAhkMgIZIBjIYCQ+gA8/dRVk0TpuW65rUlvRVjwPA7Ls51O3xcylN+JfaLyHiKWOMnbRSlJaTKlPDsrF2uHZs6k0/wBuxc0f/fQjyeL8Pbnl4fa1+dtEuZL68jUi2vJkXGMq+PBslVTcJcj1qKbft3nNPAoq4aNo5W3Uth4NlYt9U6o2NdklKTaa1tpa+dnosPKp5tqUW1/b5nzf9LQlXxOdWS73ZZHnXJY+Vx8+b1XTr5nrYcOxu3jONt8OVNKKn4fv1MsXN49lZOKmenlm+/Xyl/0p5d1N3R+L3MeX9TU2o2dpHfdt94iypRbV1bRrGo+6IdvwsTjp93QjDG6Fi8M18MPz3HQmmYtUBIZI5IdIoRyQyRyQyQgDFDgQRAYKJYpESJYstgWIyTXeYnHuGznjX5OPZbOxpvs5SXLFa6pexrpgsjGyqdc1uMotNeq0ZyipKi4yp2eV/TtU3xpSrs/bVSlNtd79NHsk+8weCx5cmT5FFOmKS3t9TcRl8y6F533Jdh2mteQiCjYyBKiuS/jr4I+wsj+Oz6ZOgoTihqTIO0urXir37xJYZFcur1rrslijpVwn/KKf0TUl4x2n6h4aktxaa9UNorvFS/HOUH8hbyaouTSsS9Oo+TXqCl/CwkEr4eS8qpzlRbS09asjrfuWNDTtWiWqMBEiEQ6NGIdDIVDw6iGY/CtRz1F2OU3S1rfkpG4jG4aks1JJLumunubaOb5n0Ns/6CkNoCGNzEKGQqHQAFIdICGQgCECCgAKDoC6sYQH/9k=",
  "time":           1568892223,
  "create_time":    1567892342,
  "data_id" +
    "" +
    "" +
    "":        "5d4a9c1d97e91b6667bcec93",
}

/**
  后台上传告警数据
*/
func TestPostVerifyInfo(t *testing.T) {
  r, reader, err := e2eTest.GeneratePostRequestBody(vb)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/assistant/verifyInfo", reader)
  _, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, "save reportData success", msg)
}

/**
  前端获取告警数据
*/
func TestGetVerifyInfo(t *testing.T) {
  r := e2eTest.GenerateGetRequestBody()
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "GET", "/api/v1/assistant/verifyInfo", nil)
  data, msg := e2eTest.AssertResponseOk(t, w, nil)
  if data == "" {
    assert.Equal(t, "no new data", msg)
    return
  }
  dataArr := data.([]interface{})
  var dataBody VerifyBody
  for _, item := range dataArr {
    arr, err := json.Marshal(item)
    assert.Nil(t, err)
    err = json.Unmarshal(arr, &dataBody)
    assert.NotEmpty(t, dataBody.AlertPosition)
    assert.NotEmpty(t, dataBody.ClusterId)
    assert.NotEmpty(t, dataBody.Description)
    assert.NotEmpty(t, dataBody.Image)
    assert.NotEmpty(t, dataBody.DeviceId)
    assert.NotEmpty(t, dataBody.CreatedTime)
    assert.NotEmpty(t, dataBody.Time)
  }
}

/**
  更新告警数据
*/
func TestUpdateVerifyInfo(t *testing.T) {
  body := gin.H{
    "event_id": "5fc8752c-c599-11e9-94fe-0242ac110012",
    "ignore":   false,
    "override": []map[string]interface{}{vb},
    "project":  "NanjingCheckAI",
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/assistant/update", reader)
  _, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, "update info success", msg)
}

/**
  创建告警信息
*/
func TestCreateInfo(t *testing.T) {
  body := gin.H{
    "override": []map[string]interface{}{vb},
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/assistant/create", reader)
  _, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, "createInfo success", msg)
}

/**
  后台上传告警数据
*/
func TestUploadDataToDataSetHandler(t *testing.T) {
  r, reader, err := e2eTest.GeneratePostRequestBody(vb)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/assistant/default/update", reader)
  _, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, "save reportData success", msg)
}

/**
  查询干预终端列表
*/
func TestQueryProjectListHandler(t *testing.T) {
  r := e2eTest.GenerateGetRequestBody()
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "GET", "/api/v1/assistant/project/query", nil)
  _, msg := e2eTest.AssertResponseOk(t, w, nil)
  assert.Equal(t, "query project list success", msg)
}

/**
  恢复告警数据
*/
func TestRecoveryAssistantDataHandler(t *testing.T) {
  body := gin.H{
    "name":       "NanjingCheckAI",
    "startIndex": 1,
    "offset":     10,
    "alertTime":  1573543473,
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/assistant/project/recovery", reader)
  _, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, "recovery reportData list success", msg)
}
