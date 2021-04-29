package v1

import (
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/datasets/models/database"
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
  "go.mongodb.org/mongo-driver/bson/primitive"
  log "k8s.io/klog"
)

/**
  dataSet 管理列表
*/
func QueryAllDataSetsForReMarkerHandler(c *gin.Context) {
  var body DataSetAdminBody
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  dataSetTarget, _ := database.QueryDataSetDataByName(body.Name)
  if dataSetTarget.Name == "" {
    response.Success(c, "dataSet data not found", "")
    return
  }
  sortString := "timestamp"
  dataSets, err := database.QueryRawDataListByIndex(body.StartIndex-1, body.Offset, sortString, dataSetTarget.Sets)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  response.Success(c, fmt.Sprintf("query %s dataSet index %d to %d success", body.Name, body.StartIndex-1, body.StartIndex+body.Offset-1), dataSets)
}

/**
  上传手动标注数据
*/
func UploadHandleMarkerImageHandler(c *gin.Context) {
  var body ImageDataBody
  err := c.BindJSON(&body)
  log.Info("upload handle image info :", body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  id, _ := primitive.ObjectIDFromHex(body.Id)
  var rawData = interfaces.Data{Path: body.Path, Name: body.Name, TimeStamp: body.TimeStamp, Id: id, LabelName: body.LabelName, Label: body.Label}
  err = database.ReplaceRawData(rawData)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  response.Success(c, fmt.Sprintf("update dataSet (%s) manage  success", body.Id), "")
}
