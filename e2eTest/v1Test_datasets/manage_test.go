package v1Test_datasets

import (
  "encoding/json"
  "fmt"
  "testing"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
  "gitlab.jiangxingai.com/jxserving/components/datasets/router/api/router"
  "gitlab.jiangxingai.com/jxserving/e2eTest"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type Label struct {
  N string   `json:"name"` // label name like fire crane..
  T int      `json:"type"` // 1 rectangle 2 cycle 3 ..
  D []string `json:"position"`
}

func TestQueryAllDataSetsForReMarkerHandler(t *testing.T) {
  type Data struct {
    Id        primitive.ObjectID `bson:"_id" json:"id"`
    Name      string             `json:"name"`
    TimeStamp int64              `json:"timestamp"`
    Path      string             `json:"path"`
    Label     []Label            `json:"label"`
  }
  body := gin.H{
    "name":       "sgcc",
    "startIndex": 1,
    "offset":     3,
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/dataSet/manage/query", reader)
  data, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("query %s dataSet index %d to %d success", body["name"], body["startIndex"].(int)-1, body["offset"].(int)+body["startIndex"].(int)-1), msg)
  dataArr := data.([]interface{})
  var dataBody Data
  for _, item := range dataArr {
    arr, err := json.Marshal(item)
    assert.Nil(t, err)
    err = json.Unmarshal(arr, &dataBody)
    assert.NotEmpty(t, dataBody.Id)
    assert.NotEmpty(t, dataBody.Name)
    assert.NotEmpty(t, dataBody.Label)
    assert.NotEmpty(t, dataBody.Path)
    assert.NotEmpty(t, dataBody.TimeStamp)
  }
}

func TestUploadHandleMarkerImageHandler(t *testing.T) {
  pos := []string{"185", "277", "268", "310"}
  labelArray1 := Label{N: "yanhuo", T: 1, D: pos}
  labelArray2 := Label{N: "fire", T: 1, D: pos}

  body := gin.H{
    "id":        "5da57df9c52c9d626d97d129",
    "timestamp": time.Now().Unix(),
    "image":     "",
    "name":      "sgcc",
    "path":      "/data/edgebox/remote/image/2019-10-15/Smoke_00018.jpg",
    "label":     []Label{labelArray1, labelArray2},
  }
  r, reader, err := e2eTest.GeneratePostRequestBody(body)
  r = router.InitRouter(r)
  w := e2eTest.PerformRequest(r, "POST", "/api/v1/dataSet/manage/update", reader)
  _, msg := e2eTest.AssertResponseOk(t, w, err)
  assert.Equal(t, fmt.Sprintf("update dataSet (%s) manage  success", body["id"]), msg)
}
