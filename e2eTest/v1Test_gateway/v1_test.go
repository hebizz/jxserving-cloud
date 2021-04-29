package v1Test_gateway

import (
  "testing"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "gitlab.jiangxingai.com/jxserving/components/gateway/router"
  "gitlab.jiangxingai.com/jxserving/e2eTest"
)

func TestImport(t *testing.T) {
  position := map[string]string{
    "left_x":  "543",
    "left_y":  "409",
    "right_x": "603",
    "right_y": "474",
  }
  appInfo := map[string]string{
    "app_name":    "通道可视化智能报警程序",
    "app_version": "1.0.0",
  }
  body := gin.H{
    "cluster_id":     "sgcc",
    "title":          "test 数据",
    "alert_type":     "crane",
    "description":    "reliability: 0.56563453244342",
    "reliability":    0.56563453244342,
    "alert_position": []map[string]string{position},
    "app_info":       appInfo,
    "create_time":      time.Now().Unix(),
    "app_id":         "8sdsfc9sdsd9s8ds89d8s9dds9s",
    "event_id":       "5fc8752c-c599-11e9-94fe-0242ac110005",
    "device_id":      "5d4a9c1d97e91b6667bcec93",
    "image_path":     "media/5d4a9c1d97e91b6667bcec93/image/2019-08-23/11-30-00.jpg",
    "image":          "/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCABLAHkDASIAAhEBAxEB/8QAGwAAAgMBAQEAAAAAAAAAAAAAAQIDBAUABgf/xAAuEAACAgIAAwgBAwUBAAAAAAAAAQIDBBEFEjETISJBUWFxgTMGIzIUUnLB0fD/xAAZAQADAQEBAAAAAAAAAAAAAAAAAQIDBAX/xAAfEQACAgIDAQEBAAAAAAAAAAAAAQIRAyESIjEEMkH/2gAMAwEAAhEDEQA/APOpDqJyQ6Wz3DzgKLfRDcrHgmh9bYrAjUSRIZIbQrACQyOSCABCBBQxDoZCIdAAyHQqGQgJEMIhgAwFFskiipxDNWBhTyHBz5fJepJi8c4JmRjzX2YNj7uW+PNDf+S/2ZzzwhLjJmsccpK0W0OkMobTlBqyG9KcHuL+GPGDet9y9iuSe0RTXommMovzIsK63Ix3OyMHFSahKK02k31XQv49Mrfxw5pPoiY5FJWVKDTorJHcppSwLVD8T3v0K1lEqpalFplKSZLTK3KFIkcQaKEBDIGhkAhkMgIZIBjIYCQ+gA8/dRVk0TpuW65rUlvRVjwPA7Ls51O3xcylN+JfaLyHiKWOMnbRSlJaTKlPDsrF2uHZs6k0/wBuxc0f/fQjyeL8Pbnl4fa1+dtEuZL68jUi2vJkXGMq+PBslVTcJcj1qKbft3nNPAoq4aNo5W3Uth4NlYt9U6o2NdklKTaa1tpa+dnosPKp5tqUW1/b5nzf9LQlXxOdWS73ZZHnXJY+Vx8+b1XTr5nrYcOxu3jONt8OVNKKn4fv1MsXN49lZOKmenlm+/Xyl/0p5d1N3R+L3MeX9TU2o2dpHfdt94iypRbV1bRrGo+6IdvwsTjp93QjDG6Fi8M18MPz3HQmmYtUBIZI5IdIoRyQyRyQyQgDFDgQRAYKJYpESJYstgWIyTXeYnHuGznjX5OPZbOxpvs5SXLFa6pexrpgsjGyqdc1uMotNeq0ZyipKi4yp2eV/TtU3xpSrs/bVSlNtd79NHsk+8weCx5cmT5FFOmKS3t9TcRl8y6F533Jdh2mteQiCjYyBKiuS/jr4I+wsj+Oz6ZOgoTihqTIO0urXir37xJYZFcur1rrslijpVwn/KKf0TUl4x2n6h4aktxaa9UNorvFS/HOUH8hbyaouTSsS9Oo+TXqCl/CwkEr4eS8qpzlRbS09asjrfuWNDTtWiWqMBEiEQ6NGIdDIVDw6iGY/CtRz1F2OU3S1rfkpG4jG4aks1JJLumunubaOb5n0Ns/6CkNoCGNzEKGQqHQAFIdICGQgCECCgAKDoC6sYQH/9k=",
    "time":           time.Now().Unix(),
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/gateway/import", reader)
  data, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, "push raw data to dataSet success & assistant success", msg)
  dataArr := data.(map[string]interface{})
  for _, item := range dataArr {
    assert.NotEmpty(t, item)
  }
}
