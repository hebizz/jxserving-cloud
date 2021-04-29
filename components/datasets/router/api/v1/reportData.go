package v1

import (
  "errors"
  "fmt"
  "net/http"
  "strconv"
  "strings"
  "time"

  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/datasets/models/database"
  "gitlab.jiangxingai.com/jxserving/components/datasets/pkg/config"
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
  "gitlab.jiangxingai.com/jxserving/components/utils"
  "go.mongodb.org/mongo-driver/bson/primitive"
  log "k8s.io/klog"
)

const (
  TypeDataMarkRectangle = 1
  TypeDataMarkCycle     = 2
)

/**
  上传真实上报数据
*/
func UploadReportDataHandler(c *gin.Context) {
  var body interfaces.ReportData
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  //解析description字段中的reliability值
  if body.Description != "" {
    arr := strings.Split(body.Description, " ")
    if len(arr) == 2 {
      reliability, err := strconv.ParseFloat(arr[1], 64)
      if err != nil {
        log.Error("reliability format error", err)
      }
      body.Reliability = reliability
    }
  }
  // 提取部分数据(name,path,label)
  //将base64字符串转image后存储到本地目录
  var rawData = interfaces.Data{}
  // time.Now().UnixNano()/1e6 : 毫秒
  // time.Now().Unix()         : 秒
  // time.Now().UnixNano()     : 纳秒
  fileName := utils.GetMD5Str(time.Now().UnixNano()/1e6) + ".jpg"
  filePath := config.ImageSavePath() + utils.GetDate()
  rawData.Path, rawData.Name, err = utils.Base64ToImage(body.Image, filePath, fileName)
  rawData.TimeStamp = body.AlertTime
  // ai assistant 上传的图片统一存入default数据集
  rawData.LabelName = "default"
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  if body.DataId == "" {
    // 通过该途径上传的数据，传递过来的AI监测数据，全部置空
    rawData.Label = make([]interfaces.Label, 0)
  } else {
    var labelList []interfaces.Label
    for _, obj := range body.AlertPosition {
      var label interfaces.Label
      label.N = body.AlertType
      label.T = TypeDataMarkRectangle
      //TODO  数据标注有可能出现cycle类型
      var marks []string
      marks = append(marks, obj.LeftX, obj.LeftY, obj.RightX, obj.RightY)
      label.D = marks
      labelList = append(labelList, label)
    }
    rawData.Label = labelList
  }
  rawDataId, err := database.InsertRawData(rawData)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  // 插入上报原始数据
  objectIDS, err := primitive.ObjectIDFromHex(rawDataId)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  body.Id = objectIDS
  body.CreateTime = time.Now().Unix()
  // 如果project为空,默认将数据归为reportData
  if body.Project == "" {
    body.Project = "reportData"
  }
  errInsertReport := database.InsertReportData(body)
  if errInsertReport != nil {
    response.ServerError(c, http.StatusInternalServerError, "", errInsertReport, "")
    return
  }
  // 插入到数据集
  var collectionName = "default"
  dataSetTarget, _ := database.QueryDataSetDataByName(collectionName)
  var errDataSets error
  if dataSetTarget.Name == "" {
    dataSetTarget.Name = collectionName
    dataSetTarget.CreatedTimestamp = time.Now().Unix()
    dataSetTarget.Sets = append(dataSetTarget.Sets, rawDataId)
    // 插入该数据
    errDataSets = database.InsertDataSetsData(*dataSetTarget)
  } else {
    errDataSets = database.UpdateDataSetsData(collectionName, rawDataId)
  }
  if errDataSets != nil {
    response.ServerError(c, http.StatusInternalServerError, "", errors.New("insert report data error"), "")
    return
  }
  response.Success(c, fmt.Sprintf("upload report data name = %s success", body.Project), rawDataId)
}
